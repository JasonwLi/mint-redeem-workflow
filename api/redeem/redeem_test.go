package redeem

import (
	"bytes"
	"encoding/json"
	"errors"
	"mint-redeem-workflow/infra/cadence"
	"mint-redeem-workflow/models"
	"mint-redeem-workflow/service"
	"mint-redeem-workflow/worker/workflows"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func mockProcessRedeem(db *gorm.DB, request *models.Request, workflowParam workflows.RedeemInput, cadenceClient cadence.WorkflowClient) error {
	return nil
}

func mockProcessRedeemError(db *gorm.DB, request *models.Request, workflowParam workflows.RedeemInput, cadenceClient cadence.WorkflowClient) error {
	return errors.New("mock process redeem error")
}

func TestHandleRedeemRequest_SuccessReturns200(t *testing.T) {
	ProcessRedeemFunc = mockProcessRedeem
	defer func() { ProcessRedeemFunc = service.ProcessRedeem }()

	reqBody, _ := json.Marshal(RedeemRequest{Amount: 10.50, Recipient: "0xnotdeadbeef"})
	req, _ := http.NewRequest(http.MethodPost, "/redeem", bytes.NewBuffer(reqBody))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	HandleRedeemRequest(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "workflow started", resp["status"])
}

func TestHandleRedeemRequest_MissingParamsReturns400(t *testing.T) {
	reqBody, _ := json.Marshal(RedeemRequest{Recipient: "0xnotdeadbeef"})
	req, _ := http.NewRequest(http.MethodPost, "/redeem", bytes.NewBuffer(reqBody))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	HandleRedeemRequest(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request payload", resp["error"])
}

func TestHandleRedeemRequest_ProcessRedeemErrorReturns500(t *testing.T) {
	ProcessRedeemFunc = mockProcessRedeemError
	defer func() { ProcessRedeemFunc = service.ProcessRedeem }()

	reqBody, _ := json.Marshal(RedeemRequest{Amount: 10.50, Recipient: "0xnotdeadbeef"})
	req, _ := http.NewRequest(http.MethodPost, "/redeem", bytes.NewBuffer(reqBody))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	HandleRedeemRequest(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "mock process redeem error", resp["error"])
}
