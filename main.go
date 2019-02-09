package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"log"
)

type Trainer struct {
	Name string
	Age int
	City string
}


func main() {
	/* =============================================================== */
	/* | Sample Data                                                 | */
	/* =============================================================== */
	ash := Trainer{"Ash", 10, "Pallet Twon"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	/* =============================================================== */
	/* | Connected to MongoDB                                        | */
	/* =============================================================== */
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MangoDB!")

	/* =============================================================== */
	/* | Get collection                                              | */
	/* =============================================================== */
	collection := client.Database("test").Collection("trainers")

	/* =============================================================== */
	/* | Insert a single document                                    | */
	/* =============================================================== */
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document:", insertResult.InsertedID)

	/* =============================================================== */
	/* | Insert multiple documents at a time                         | */
	/* =============================================================== */
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents", insertManyResult.InsertedIDs)

	/* =============================================================== */
	/* | Update documents                                            | */
	/* =============================================================== */
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
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	/* =============================================================== */
	/* | Find single document                                        | */
	/* =============================================================== */
	//create a value into which the result can be decoded
	var result Trainer
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)

	/* =============================================================== */
	/* | Find multiple documents                                     | */
	/* =============================================================== */
	options := options.Find()
	options.SetLimit(2)
	var results []*Trainer
	cur, err := collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	/* =============================================================== */
	/* | Delete single Document                                      | */
	/* =============================================================== */
	delFilter := bson.D{{"name", "Ash"}}
	deleteResult, err := collection.DeleteOne(context.TODO(), delFilter, nil)

	fmt.Printf("Found a single document: %+v\n", deleteResult)
	/* =============================================================== */
	/* | Connection to MangoDB closed.                               | */
	/* =============================================================== */
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MangoDB closed.")

}