package routes

import (
	"net/http"

	"github.com/Sasank-V/CIMP-Golang-Backend/controllers"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"github.com/Sasank-V/CIMP-Golang-Backend/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUserRoutes(r *gin.RouterGroup) {
	r.GET("/get/:id", getUserInfo)
	r.GET("/get-requests/:id", getUserRequests)
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

func getUserRequests(c *gin.Context) {
	userID := c.Param("id")
	user, err := controllers.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, types.AuthResponse{
				Message: "No User found with the given ID",
			})
		}
	}
}

func getUserContributionData(c *gin.Context) {

}
