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
	DB               *databaseData
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
	authCredentials := options.Credential{Username: d.DB.User, Password: d.DB.Password}
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

// GetData returns information about the database, which was parsed from JSON.
func (d *mongoDatabase) GetData() *databaseData {
	return d.DB
}

// Init creates the db connection string and context object.
func (d *mongoDatabase) Init() {
	d.connectionString = fmt.Sprintf(`mongodb://%s:%d/%s`,
		d.DB.Host,
		d.DB.Port,
		d.DB.Name,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	d.ctx = ctx
	d.close = cancel
}

// Select selects data from the database, with or without filters.
func (d *mongoDatabase) Select(tableName string, conditions string) []map[string]interface{} {
	var allDocuments []map[string]interface{}

	// defer d.close()

	d.TestConnection()

	client := d.GetClient()
	collection := client.Database(d.DB.Name).Collection(tableName)

	var bsonConditions interface{}
	if conditions != "" && conditions != "*" {
		err := bson.UnmarshalExtJSON([]byte(conditions), true, &bsonConditions)
		if err != nil {
			panic(err)
		}
	} else {
		bsonConditions = bson.M{}
	}

	cur, err := collection.Find(d.ctx, bsonConditions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var documentMap map[string]interface{} = make(map[string]interface{})
		err := cur.Decode(documentMap)
		if err != nil {
			log.Fatal(err)
		}
		allDocuments = append(allDocuments, documentMap)
	}

	defer cur.Close(context.TODO())

	// for _, element := range allDocuments {
	// 	book := element
	// 	fmt.Println(book)
	// }

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

// Update updates a document with the privided key.
func (d *mongoDatabase) Update(table string, key interface{}, column string, val interface{}) (bool, error) {
	client := d.GetClient()
	collection := client.Database(d.DB.Name).Collection(table)
	filter := bson.D{{"name", "Ash"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(updateResult)
	return false, nil
}
