package tokenizer

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"tokentarpon/crypto"
	"tokentarpon/tokenizer/datastore"

	"github.com/google/uuid"
)

var DefaultMaxLength int = 255
var MaxRecords int64 = 100
var UnitTest = false
var tokenRecordType = "token"

type Token struct {
	Uuid           string `bson:"uuid" json:"uuid"`
	DomainUuid     string `bson:"domainUuid" json:"domainUuid"`
	Value          string `bson:"value" json:"value"`
	EncryptedValue string `bson:"encryptedValue" json:"encryptedValue"`
	IsDeleted      bool   `json:"isdeleted" bson:"isdeleted"`
	DocumentType   string `bson:"documentType" json:"documentType"`
	Version        string `bson:"version" json:"version"`
	Created        int64  `json:"created" bson:"created"`
	Updated        int64  `json:"updated" bson:"updated"`
	Check          string `bson:"check" json:"check"`
}

type Token_v001 struct {
	Uuid           string `bson:"uuid" json:"uuid"`
	DomainUuid     string `bson:"domainUuid" json:"domainUuid"`
	Value          string `bson:"value" json:"value"`
	EncryptedValue string `bson:"encryptedValue" json:"encryptedValue"`
	IsDeleted      bool   `json:"isdeleted" bson:"isdeleted"`
	DocumentType   string `bson:"documentType" json:"documentType"`
	Version        string `bson:"version" json:"version"`
	Created        int64  `json:"created" bson:"created"`
	Updated        int64  `json:"updated" bson:"updated"`
	Check          string `bson:"check" json:"check"`
}

type TokenQuery struct {
	DomainUuid string   `bson:"domainUuid" json:"domainUuid"`
	Uuids      []string `bson:"uuids" json:"uuids"`
}

type TokenError struct {
	Token Token  `bson:"token" json:"token"`
	Error string `bson:"error" json:"error"`
}

var (
	ErrEmptyValue      = errors.New("cannot store empty value")
	ErrValueTooBig     = errors.New("value too large for storage")
	ErrNoMatchingToken = errors.New("no token found for provided domain and token id")
)

func CreateToken(collectionName string, domainUuid string, value string) (Token, error) {
	var tok Token

	err := errors.New("data: token incomplete, need domain id, value")
	if len(strings.TrimSpace(domainUuid)) == 0 {
		return tok, err
	}
	if len(strings.TrimSpace(value)) == 0 {
		return tok, err
	}
	aUuid := uuid.New()
	tok.DomainUuid = domainUuid
	tok.Value = value
	tok.Uuid = aUuid.String()
	tok.Created = time.Now().Unix()

	if UnitTest {
		return tok, nil
	}

	datastore.CollectionName = collectionName
	result, errInsert := datastore.InsertRecord(tokenRecordType, tok)
	if errInsert != nil {
		return tok, errInsert
	} else if nil == result {
		return tok, datastore.ErrNotFound
	} else {
		insertedToken := result.(Token)
		return insertedToken, nil
	}
}

func CreateTokens(collectionName, domainUuid string, tokens []Token) ([]Token, []TokenError) {
	var createdTokens []Token
	var errorTokens []TokenError
	datastore.CollectionName = collectionName

	for _, tokenObj := range tokens {
		if len(strings.TrimSpace(tokenObj.DomainUuid)) == 0 {
			e := TokenError{Token: tokenObj, Error: "Missing Domain ID"}
			errorTokens = append(errorTokens, e)
		} else if tokenObj.DomainUuid != domainUuid {
			e := TokenError{Token: tokenObj, Error: "Invalid Domain ID"}
			errorTokens = append(errorTokens, e)
		} else if len(strings.TrimSpace(tokenObj.Value)) == 0 {
			e := TokenError{Token: tokenObj, Error: "Missing Token Value"}
			errorTokens = append(errorTokens, e)
		} else {
			aUuid := uuid.New()
			tokenObj.Uuid = aUuid.String()
			tokenObj.Created = time.Now().Unix()

			if UnitTest {
				createdTokens = append(createdTokens, tokenObj)
			} else {
				result, errInsert := datastore.InsertRecord(tokenRecordType, tokenObj)
				if errInsert != nil {
					errmsg := fmt.Sprint(errInsert)
					e := TokenError{Token: tokenObj, Error: errmsg}
					errorTokens = append(errorTokens, e)
				} else if nil == result {
					e := TokenError{Token: tokenObj, Error: "The token was inserted but not returned"}
					errorTokens = append(errorTokens, e)
				} else {
					tok := result.(Token)
					createdTokens = append(createdTokens, tok)
				}
			}
		}
	}
	return createdTokens, errorTokens
}

func GetToken(collectionName string, domainUuid string, tokenUuid string) (Token, error) {
	var tok Token
	err := errors.New("data: need domain id, token id")
	if len(strings.TrimSpace(domainUuid)) == 0 {
		return tok, err
	}
	if len(strings.TrimSpace(tokenUuid)) == 0 {
		return tok, err
	}

	if UnitTest {
		tok.DomainUuid = domainUuid
		tok.Uuid = tokenUuid
		return tok, nil
	}

	filter := datastore.MakeDomainQuery(domainUuid, "uuid", tokenUuid, true)
	datastore.CollectionName = collectionName
	result, err := datastore.GetRecord(tokenRecordType, filter)
	if err == nil {
		tok = result.(Token)
	}
	return tok, err
}

func DeleteToken(collectionName string, tokenObj Token) (Token, error) {
	var empty Token
	err := errors.New("data: token incomplete, need domain id, token id, value")
	if len(strings.TrimSpace(tokenObj.DomainUuid)) == 0 {
		return empty, err
	}
	if len(strings.TrimSpace(tokenObj.Uuid)) == 0 {
		return empty, err
	}

	if UnitTest {
		tokenObj.IsDeleted = true
		return tokenObj, nil
	}

	filter := datastore.MakeDomainQuery(tokenObj.DomainUuid, "uuid", tokenObj.Uuid, false)
	datastore.CollectionName = collectionName
	result, err := datastore.GetRecord(tokenRecordType, filter)
	if err != nil {
		return empty, err
	}

	var t Token = result.(Token)
	t.IsDeleted = true
	updateResult, err := datastore.UpdateRecord(tokenRecordType, filter, "and", t)
	updatedToken := updateResult.(Token)
	return updatedToken, err
}

func CreateMultiTokenQuery(tokenQuery TokenQuery) []datastore.DataQueryGroup {

	var filters = make([]datastore.DataQueryGroup, 2)
	var nvq datastore.DataQueryGroup

	nvq.Operator = "and"
	nvq.DataQueries = make([]datastore.DataQuery, 1)
	nvq.DataQueries[0].FieldName = "domainUuid"
	nvq.DataQueries[0].FieldValue = tokenQuery.DomainUuid
	nvq.DataQueries[0].Wildcard = false
	nvq.DataQueries[0].CaseSensitive = true
	filters[0] = nvq

	nvq.Operator = "or"
	nvq.DataQueries = make([]datastore.DataQuery, len(tokenQuery.Uuids))
	for idx, uuid := range tokenQuery.Uuids {
		//fmt.Println("Slice item", i, "is", numbers[i])
		nvq.DataQueries[idx].FieldName = "uuid"
		nvq.DataQueries[idx].FieldValue = uuid
		nvq.DataQueries[idx].Wildcard = false
		nvq.DataQueries[idx].CaseSensitive = true
	}
	filters[1] = nvq
	return filters
}

func GetTokens(collectionName string, domainUuid string) ([]Token, error) {
	var empty []Token
	err := errors.New("data: need domain id")
	if len(strings.TrimSpace(domainUuid)) == 0 {
		return empty, err
	}

	if UnitTest {
		return empty, err
	}

	datastore.CollectionName = collectionName
	filter := datastore.MakeSimpleQuery("domainUuid", domainUuid, false)
	results, err := datastore.GetRecords(tokenRecordType, filter, "and", 0, 100)
	if err != nil {
		return nil, datastore.ErrQueryError
	} else {
		tokens := make([]Token, 0)
		for _, x := range results {
			tokens = append(tokens, x.(Token))
		}
		return tokens, nil
	}
}

func GetTokenValues(collectionName string, tokenQuery TokenQuery) ([]string, error) {
	var empty []string
	err := errors.New("data: need domain id")
	if len(strings.TrimSpace(tokenQuery.DomainUuid)) == 0 {
		return empty, err
	}
	err = errors.New("data: need uuids")
	if len(tokenQuery.Uuids) == 0 {
		return empty, err
	}

	if UnitTest {
		return empty, err
	}

	filter := CreateMultiTokenQuery(tokenQuery)
	datastore.CollectionName = collectionName
	results, err := datastore.GetRecords(tokenRecordType, filter, "and", 0, MaxRecords)
	if err != nil {
		return nil, datastore.ErrQueryError
	} else {
		tokenValues := make([]string, 0)

		// return the token values at the same indices as their uuids
		for _, uuid := range tokenQuery.Uuids {
			for _, tv := range results {
				if uuid == tv.(Token).Uuid {
					tokenValues = append(tokenValues, tv.(Token).Value)
				}
			}
		}
		return tokenValues, nil
	}

}

func EncryptValues(plaintext []string, key string) []string {
	encryptedValues := make([]string, 0)
	for _, plain := range plaintext {
		encrypted, err := crypto.EncryptAES(plain, key)
		if nil == err {
			encryptedValues = append(encryptedValues, encrypted)
		}
	}
	return encryptedValues
}

func DecryptValues(encryptedValues []string, key string) []string {
	plaintext := make([]string, 0)
	for _, enc := range encryptedValues {
		encrypted, err := crypto.DecryptAES(enc, key)
		if nil == err {
			plaintext = append(plaintext, encrypted)
		}
	}
	return plaintext
}
