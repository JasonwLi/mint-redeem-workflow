package mint

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

func mockProcessMint(db *gorm.DB, request *models.Request, workflowParam workflows.MintInput, cadenceClient cadence.WorkflowClient) error {
	return nil
}

func mockProcessMintError(db *gorm.DB, request *models.Request, workflowParam workflows.MintInput, cadenceClient cadence.WorkflowClient) error {
	return errors.New("mock process mint error")
}

func TestHandleMintRedeemRequest_SuccessReturns200(t *testing.T) {
	ProcessMintFunc = mockProcessMint
	defer func() { ProcessMintFunc = service.ProcessMint }()

	reqBody, _ := json.Marshal(MintRedeemRequest{Amount: 10.50, Recipient: "0xnotdeadbeef"})
	req, _ := http.NewRequest(http.MethodPost, "/mint", bytes.NewBuffer(reqBody))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	HandleMintRedeemRequest(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "workflow started", resp["status"])
}

func TestHandleMintRedeemRequest_MissingParamsReturns400(t *testing.T) {
	reqBody, _ := json.Marshal(MintRedeemRequest{Recipient: "0xnotdeadbeef"})
	req, _ := http.NewRequest(http.MethodPost, "/mint", bytes.NewBuffer(reqBody))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	HandleMintRedeemRequest(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request payload", resp["error"])
}

func TestHandleMintRedeemRequest_ProcessMintErrorReturns500(t *testing.T) {

	ProcessMintFunc = mockProcessMintError
	defer func() { ProcessMintFunc = service.ProcessMint }()

	reqBody, _ := json.Marshal(MintRedeemRequest{Amount: 10.50, Recipient: "0xnotdeadbeef"})
	req, _ := http.NewRequest(http.MethodPost, "/mint", bytes.NewBuffer(reqBody))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	HandleMintRedeemRequest(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "mock process mint error", resp["error"])
}
