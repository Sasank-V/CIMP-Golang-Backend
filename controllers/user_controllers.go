package controllers

import (
	"fmt"
	"sync"

	"github.com/Sasank-V/CIMP-Golang-Backend/database"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var UserColl *mongo.Collection
var userConnect sync.Once

func connectUserCollection() {
	userConnect.Do(func() {
		db := database.InitDB()
		schemas.CreateUserCollection(db)
		UserColl = db.Collection("users")
	})
}

func UserExist(id string) bool {
	connectUserCollection()
	ctx, cancel := database.GetContext()
	defer cancel()

	count, err := UserColl.CountDocuments(ctx, bson.M{"id": id})
	if err != nil {
		fmt.Printf("Error checking for existing user: %v\n", err)
		return false
	}
	return count > 0
}

func GetUserByID(id string) (schemas.User, error) {
	connectUserCollection()
	ctx, cancel := database.GetContext()
	defer cancel()

	var user schemas.User
	err := UserColl.FindOne(ctx, bson.D{{"id", id}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return schemas.User{}, err
		}
		fmt.Printf("Error fetching user data: %v", err)
		return schemas.User{}, err
	}
	return user, nil
}

func AddUser(user schemas.User) error {
	connectUserCollection()
	ctx, cancel := database.GetContext()
	defer cancel()
	if UserExist(user.ID) {
		return fmt.Errorf("User Already exist\n")
	}
	_, err := UserColl.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(user_id string) error {
	connectUserCollection()
	ctx, cancel := database.GetContext()
	defer cancel()
	res := UserColl.FindOneAndDelete(ctx, bson.D{{"id", user_id}})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no User found with the given ID")
		}
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
