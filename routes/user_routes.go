package routes

import (
	"net/http"

	"github.com/Sasank-V/CIMP-Golang-Backend/utils"
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
	user, err := utils.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No User found with the given ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching user data, Try Again Later",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User Data retrived successfully",
		"user":    user,
	})
}

func getUserRequests(c *gin.Context) {

}

func getUserContributionData(c *gin.Context) {

}
