package datastoremongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MaxRecords int64 = 100
var Client *mongo.Client
var Ctx context.Context
var Cancel context.CancelFunc
var Connected bool = false

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancels the context.
func Close() {

	if Connected {
		// CancelFunc to cancel to context
		defer Cancel()

		// Client provides a method to close
		// a mongoDB connection.
		defer func() {

			Connected = false
			// Client.Disconnect method also has deadline.
			// returns error if any,
			if err := Client.Disconnect(Ctx); err != nil {
				if err == mongo.ErrClientDisconnected {
					//  ignore
				} else {
					panic(err)
				}
			}
		}()
	}
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resources associated with it.
func Connect(uri string) error {

	// Ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	Ctx, Cancel = context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	var err error = nil
	if !Connected {
		Client, err = mongo.Connect(Ctx, options.Client().ApplyURI(uri))
		Connected = true
	}
	return err
}

func Ping() error {
	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := Client.Ping(Ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func InsertOne(dataBase string, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	// select database and collection with Client.Database method
	// and Database.Collection method
	collection := Client.Database(dataBase).Collection(col)
	// InsertOne accepts two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(Ctx, doc)
	return result, err
}

func UpdateOne(dataBase string, collectionName string,
	filter bson.M, operator string, doc interface{}) (*mongo.UpdateResult, error) {

	// select database and collection with Client.Database method
	// and Database.Collection method
	collection := Client.Database(dataBase).Collection(collectionName)

	result, err := collection.UpdateOne(Ctx, filter, bson.M{"$set": doc})
	return result, err
}

func DeleteRecordByUuid(dataBase string, collectionName string, uuid string) error {
	collection := Client.Database(dataBase).Collection(collectionName)
	filter := bson.M{"uuid": uuid}
	_, err := collection.DeleteOne(Ctx, filter)
	return err
}

func DeleteCollectionRecords(dataBase string, collectionName string, filter bson.M) error {
	collection := Client.Database(dataBase).Collection(collectionName)
	_, err := collection.DeleteMany(Ctx, filter)
	return err
}

func GetRecords(dataBase string, collectionName string,
	start int64, limit int64, filter bson.M) (*mongo.Cursor, error) {

	collection := Client.Database(dataBase).Collection(collectionName)

	if start < 0 {
		start = 0
	}
	if limit > MaxRecords {
		limit = MaxRecords
	}

	findOptions := options.FindOptions{
		Skip:  &start,
		Limit: &limit,
	}
	cursor, err := collection.Find(Ctx, filter, &findOptions)

	return cursor, err
}

func GetRecord(dataBase string, collectionName string, filter bson.M) *mongo.SingleResult {
	collection := Client.Database(dataBase).Collection(collectionName)
	result := collection.FindOne(Ctx, filter)
	return result
}
