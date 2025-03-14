package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sasank-V/CIMP-Golang-Backend/controllers"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"github.com/Sasank-V/CIMP-Golang-Backend/types"
	"github.com/Sasank-V/CIMP-Golang-Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupAuthRoutes(r *gin.RouterGroup) {
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)
}

func signupHandler(c *gin.Context) {
	var user types.UserSignUpInfo

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		fmt.Printf("Error in signup: %v", err)
		c.JSON(http.StatusBadRequest, types.AuthResponse{
			Message: "Error in the data sent , Not a JSON Object",
			Token:   "",
		})
		return
	}

	new_user := schemas.User{
		ID:        utils.GetUserIDFromRegNumber(user.RegNumber),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  utils.HashSHA256(user.Password),
		RegNumber: strings.ToUpper(user.RegNumber),
		IsLead:    false,
	}
	if err := controllers.AddUser(new_user); err != nil {
		fmt.Printf("Error adding user: %v\n", err)
		c.JSON(http.StatusInternalServerError, types.AuthResponse{
			Message: fmt.Sprintf("Error adding user: %s , Try Again Later", err),
			Token:   "",
		})
		return
	}

	payload := types.TokenPayload{
		ID:     new_user.ID,
		Name:   new_user.FirstName + " " + new_user.LastName,
		IsLead: new_user.IsLead,
	}
	token, err := controllers.GenerateToken(payload)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		controllers.DeleteUser(new_user.ID)
		c.JSON(http.StatusInternalServerError, types.AuthResponse{
			Message: "Error while generating Token , Try Again Later",
			Token:   "",
		})
		return
	}

	c.JSON(http.StatusOK, types.AuthResponse{
		Message: "User Signup Successfull",
		Token:   token,
	})
}

func loginHandler(c *gin.Context) {
	var login types.UserLoginInfo

	if jerr := c.ShouldBindBodyWithJSON(&login); jerr != nil {
		fmt.Printf("Error logining in : %v", jerr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error loggin in the data sent, Not a JSON object",
		})
		return
	}
	user, uerr := controllers.GetUserByID(utils.GetUserIDFromRegNumber(login.RegNumber))
	if uerr != nil {
		if uerr == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, types.AuthResponse{
				Message: "No User found with the given ID",
				Token:   "",
			})
			return
		}
		fmt.Printf("Error loggin in : %v", uerr)
		c.JSON(http.StatusInternalServerError, types.AuthResponse{
			Message: "Some Error occured while logging in , Try Again Later",
			Token:   "",
		})
		return
	}

	pass_hash := utils.HashSHA256(login.Password)
	if user.Password != pass_hash {
		c.JSON(http.StatusBadRequest, types.AuthResponse{
			Message: "Incorrect Password",
			Token:   "",
		})
		return
	}

	payload := types.TokenPayload{
		ID:     user.ID,
		Name:   user.FirstName + " " + user.LastName,
		IsLead: user.IsLead,
	}
	token, terr := controllers.GenerateToken(payload)
	if terr != nil {
		fmt.Printf("Error while generating token: %v", terr)
		c.JSON(http.StatusInternalServerError, types.AuthResponse{
			Message: "Server Error while creating token , Try again later",
			Token:   "",
		})
		return
	}

	c.JSON(http.StatusOK, types.AuthResponse{
		Message: "User Logged in Successfully",
		Token:   token,
	})

}
