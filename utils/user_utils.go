package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Sasank-V/CIMP-Golang-Backend/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var User *mongo.Collection
var userConnect sync.Once

func connectUserCollection() {
	userConnect.Do(func() {
		db := database.InitDB()
		User = db.Collection("users")
	})
}

func GetUserByID(id string) (bson.M, error) {
	connectUserCollection()
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
