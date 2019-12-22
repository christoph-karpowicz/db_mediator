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

type MongoDatabase struct {
	DB               *DatabaseData
	Tst              int
	connectionString string
	ctx              context.Context
	close            context.CancelFunc
}

func (d *MongoDatabase) CloseConnection() {
	d.close()
}

func (d *MongoDatabase) GetClient() *mongo.Client {
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

func (d *MongoDatabase) GetData() *DatabaseData {
	return d.DB
}

func (d *MongoDatabase) Init() {
	d.connectionString = fmt.Sprintf(`mongodb://%s:%d/%s`,
		d.DB.Host,
		d.DB.Port,
		d.DB.Name,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	d.ctx = ctx
	d.close = cancel
}

func (d *MongoDatabase) Select(tableName string, conditions string) []map[string]interface{} {
	var allDocuments []map[string]interface{}

	// defer d.close()

	d.TestConnection()

	client := d.GetClient()
	collection := client.Database(d.DB.Name).Collection(tableName)

	var bsonConditions interface{}
	if conditions != "-" && conditions != "" {
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

func (d *MongoDatabase) TestConnection() {
	c := d.GetClient()
	err := c.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
}
