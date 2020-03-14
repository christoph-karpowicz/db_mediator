package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDatabase implements Database interface for MongoDB database.
type mongoDatabase struct {
	cfg              *config
	connectionString string
	ctx              context.Context
	close            context.CancelFunc
}

// CloseConnection closes the db connection.
func (d *mongoDatabase) CloseConnection() {
	d.close()
}

// GetClient returns a connection client object.
func (d *mongoDatabase) GetClient() *mongo.Client {
	authCredentials := options.Credential{Username: d.cfg.User, Password: d.cfg.Password}
	clientOptions := options.Client().ApplyURI(d.connectionString).SetAuth(authCredentials)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(d.ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// GetConfig returns information about the database, which was parsed from JSON.
func (d *mongoDatabase) GetConfig() *config {
	return d.cfg
}

// Init creates the db connection string and context object.
func (d *mongoDatabase) Init() {
	d.connectionString = fmt.Sprintf(`mongodb://%s:%d/%s`,
		d.cfg.Host,
		d.cfg.Port,
		d.cfg.Name,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	d.ctx = ctx
	d.close = cancel

	d.TestConnection()
}

// Insert inserts one row into a given collection.
func (d *mongoDatabase) Insert(table string, keyName string, keyVal interface{}, values map[string]interface{}) (bool, error) {
	client := d.GetClient()
	collection := client.Database(d.cfg.Name).Collection(table)

	insertResult, err := collection.InsertOne(context.TODO(), values)
	if err != nil {
		dbErr := &DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error(), KeyName: keyName, KeyVal: keyVal}
		return false, dbErr
	}
	if insertResult.InsertedID == nil {
		dbErr := &DatabaseError{DBName: d.cfg.Name, ErrMsg: "could not insert document", KeyName: keyName, KeyVal: keyVal}
		return false, dbErr
	}

	return true, nil
}

// Select selects data from the database, with or without filters.
func (d *mongoDatabase) Select(tableName string, conditions string) []map[string]interface{} {
	var allDocuments []map[string]interface{}

	// defer d.close()

	client := d.GetClient()
	collection := client.Database(d.cfg.Name).Collection(tableName)

	var bsonConditions interface{}
	if conditions != "" && conditions != "*" {
		err := bson.UnmarshalExtJSON([]byte(conditions), true, &bsonConditions)
		if err != nil {
			panic(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
		}
	} else {
		bsonConditions = bson.M{}
	}

	cur, err := collection.Find(d.ctx, bsonConditions)
	if err != nil {
		log.Fatal(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
	}

	for cur.Next(context.TODO()) {
		var documentMap map[string]interface{} = make(map[string]interface{})
		err := cur.Decode(documentMap)
		if err != nil {
			log.Fatal(&DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error()})
		}
		allDocuments = append(allDocuments, documentMap)
	}

	defer cur.Close(context.TODO())

	return allDocuments
}

// TestConnection pings the database.
func (d *mongoDatabase) TestConnection() {
	c := d.GetClient()
	err := c.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
}

// Update updates a document with the provided key.
func (d *mongoDatabase) Update(table string, keyName string, keyVal interface{}, column string, val interface{}) (bool, error) {
	client := d.GetClient()
	collection := client.Database(d.cfg.Name).Collection(table)
	filter := bson.D{{keyName, keyVal}}
	update := bson.D{
		{"$set", bson.D{
			{column, val},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		dbErr := &DatabaseError{DBName: d.cfg.Name, ErrMsg: err.Error(), KeyName: keyName, KeyVal: keyVal}
		return false, dbErr
	}
	if updateResult.MatchedCount == 0 {
		dbErr := &DatabaseError{DBName: d.cfg.Name, ErrMsg: "document with given key not found", KeyName: keyName, KeyVal: keyVal}
		return false, dbErr
	}
	// if updateResult.ModifiedCount == 0 {
	// 	dbErr := &DatabaseError{DBName: d.cfg.Name, ErrMsg: "no documents modified", KeyName: keyName, KeyVal: keyVal}
	// 	return false, dbErr
	// }

	return true, nil
}
