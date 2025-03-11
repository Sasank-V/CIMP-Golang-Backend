package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Sasank-V/CIMP-Golang-Backend/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var DBClient *mongo.Client
var User *mongo.Collection

func InitDB() {
	DBClient = database.ConnectDB()
	DBname := os.Getenv("DATABASE_NAME")
	if DBname == "" {
		log.Fatal("No DATBASE_NAME found in env.")
	}
	User = DBClient.Database(DBname).Collection("users")
}

func GetUserByID(id string) (bson.M, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user bson.M
	err := User.FindOne(ctx, bson.D{{"id", id}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bson.M{}, err
		}
		fmt.Printf("Error fetching user data: %v", err)
		return bson.M{}, err
	}
	return user, nil
}
