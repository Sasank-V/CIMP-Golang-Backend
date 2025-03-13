package controllers

import (
	"fmt"
	"log"
	"sync"

	"github.com/Sasank-V/CIMP-Golang-Backend/database"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"github.com/Sasank-V/CIMP-Golang-Backend/lib"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var UserColl *mongo.Collection
var userConnect sync.Once

func ConnectUserCollection() {
	userConnect.Do(func() {
		db := database.InitDB()
		schemas.CreateUserCollection(db)
		UserColl = db.Collection(lib.UserCollName)
	})
}

func UserExist(id string) bool {
	ctx, cancel := database.GetContext()
	defer cancel()

	count, err := UserColl.CountDocuments(ctx, bson.M{"id": id})
	if err != nil {
		log.Printf("error checking for existing user: %v", err)
		return false
	}
	return count > 0
}

func GetUserByID(id string) (schemas.User, error) {
	ctx, cancel := database.GetContext()
	defer cancel()

	var user schemas.User
	err := UserColl.FindOne(ctx, bson.D{{"id", id}}).Decode(&user)
	if err != nil {
		log.Printf("error fetching user data: %v", err)
		if err == mongo.ErrNoDocuments {
			return schemas.User{}, err
		}
		return schemas.User{}, err
	}
	return user, nil
}

func AddUser(user schemas.User) error {
	ctx, cancel := database.GetContext()
	defer cancel()
	if UserExist(user.ID) {
		return fmt.Errorf("user already exist")
	}
	_, err := UserColl.InsertOne(ctx, user)
	if err != nil {
		log.Printf("error adding user : %v", err)
		return err
	}
	return nil
}

func DeleteUser(user_id string) error {
	ctx, cancel := database.GetContext()
	defer cancel()
	res := UserColl.FindOneAndDelete(ctx, bson.D{{"id", user_id}})
	if err := res.Err(); err != nil {
		log.Printf("error deleting user: %v", err)
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no User found with the given ID")
		}
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
