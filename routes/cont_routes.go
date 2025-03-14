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
)

func SetupContributionRoutes(r *gin.RouterGroup) {
	r.POST("/add", addContribution)
	r.PATCH("/update", updateContribution)
}

type ContributionInfo struct {
	Title       string   `json:"title"`
	Points      int      `json:"points"`
	UserID      string   `json:"user_id"`
	Description string   `json:"description"`
	ProofFiles  []string `json:"proof_files,omitempty"`
	Target      string   `json:"target"`
	secTargets  []string `json:"secTargets,omitempty"`
	ClubID      string   `json:"club_id"`
	Department  string   `json:"department"`
}

func addContribution(c *gin.Context) {
	var contInfo ContributionInfo
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
		SecTargets:  contInfo.secTargets,
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

func updateContribution(c *gin.Context) {

}
