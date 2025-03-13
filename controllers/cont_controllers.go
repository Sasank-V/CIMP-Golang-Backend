package controllers

import (
	"log"
	"sync"

	"github.com/Sasank-V/CIMP-Golang-Backend/database"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"github.com/Sasank-V/CIMP-Golang-Backend/lib"
	"github.com/Sasank-V/CIMP-Golang-Backend/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ContColl *mongo.Collection
var ContConnect sync.Once

func ConnectContributionCollection() {
	ContConnect.Do(func() {
		db := database.InitDB()
		schemas.CreateContributionCollection(db)
		ContColl = db.Collection(lib.ContributionCollName)
	})
}

func GetContributionByID(id string) (types.FullContribution, error) {
	ctx, cancel := database.GetContext()
	defer cancel()
	var cont schemas.Contribution
	err := ContColl.FindOne(ctx, bson.D{{"id", id}}).Decode(&cont)
	if err != nil {
		log.Printf("error getting contribution data: %v", err)
		return types.FullContribution{}, err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var club schemas.Club
	var dept schemas.Department
	var clubErr, deptErr error

	go func() {
		club, clubErr = GetClubByID(cont.ClubID)
		defer wg.Done()
	}()

	go func() {
		dept, deptErr = GetDepartmentByID(cont.Department)
		defer wg.Done()
	}()

	wg.Wait()

	if clubErr != nil {
		return types.FullContribution{}, err
	}
	if deptErr != nil {
		return types.FullContribution{}, err
	}

	return types.FullContribution{
		Contribution:   cont,
		ClubName:       club.Name,
		DepartmentName: dept.Name,
	}, nil
}
