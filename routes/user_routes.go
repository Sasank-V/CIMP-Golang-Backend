package routes

import (
	"log"
	"net/http"
	"sync"

	"github.com/Sasank-V/CIMP-Golang-Backend/controllers"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"github.com/Sasank-V/CIMP-Golang-Backend/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Route : /api/user
func SetupUserRoutes(r *gin.RouterGroup) {
	r.GET("/info/:id", getUserInfo)
	r.GET("/contributions/:id", getUserContributions)
	r.GET("/get-contribution-data/:id", getUserContributionData)
}

func getUserInfo(c *gin.Context) {
	userID := c.Param("id")
	user, err := controllers.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, types.GetUserResponse{
				Message: "No User found with the given ID",
				User:    schemas.User{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, types.GetUserResponse{
			Message: "Error fetching user data, Try Again Later",
			User:    schemas.User{},
		})
		return
	}
	c.JSON(http.StatusOK, types.GetUserResponse{
		Message: "User Data retrived successfully",
		User:    user,
	})
}

func getUserContributions(c *gin.Context) {
	userID := c.Param("id")
	user, err := controllers.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, types.GetUserRequestsResponse{
				Message:       "No User found with the given ID",
				Contributions: []types.FullContribution{},
			})
		}
	}

	contChan := make(chan types.FullContribution, len(user.Contributions))
	errChan := make(chan error, len(user.Contributions))
	var wg sync.WaitGroup

	for _, contID := range user.Contributions {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			fullCont, err := controllers.GetContributionByID(id)
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
		log.Printf("Error fetching Contribution :", err)
		c.JSON(http.StatusInternalServerError, types.GetUserRequestsResponse{
			Message:       "Error while getting Contribution Details",
			Contributions: []types.FullContribution{},
		})
		return
	}

	c.JSON(http.StatusOK, types.GetUserRequestsResponse{
		Message:       "User request fetched successfully",
		Contributions: contributions,
	})
}

func getUserContributionData(c *gin.Context) {

}
