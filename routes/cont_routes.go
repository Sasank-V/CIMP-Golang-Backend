package routes

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/Sasank-V/CIMP-Golang-Backend/controllers"
	"github.com/Sasank-V/CIMP-Golang-Backend/database/schemas"
	"github.com/Sasank-V/CIMP-Golang-Backend/types"
	"github.com/Sasank-V/CIMP-Golang-Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupContributionRoutes(r *gin.RouterGroup) {
	r.POST("/add", addContribution)
	r.PATCH("/update/details", updateContributionDetails)
	r.PATCH("/update/status", updateContributionStatus)
}

func addContribution(c *gin.Context) {
	var contInfo types.ContributionInfo
	if err := c.ShouldBindBodyWithJSON(&contInfo); err != nil {
		log.Printf("Error decoding the JSON Body: %v", err)
		c.JSON(http.StatusBadRequest, types.AddContributionResponse{
			Message:        "Error decoding the JSON Body",
			ContributionID: "",
		})
		return
	}

	if contInfo.Points < 0 {
		log.Printf("Points less than zero")
		c.JSON(http.StatusBadRequest, types.AddContributionResponse{
			Message:        "Points value less than zero",
			ContributionID: "",
		})
		return
	}

	new_cont := schemas.Contribution{
		ID:          utils.GetContributionID(contInfo.UserID),
		Title:       contInfo.Title,
		Points:      uint(contInfo.Points),
		UserID:      contInfo.UserID,
		Description: contInfo.Description,
		ProofFiles:  contInfo.ProofFiles,
		Target:      contInfo.Target,
		SecTargets:  contInfo.SecTargets,
		ClubID:      contInfo.ClubID,
		Department:  contInfo.Department,
		Status:      schemas.Pending,
		CreatedAt:   time.Now(),
	}

	newContID, err := controllers.AddContribution(new_cont)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.AddContributionResponse{
			Message:        "Error adding contribution",
			ContributionID: "",
		})
		return
	}
	c.JSON(http.StatusOK, types.AddContributionResponse{
		Message:        "Contribution Added Successfully",
		ContributionID: newContID,
	})
}

func updateContributionDetails(c *gin.Context) {
	var ContInfo types.ContributionUpdateInfo
	if err := c.ShouldBindBodyWithJSON(&ContInfo); err != nil {
		log.Printf("error parsing JSON body: %v", err)
		c.JSON(http.StatusBadRequest, types.MessageResponse{
			Message: "Error parsing JSON body",
		})
		return
	}

	cont, err := controllers.GetContributionByID(ContInfo.ContributionID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No Contribution found with the given contribution id: %v", err)
			c.JSON(http.StatusNotFound, types.MessageResponse{
				Message: "No Contribution found with the given ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, types.MessageResponse{
			Message: "Error fetching contribution details",
		})
		return
	}

	if cont.Contribution.UserID != ContInfo.UserID {
		c.JSON(http.StatusUnauthorized, types.MessageResponse{
			Message: "User with give user id is not the owner of this contribution",
		})
		return
	}

	if cont.Contribution.Status == schemas.Approved {
		c.JSON(http.StatusBadRequest, types.MessageResponse{
			Message: "User Contribution Already Approved",
		})
		return
	}

	err = controllers.UpdateContributionDetails(ContInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.MessageResponse{
			Message: "Error updating the details",
		})
		return
	}
	c.JSON(http.StatusOK, types.MessageResponse{
		Message: "Contribution Details Updated Successfully",
	})

}

func updateContributionStatus(c *gin.Context) {
	var updateInfo types.ContributionStatusInfo
	if err := c.ShouldBindBodyWithJSON(&updateInfo); err != nil {
		log.Printf("error parsing JSON body: %v", err)
		c.JSON(http.StatusBadRequest, types.MessageResponse{
			Message: "Error parsing JSON body data",
		})
		return
	}

	var user schemas.User
	var cont types.FullContribution
	var userErr, contErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func(id string) {
		defer wg.Done()
		user, userErr = controllers.GetUserByID(id)
	}(updateInfo.LeadUserID)

	go func(id string) {
		defer wg.Done()
		cont, contErr = controllers.GetContributionByID(id)
	}(updateInfo.ContributionID)

	wg.Wait()

	if userErr != nil {
		if userErr == mongo.ErrNoDocuments {
			log.Printf("no user found with the given id: %v", userErr)
			c.JSON(http.StatusNotFound, types.MessageResponse{
				Message: "No Lead User found with the given ID",
			})
			return
		}
		log.Printf("error fetching lead user details: %v", userErr)
		c.JSON(http.StatusInternalServerError, types.MessageResponse{
			Message: "Error fetching lead user details",
		})
		return
	}

	if contErr != nil {
		if contErr == mongo.ErrNoDocuments {
			log.Printf("no contribution found with the given id: %v", contErr)
			c.JSON(http.StatusNotFound, types.MessageResponse{
				Message: "No Contribution found with the given ID",
			})
			return
		}
		log.Printf("error fetching contribution details: %v", contErr)
		c.JSON(http.StatusInternalServerError, types.MessageResponse{
			Message: "Error fetching contribution details",
		})
		return
	}

	if !user.IsLead {
		c.JSON(http.StatusUnauthorized, types.MessageResponse{
			Message: "Give Lead User is NOT a Lead",
		})
		return
	}

	if cont.Contribution.Target != updateInfo.LeadUserID && !slices.Contains(cont.Contribution.SecTargets, updateInfo.LeadUserID) {
		c.JSON(http.StatusUnauthorized, types.MessageResponse{
			Message: "You are not the Targeted Lead",
		})
		return
	}

	err := controllers.UpdateContributionStatus(updateInfo.ContributionID, updateInfo.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.MessageResponse{
			Message: fmt.Sprintf("Error updating status: %v", err),
		})
		return
	}

	//Increase the total_points field for the user if pending -> approved
	//Decrease the total_points field for the user if approved -> rejected

	c.JSON(http.StatusOK, types.MessageResponse{
		Message: "Contribution Status Updated Successfully",
	})

}
