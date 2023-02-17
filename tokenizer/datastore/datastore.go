// Datastore encapsulates the nitty gritty of getting and setting data
// Currently using mongodb for storage, but could be any storage mechanism
package datastore

import (
	"errors"
	"reflect"
	"strings"

	"tokentarpon/tokenizer/datastore/datastoremongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var MaxRecords int64 = 100
var MongoUri = ""
var MongoDatabase = ""
var CollectionName = ""

var tokenVersion string = "001"

// type Token struct {
// 	Uuid           string `bson:"uuid" json:"uuid"`
// 	DomainUuid     string `bson:"domainUuid" json:"domainUuid"`
// 	Value          string `bson:"value" json:"value"`
// 	EncryptedValue string `bson:"encryptedValue" json:"encryptedValue"`
// 	IsDeleted      bool   `json:"isdeleted" bson:"isdeleted"`
// 	DocumentType   string `bson:"documentType" json:"documentType"`
// 	Version        string `bson:"version" json:"version"`
// 	Created        int64  `json:"created" bson:"created"`
// 	Updated        int64  `json:"updated" bson:"updated"`
// 	Check          string `bson:"check" json:"check"`
// }

// type Token_v001 struct {
// 	Uuid           string `bson:"uuid" json:"uuid"`
// 	DomainUuid     string `bson:"domainUuid" json:"domainUuid"`
// 	Value          string `bson:"value" json:"value"`
// 	EncryptedValue string `bson:"encryptedValue" json:"encryptedValue"`
// 	IsDeleted      bool   `json:"isdeleted" bson:"isdeleted"`
// 	DocumentType   string `bson:"documentType" json:"documentType"`
// 	Version        string `bson:"version" json:"version"`
// 	Created        int64  `json:"created" bson:"created"`
// 	Updated        int64  `json:"updated" bson:"updated"`
// 	Check          string `bson:"check" json:"check"`
// }

// used for querying the db using grouped name-value pairs
// e.g. DataQueries contains 2 DataQuery values
type DataQueryGroup struct {
	DataQueries []DataQuery `bson:"dataQueries" json:"dataQueries"`
	Operator    string      `bson:"operator" json:"operator"`
}

type DataQuery struct {
	FieldName     string `bson:"fieldName" json:"fieldName"`
	FieldValue    string `bson:"fieldValue" json:"fieldValue"`
	Negate        bool   `bson:"negate" json:"negate"`
	IsBool        bool   `bson:"isBool" json:"isBool"`
	BoolValue     bool   `bson:"boolValue" json:"boolValue"`
	IdValue       string `bson:"idValue" json:"idValue"`
	CaseSensitive bool   `bson:"caseSensitive" json:"caseSensitive"`
	Wildcard      bool   `bson:"wildcard" json:"wildcard"`
}

var (
	ErrServerError    = errors.New("data: a server error occurred")
	ErrDatastoreError = errors.New("data: a datastore error occurred")
	ErrQueryError     = errors.New("data: query returned an error and may be malformed")
	ErrConflict       = errors.New("data: record cannot be overwritten")
	ErrNotFound       = errors.New("data: record not found")
)

// GetRecord takes an incoming query against a target table/collection
// And returns a single record in the form of an interface
func GetRecord(queryParams []DataQueryGroup, record interface{}) error {
	err := datastoremongo.Connect(MongoUri)
	if err != nil {
		return err
	}
	filter := CreateMongoFilter(queryParams, "and")
	result := datastoremongo.GetRecord(MongoDatabase, CollectionName, filter)
	if nil == result {
		return ErrNotFound
	} else if nil != result.Err() {
		return result.Err()
	}
	mongoerr := result.Decode(record)
	return mongoerr
}

// GetRecords takes an incoming query against a target table/collection
// And returns an array of records in the form of an array of interfaces
func GetRecords(queryParams []DataQueryGroup,
	operator string, start int64, limit int64, record interface{}) ([]interface{}, error) {

	var results []interface{}

	connecterr := datastoremongo.Connect(MongoUri)
	if connecterr != nil {
		return results, ErrDatastoreError
	}

	filter := CreateMongoFilter(queryParams, operator)
	mongocursor, mongoerr := datastoremongo.GetRecords(MongoDatabase,
		CollectionName, start, limit, filter)
	if nil != mongoerr {
		return results, ErrQueryError
	}

	myType := reflect.TypeOf(record)
	tokens := reflect.MakeSlice(reflect.SliceOf(myType), 0, 0).Interface()
	if err := mongocursor.All(datastoremongo.Ctx, &tokens); err != nil {
		return nil, ErrQueryError
	}

	// because of using reflection we have to repack the results to return
	tokenValues := reflect.ValueOf(tokens)
	for i := 0; i < tokenValues.Len(); i++ {
		v := tokenValues.Index(i).Interface()
		results = append(results, v)
	}
	return results, nil
}

func InsertRecord(recordType string, document interface{}) error {

	err := datastoremongo.Connect(MongoUri)
	if err != nil {
		return ErrDatastoreError
	}

	datastoremongo.InsertOne(MongoDatabase,
		CollectionName, document)

	result, err := datastoremongo.InsertOne(MongoDatabase,
		CollectionName, document)
	if err != nil {
		return ErrQueryError
	}

	// Get the whole record-
	// use the returned InsertedID (bson.ObjectID) from the mongo.InsertOneResult
	var idString string = hexFromObjectId(result.InsertedID)
	filter := makeIdQuery(idString)
	geterr := GetRecord(filter, &document)
	return geterr
}

func DeleteRecord(uuid string) error {
	err := datastoremongo.Connect(MongoUri)
	if err != nil {
		return err
	}
	mongoerr := datastoremongo.DeleteRecordByUuid(MongoDatabase, CollectionName, uuid)
	return mongoerr
}

func DeleteRecords(queryParams []DataQueryGroup, operator string) error {
	err := datastoremongo.Connect(MongoUri)
	if err != nil {
		return err
	}
	filter := CreateMongoFilter(queryParams, operator)
	mongoerr := datastoremongo.DeleteCollectionRecords(MongoDatabase, CollectionName, filter)
	return mongoerr
}

func UpdateRecord(recordType string,
	queryParams []DataQueryGroup, operator string,
	document interface{}) (interface{}, error) {

	err := datastoremongo.Connect(MongoUri)
	if err != nil {
		return nil, ErrDatastoreError
	}

	filter := CreateMongoFilter(queryParams, operator)
	// updateDoc := bson.D{
	// 	{"$set", bson.D{doc}},
	// }

	_, updateerr := datastoremongo.UpdateOne(MongoDatabase, CollectionName, filter, "and", document)
	if updateerr != nil {
		//fmt.Println(result)
		return nil, ErrQueryError
	}
	return document, nil
}

// func updateChecksum(recordType string, document interface{}) interface{} {

// 	// update the Check hash
// 	if recordType == "token" {
// 		var x Token = document.(Token)
// 		x.Updated = time.Now().Unix()
// 		if x.Created == 0 {
// 			x.Created = x.Updated
// 		}
// 		x.DocumentType = "token"
// 		x.Version = tokenVersion
// 		x.Check = ""
// 		j, _ := json.Marshal(x)
// 		x.Check = crypto.GetHashForByteArray(j)
// 		document = x
// 	}
// 	return document
// }

// func ValidateChecksum(recordType string, recordId string, record interface{}) (bool, error) {

// 	var checkOk bool = false
// 	var err error = nil
// 	var docTypeNotSpecified error = errors.New("Record is missing document type or version")
// 	var noSuchTypeErr error = errors.New("Record document type or version not found")

// 	// the document itself contains the DocumentType and Version
// 	// we need both of those in order to test the checksum

// 	// first we have to load the record into the specified recordType
// 	// in order to get the DocumentType and Version

// 	// then use structure DocumentType_Version to load the document,
// 	// marshal it, and generate the checksum to compare to the existing check

// 	var existingChecksum string = ""
// 	var storedTypeVersion string = ""

// 	query := MakeSimpleQuery("uuid", recordId, true)
// 	geterr := GetRecord(query, record)
// 	if geterr != nil {
// 		return checkOk, geterr
// 	}

// 	jsonDoc, err := json.Marshal(record)

// 	if recordType == "token" {
// 		var x Token = record.(Token)
// 		existingChecksum = x.Check
// 		if len(x.DocumentType) == 0 || len(x.Version) == 0 {
// 			err = docTypeNotSpecified
// 		} else {
// 			storedTypeVersion = x.DocumentType + "_v" + x.Version
// 			if storedTypeVersion == "token_v001" {
// 				var x2 Token_v001
// 				json.Unmarshal(jsonDoc, &x2)

// 				// unset the check value and marshal back into json
// 				// to generate a comparison hash
// 				x2.Check = ""
// 				j, _ := json.Marshal(x2)
// 				x2.Check = crypto.GetHashForByteArray(j)
// 				if x2.Check == existingChecksum {
// 					checkOk = true
// 				}
// 			} else {
// 				err = noSuchTypeErr
// 			}
// 		}
// 	}
// 	return checkOk, err
// }

// MakeSimpleQuery creates a simple name-value pair and generates
// the array necessary to call any of the datastore Get functions
func MakeSimpleQuery(fieldName string, fieldValue string, caseSensitive bool) []DataQueryGroup {

	var filters = make([]DataQueryGroup, 1)
	var nvq DataQueryGroup

	nvq.Operator = "and"
	nvq.DataQueries = make([]DataQuery, 2)

	nvq.DataQueries[0].FieldName = fieldName
	nvq.DataQueries[0].FieldValue = fieldValue
	nvq.DataQueries[0].Wildcard = false
	nvq.DataQueries[0].CaseSensitive = caseSensitive

	nvq.DataQueries[1].FieldName = "isDeleted"
	nvq.DataQueries[1].BoolValue = false
	nvq.DataQueries[1].IsBool = true

	filters[0] = nvq

	return filters
}

func MakeDomainQuery(domainUuid string, fieldName string, fieldValue string,
	caseSensitive bool) []DataQueryGroup {

	var filters = make([]DataQueryGroup, 1)
	var nvq DataQueryGroup

	nvq.Operator = "and"
	nvq.DataQueries = make([]DataQuery, 3)

	nvq.DataQueries[0].FieldName = "domainUuid"
	nvq.DataQueries[0].FieldValue = domainUuid
	nvq.DataQueries[0].Wildcard = false
	nvq.DataQueries[0].CaseSensitive = caseSensitive

	nvq.DataQueries[1].FieldName = fieldName
	nvq.DataQueries[1].FieldValue = fieldValue
	nvq.DataQueries[1].Wildcard = false
	nvq.DataQueries[1].CaseSensitive = caseSensitive

	nvq.DataQueries[2].FieldName = "isDeleted"
	nvq.DataQueries[2].BoolValue = false
	nvq.DataQueries[2].IsBool = true

	filters[0] = nvq
	return filters
}

func CreateMongoFilter(queryValues []DataQueryGroup, operator string) bson.M {

	var result bson.M

	resultValues := make(bson.A, len(queryValues))
	for k, v := range queryValues {
		if len(v.DataQueries) == 1 {
			// no need for the operator just create a simple bson.M query
			result = createSimpleFilter(v.DataQueries[0])
			resultValues[k] = result
		} else if len(v.DataQueries) > 1 {
			result = createComplexFilter(v)
			resultValues[k] = result
		}
	}

	if len(resultValues) > 1 {
		result = make(bson.M, 1)
		result["$"+operator] = resultValues
	}
	return result
}

func Close() {
	defer datastoremongo.Close()
}

func createSimpleFilter(dataQuery DataQuery) bson.M {

	var bsonQuery bson.M

	//if len(strings.TrimSpace(dataQuery.IdValue)) > 0 {
	//@todo instead check if IdValue has a value

	if len(strings.TrimSpace(dataQuery.IdValue)) > 0 {
		idVal, err := primitive.ObjectIDFromHex(dataQuery.IdValue)
		if nil == err {
			bsonQuery = bson.M{"_id": idVal}
			// bson.ObjectIdHex(dataQuery.IdValue)
			// bson.ObjectID(dataQuery.IdValue)}
		}
	} else if dataQuery.IsBool {
		bsonQuery = bson.M{dataQuery.FieldName: bson.M{"$eq": dataQuery.BoolValue}}
	} else if dataQuery.Wildcard {
		if dataQuery.CaseSensitive {
			bsonQuery = bson.M{dataQuery.FieldName: bson.M{"$regex": dataQuery.FieldValue}}
		} else {
			// case-insensitive wildcard search using value v against field k
			bsonQuery = bson.M{dataQuery.FieldName: bson.M{"$regex": dataQuery.FieldValue, "$options": "i"}}
		}
	} else {
		if dataQuery.CaseSensitive {
			if dataQuery.Negate {
				bsonQuery = bson.M{dataQuery.FieldName: bson.M{"$ne": dataQuery.FieldValue}}
			} else {
				bsonQuery = bson.M{dataQuery.FieldName: dataQuery.FieldValue}
			}
		} else {
			// note: no forward slashes!
			// most of the documentation indicates you should include forward slashes in the regular expression,
			// but if you do this the expression will work with .Find() but not with .FindOne()
			// super confusing and can take a while to troubleshoot
			bsonQuery = bson.M{dataQuery.FieldName: bson.M{"$regex": "^" + dataQuery.FieldValue + "$", "$options": "i"}}
		}
	}

	return bsonQuery
}

func createComplexFilter(dataQueryGroup DataQueryGroup) bson.M {

	var bsonQuery bson.M
	subQueries := make(bson.A, len(dataQueryGroup.DataQueries))

	for k, v := range dataQueryGroup.DataQueries {
		//bsonA[k] = CreateSimpleFilter(v)
		subQueries[k] = createSimpleFilter(v)
	}

	// bson.M{
	// 	"$and": bson.A{
	// 		bson.M{v.QueryName.Name: v.QueryName.Value},
	// 		bson.M{v.QueryValue.Name: v.QueryValue.Value},
	// 	},
	// }

	if len(subQueries) > 0 {
		bsonQuery = make(bson.M, 1)
		bsonQuery["$"+dataQueryGroup.Operator] = subQueries // bson.A{subQueries}
	}

	return bsonQuery
}
func hexFromObjectId(id interface{}) string {
	objectId := id.(primitive.ObjectID)
	return objectId.Hex()
}

// MakeIdQuery creates a simple query against the mongo record's _id value
// returning the array necessary to call any of the data Get functions
func makeIdQuery(id string) []DataQueryGroup {

	var filters = make([]DataQueryGroup, 1)
	var nvq DataQueryGroup

	nvq.Operator = "and"
	nvq.DataQueries = make([]DataQuery, 1)
	nvq.DataQueries[0].IdValue = id
	//  interface {} is primitive.ObjectID, not []uint8
	filters[0] = nvq

	return filters
}
