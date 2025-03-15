package controllers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Sasank-V/CIMP-Golang-Backend/api/types"
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

//Get Functions

func GetUserByID(id string) (schemas.User, error) {
	ctx, cancel := database.GetContext()
	defer cancel()

	var user schemas.User
	err := UserColl.FindOne(ctx, bson.D{{"id", id}}).Decode(&user)
	if err != nil {
		log.Printf("error fetching user data: %v", err)
		return schemas.User{}, err
	}
	return user, nil
}

func GetAllUserInClub(id string) ([]schemas.User, error) {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"clubs": bson.M{"$in": []string{id}},
	}

	cursor, err := UserColl.Find(ctx, filter)
	if err != nil {
		log.Printf("error fetching club users: %v", err)
		return []schemas.User{}, err
	}

	var members []schemas.User
	if err = cursor.All(ctx, &members); err != nil {
		log.Printf("cursor error: %v", err)
		return []schemas.User{}, err
	}
	return members, nil
}

func GetAllUserContributions(id string) ([]types.FullContribution, error) {
	user, err := GetUserByID(id)
	if err != nil {
		log.Printf("error fetching user Data: %v", err)
		return []types.FullContribution{}, err
	}

	contChan := make(chan types.FullContribution, len(user.Contributions))
	errChan := make(chan error, len(user.Contributions))
	var wg sync.WaitGroup

	for _, contID := range user.Contributions {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			fullCont, err := GetContributionByID(id)
			if err != nil {
				errChan <- err
				return
			}
			contChan <- fullCont
		}(contID)
	}

	go func() {
		wg.Wait()
		close(contChan)
		close(errChan)
	}()

	var contributions []types.FullContribution
	for cont := range contChan {
		contributions = append(contributions, cont)
	}

	if len(errChan) > 0 {
		err = <-errChan
		log.Printf("Error fetching Contribution :", err)
		return []types.FullContribution{}, err
	}
	return contributions, nil
}

//Add/Update Functions

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

func AddContributionIDToUser(user_id string, cont_id string) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"id": user_id,
	}
	update := bson.M{
		"$push": bson.M{"contributions": cont_id},
	}

	res, err := UserColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error adding contID to user Contributions: %v", err)
		return err
	}

	if res.ModifiedCount == 0 {
		log.Printf("No User found with the given ID")
		return mongo.ErrNoDocuments
	}
	return nil
}

func UpdateUserTotalPoints(userID string, points int) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	user, err := GetUserByID(userID)
	if err != nil {
		log.Printf("error fetching user: %v", err)
		return err
	}

	if points < 0 && user.TotalPoints < points {
		log.Printf("user points cannot be negative")
		return fmt.Errorf("user points cannot be negative")
	}
	filter := bson.M{
		"id": userID,
	}
	update := bson.M{
		"$inc": bson.M{"total_points": points},
	}
	res, err := UserColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating user total_points: %v", err)
		return err
	}

	if res.ModifiedCount == 0 {
		log.Printf("No user found with the given ID")
		return mongo.ErrNoDocuments
	}

	return nil
}

//OTP and Pass Reset Functions

func SetResetOTPToUser(user_id string, otp string) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"id": user_id,
	}
	update := bson.M{
		"$set": bson.M{
			"otp": otp,
		},
	}
	res, err := UserColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error setting up user OTP: %v", err)
		return err
	}
	if res.ModifiedCount == 0 {
		log.Printf("No User found to set OTP")
		return mongo.ErrNoDocuments
	}
	return nil
}

func IncreaseUserOTPRetries(user_id string) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"id": user_id,
	}
	update := bson.M{
		"$inc": bson.M{
			"otp_retries": 1,
		},
	}
	res, err := UserColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error increasing the opt retry count: %v", err)
		return err
	}
	if res.ModifiedCount == 0 {
		log.Printf("No User found with the given ID")
		return mongo.ErrNoDocuments
	}
	return nil
}

func LockUserAccountPasswordReset(user_id string, time time.Time) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"id": user_id,
	}
	update := bson.M{
		"$set": bson.M{
			"locked_till": time,
		},
	}

	res, err := UserColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error locking the account: %v", err)
		return err
	}
	if res.ModifiedCount == 0 {
		log.Printf("No User found to lock")
		return mongo.ErrNoDocuments
	}
	return nil
}

func ResetOTPandLockValuesForUser(user_id string) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"id": user_id,
	}
	update := bson.M{
		"$unset": bson.M{
			"otp":         "",
			"locked_till": nil,
			"otp_retries": 0,
		},
	}
	res, err := UserColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error reseting the Values: %v", err)
		return err
	}
	if res.ModifiedCount == 0 {
		log.Printf("No user found with the ID to reset")
		return mongo.ErrNoDocuments
	}
	return nil
}

func SetNewPasswordToUser(user_id string, pass string) error {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"id": user_id,
	}
	update := bson.M{
		"$set": bson.M{
			"password": pass,
		},
	}
	res, err := UserColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error setting up new password to user")
		return err
	}
	if !res.Acknowledged {
		log.Printf("No User was set the new Password: %v", res)
		return mongo.ErrNoDocuments
	}
	return nil

}

// Delete Functions
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
