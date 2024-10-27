package service

import (
	"context"
	"errors"
	"mint-redeem-workflow/db"
	"mint-redeem-workflow/models"
	"mint-redeem-workflow/worker/workflows"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/cadence/client"
)

type MockCadenceClient struct {
	mock.Mock
}

func (m *MockCadenceClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	mockArgs := m.Called(ctx, options, workflow, args)
	return mockArgs.Get(0).(client.WorkflowRun), mockArgs.Error(1)
}

type MockWorkflowRun struct {
	mock.Mock
}

func (m *MockWorkflowRun) GetID() string {
	return "mock-workflow-id"
}

func (m *MockWorkflowRun) GetRunID() string {
	return "mock-run-id"
}

func (m *MockWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error {
	return nil
}

func InitTestDB() {
	db.InitDB()
	db.Db.Exec("DELETE FROM requests")
}

func TestProcessMint_Success_SavesRequestToDbUpdatesToStarted(t *testing.T) {
	InitTestDB()

	mockCadenceClient := new(MockCadenceClient)
	mockWorkflowRun := new(MockWorkflowRun)

	requestID := uuid.New()
	request := models.Request{
		ID:        requestID,
		Type:      "mint",
		Amount:    100.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}
	mockCadenceClient.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)

	workflowInput := workflows.MintInput{
		Amount:    100.50,
		Recipient: "0xnotdeadbeef",
		RequestID: requestID.String(),
	}

	err := ProcessMint(db.Db, &request, workflowInput, mockCadenceClient)
	assert.NoError(t, err)

	var dbRequest models.Request
	err = db.Db.First(&dbRequest, "id = ?", request.ID.String()).Error
	assert.NoError(t, err)
	assert.Equal(t, "started", dbRequest.Status)

	mockCadenceClient.AssertExpectations(t)
	mockWorkflowRun.AssertExpectations(t)
}

func TestProcessMint_WorkflowExecutionError_SavesRequestDoesNotUpdateStatus(t *testing.T) {
	InitTestDB()

	mockCadenceClient := new(MockCadenceClient)
	mockWorkflowRun := new(MockWorkflowRun)

	requestID := uuid.New()
	request := models.Request{
		ID:        requestID,
		Type:      "mint",
		Amount:    100.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	workflowInput := workflows.MintInput{
		Amount:    100.50,
		Recipient: "0xnotdeadbeef",
		RequestID: requestID.String(),
	}

	mockCadenceClient.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockWorkflowRun, errors.New("workflow execution error"))

	err := ProcessMint(db.Db, &request, workflowInput, mockCadenceClient)
	assert.EqualError(t, err, "workflow execution error")

	var dbRequest models.Request
	err = db.Db.First(&dbRequest, "id = ?", request.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "pending", dbRequest.Status)

	mockCadenceClient.AssertExpectations(t)

}

func TestProcessRedeem_SavesRequestToDbUpdatesToStarted(t *testing.T) {
	InitTestDB()

	mockCadenceClient := new(MockCadenceClient)
	mockWorkflowRun := new(MockWorkflowRun)

	requestID := uuid.New()
	request := models.Request{
		ID:        requestID,
		Type:      "redeem",
		Amount:    100.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	mockCadenceClient.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockWorkflowRun, nil)

	workflowInput := workflows.RedeemInput{
		Amount:    request.Amount,
		Recipient: request.Recipient,
		RequestID: request.ID.String(),
	}

	err := ProcessRedeem(db.Db, &request, workflowInput, mockCadenceClient)
	assert.NoError(t, err)

	var dbRequest models.Request
	err = db.Db.First(&dbRequest, "id = ?", request.ID.String()).Error
	assert.NoError(t, err)
	assert.Equal(t, "started", dbRequest.Status)

	mockCadenceClient.AssertExpectations(t)
	mockWorkflowRun.AssertExpectations(t)
}

func TestProcessRedeem_WorkflowExecutionError(t *testing.T) {
	InitTestDB()

	mockCadenceClient := new(MockCadenceClient)
	mockWorkflowRun := new(MockWorkflowRun)

	requestID := uuid.New()
	request := models.Request{
		ID:        requestID,
		Type:      "redeem",
		Amount:    10.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	workflowInput := workflows.RedeemInput{
		Amount:    request.Amount,
		Recipient: request.Recipient,
		RequestID: request.ID.String(),
	}

	mockCadenceClient.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockWorkflowRun, errors.New("workflow execution error"))

	err := ProcessRedeem(db.Db, &request, workflowInput, mockCadenceClient)
	assert.EqualError(t, err, "workflow execution error")

	var dbRequest models.Request
	err = db.Db.First(&dbRequest, "id = ?", request.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "pending", dbRequest.Status)

	mockCadenceClient.AssertExpectations(t)
}
