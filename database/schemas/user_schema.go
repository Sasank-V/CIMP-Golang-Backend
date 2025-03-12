package schemas

import (
	"log"

	"slices"

	"github.com/Sasank-V/CIMP-Golang-Backend/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID            string   `bson:"id" json:"id"`
	RegNumber     string   `bson:"reg_number" json:"reg_number"`
	FirstName     string   `bson:"first_name" json:"first_name"`
	LastName      string   `bson:"last_name" json:"last_name"`
	Email         string   `bson:"email" json:"email"`
	Password      string   `bson:"password" json:"password"`
	IsLead        bool     `bson:"is_lead" json:"is_lead"`
	Departments   []string `bson:"departments,omitempty" json:"departments,omitempty"`
	Clubs         []string `bson:"clubs,omitempty" json:"clubs,omitempty"`
	Contributions []string `bson:"contributions,omitempty" json:"contributions,omitempty"`
}

func SetUserDefaults(user *User) {
	if user.Departments == nil {
		user.Departments = []string{}
	}
	if user.Clubs == nil {
		user.Departments = []string{}
	}
	if user.Contributions == nil {
		user.Contributions = []string{}
	}
}

func setUserUniqueKeys(coll *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"id", 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := database.GetContext()
	defer cancel()

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal("Error while setting up unique key: ", err)
	}
}

func CreateUserCollection(db *mongo.Database) {
	ctx, cancel := database.GetContext()
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("Error getting the Collection Names: ", err)
		return
	}
	if slices.Contains(collections, "users") {
		log.Printf("User Collection already exist , skipping creation")
		return
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"id", "reg_number", "first_name", "last_name", "email", "password", "is_lead"},
		"properties": bson.M{
			"id": bson.M{
				"bsonType": "string",
			},
			"reg_number": bson.M{
				"bsonType": "string",
			},
			"first_name": bson.M{
				"bsonType": "string",
			},
			"last_name": bson.M{
				"bsonType": "string",
			},
			"email": bson.M{
				"bsonType": "string",
			},
			"password": bson.M{
				"bsonType": "string",
			},
			"is_lead": bson.M{
				"bsonType": "bool",
			},
			"departments": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"bsonType": "string",
				},
			},
			"clubs": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"bsonType": "string",
				},
			},
			"contributions": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"bsonType": "string",
				},
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	opts := options.CreateCollection().SetValidator(validator)
	err = db.CreateCollection(ctx, "users", opts)
	if err != nil {
		log.Fatal("Failed to Create User Collection: ", err)
		return
	}

	setUserUniqueKeys(db.Collection("users"))
}
