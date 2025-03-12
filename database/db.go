package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	DBClient   *mongo.Client
	clientOnce sync.Once
)

var User *mongo.Collection

func InitDB() *mongo.Database {
	connectDB()
	DBname := os.Getenv("DATABASE_NAME")
	if DBname == "" {
		log.Fatal("No DATBASE_NAME found in env.")
	}
	return DBClient.Database(DBname)
}

func connectDB() *mongo.Client {
	clientOnce.Do(func() {
		uri := os.Getenv("CONNECTION_STRING")
		if uri == "" {
			log.Fatal("Set your 'CONNECTION_STRING' environment variable. ")
		}
		dbClient, err := mongo.Connect(options.Client().
			ApplyURI(uri))
		if err != nil {
			log.Fatal("[MONGO-DB] Failed to connect to MongoDB: ", err)
		}
		if err := dbClient.Ping(context.TODO(), nil); err != nil {
			log.Fatal("[MONGO-DB] MongoDB connection test failed: ", err)
		}
		fmt.Printf("[MONGO-DB] MongoDB Connected\n")
		DBClient = dbClient
	})
	return DBClient
}
