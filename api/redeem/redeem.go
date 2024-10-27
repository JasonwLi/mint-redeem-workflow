package redeem

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

var ProcessRedeemFunc = service.ProcessRedeem

type RedeemRequest struct {
	Amount    float64 `json:"amount" binding:"required"`
	Recipient string  `json:"recipient" binding:"required"`
}

func HandleRedeemRequest(c *gin.Context) {
	var req RedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	request := models.Request{
		ID:        uuid.New(),
		Type:      "redeem",
		Amount:    req.Amount,
		Recipient: req.Recipient,
	}

	workflowInput := workflows.RedeemInput{
		Amount:    req.Amount,
		Recipient: req.Recipient,
		RequestID: request.ID.String(),
	}

	cadenceClient, err := deps.BuildCadenceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ProcessRedeemFunc(db.Db, &request, workflowInput, cadenceClient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": request.ID.String(), "status": "workflow started"})
}
