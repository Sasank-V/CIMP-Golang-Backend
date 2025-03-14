package routes

import (
	"log"
	"net/http"
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
		c.JSON(http.StatusBadRequest, types.UpdateContributionDetailsResponse{
			Message: "Error parsing JSON body",
		})
		return
	}

	cont, err := controllers.GetContributionByID(ContInfo.ContributionID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No Contribution found with the given contribution id: %v", err)
			c.JSON(http.StatusNotFound, types.UpdateContributionDetailsResponse{
				Message: "No Contribution found with the given ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, types.UpdateContributionDetailsResponse{
			Message: "Error fetching contribution details",
		})
		return
	}

	if cont.Contribution.UserID != ContInfo.UserID {
		c.JSON(http.StatusUnauthorized, types.UpdateContributionDetailsResponse{
			Message: "User with give user id is not the owner of this contribution",
		})
		return
	}

	if cont.Contribution.Status == schemas.Approved {
		c.JSON(http.StatusBadRequest, types.UpdateContributionDetailsResponse{
			Message: "User Contribution Already Approved",
		})
		return
	}

	err = controllers.UpdateContributionDetails(ContInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.UpdateContributionDetailsResponse{
			Message: "Error updating the details",
		})
		return
	}
	c.JSON(http.StatusOK, types.UpdateContributionDetailsResponse{
		Message: "Contribution Details Updated Successfully",
	})

}
