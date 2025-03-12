package routes

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup) {
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)
}

func hashSHA256(inp string) string {
	hash := sha256.Sum256([]byte(inp))
	return fmt.Sprintf("%x", hash)
}

type User struct {
	RegNumber string `json:"reg_number"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func signupHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		fmt.Printf("Error in signup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in the data sent , Not a JSON Object",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Signup Successfull",
	})
}

func loginHandler(c *gin.Context) {

}
