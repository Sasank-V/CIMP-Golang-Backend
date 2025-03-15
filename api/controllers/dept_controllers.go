package controllers

import (
	"log"
	"sync"

	"github.com/Sasank-V/CIMP-Golang-Backend/database"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"github.com/Sasank-V/CIMP-Golang-Backend/lib"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var DeptColl *mongo.Collection
var DeptConnect sync.Once

func ConnectDepartmentCollection() {
	DeptConnect.Do(func() {
		db := database.InitDB()
		schemas.CreateDepartmentCollection(db)
		DeptColl = db.Collection(lib.DepartmentCollName)
	})
}

func GetDepartmentByID(id string) (schemas.Department, error) {
	ctx, cancel := database.GetContext()
	defer cancel()

	var dept schemas.Department
	err := DeptColl.FindOne(ctx, bson.D{{"id", id}}).Decode(&dept)
	if err != nil {
		log.Printf("error getting deparment data: ", err)
		return schemas.Department{}, err
	}
	return dept, nil
}

func GetAllDepartmentsInClub(id string) ([]schemas.Department, error) {
	ctx, cancel := database.GetContext()
	defer cancel()

	filter := bson.M{
		"club_id": id,
	}

	cursor, err := DeptColl.Find(ctx, filter)
	if err != nil {
		log.Printf("Error getting departments in Club: %v", err)
		return []schemas.Department{}, err
	}

	var departments []schemas.Department
	if err = cursor.All(ctx, &departments); err != nil {
		log.Printf("cursor error: %v", err)
		return []schemas.Department{}, err
	}
	return departments, nil
}
