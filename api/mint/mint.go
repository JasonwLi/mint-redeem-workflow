package mint

import (
	"mint-redeem-workflow/db"
	"mint-redeem-workflow/deps"
	"mint-redeem-workflow/models"
	"mint-redeem-workflow/service"
	"mint-redeem-workflow/worker/workflows"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ProcessMintFunc = service.ProcessMint

type MintRedeemRequest struct {
	Amount    float64 `json:"amount" binding:"required"`
	Recipient string  `json:"recipient" binding:"required"`
}

func HandleMintRedeemRequest(c *gin.Context) {
	var req MintRedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	request := models.Request{
		ID:        uuid.New(),
		Type:      "mint",
		Amount:    req.Amount,
		Recipient: req.Recipient,
	}

	workflowInput := workflows.MintInput{
		Amount:    req.Amount,
		Recipient: req.Recipient,
		RequestID: request.ID.String(),
	}

	db := db.Db

	cadenceClient, err := deps.BuildCadenceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ProcessMintFunc(db, &request, workflowInput, cadenceClient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": request.ID.String(), "status": "workflow started"})
}
