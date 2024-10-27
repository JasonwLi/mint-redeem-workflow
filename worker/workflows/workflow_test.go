package workflows

import (
	"context"
	"errors"
	"mint-redeem-workflow/activities"
	"mint-redeem-workflow/db"
	"mint-redeem-workflow/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.uber.org/cadence"
	"go.uber.org/cadence/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func InitTestDB() {
	db.InitDB()
	db.Db.Exec("DELETE FROM requests")
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()

	s.env.RegisterActivity(activities.MintActivity)
	s.env.RegisterActivity(activities.UpdateStatusActivity)
	s.env.RegisterActivity(activities.RedeemActivity)

}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_MintWorkFlow_Success_RequestIsMarkedCompleted() {
	InitTestDB()
	request := models.Request{
		ID:        uuid.New(),
		Type:      "mint",
		Amount:    10.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	db.Db.Create(&request)

	s.env.ExecuteWorkflow(MintWorkflow, request.Amount, request.Recipient, request.ID.String())

	s.True(s.env.IsWorkflowCompleted())

	s.NoError(s.env.GetWorkflowError())

	var req models.Request
	db.Db.First(&req, "id = ?", request.ID)
	s.Equal("completed", req.Status)
}

func (s *UnitTestSuite) Test_MintWorkflow_ActivityParamPassedCorrectly() {
	InitTestDB()
	request := models.Request{
		ID:        uuid.New(),
		Type:      "mint",
		Amount:    10.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	db.Db.Create(&request)

	s.env.OnActivity(activities.MintActivity, mock.Anything, request.Amount, request.Recipient, request.ID.String()).Return(
		func(ctx context.Context, amount float64, recepient, requestID string) (activities.MintActivityResponse, error) {
			s.Equal(request.Recipient, recepient)
			s.Equal(request.ID.String(), requestID)
			return activities.MintActivityResponse{RequestId: requestID}, nil
		},
	)

	s.env.OnActivity(activities.UpdateStatusActivity, mock.Anything, request.ID.String(), "completed").Return(
		func(ctx context.Context, requestID, status string) error {
			s.Equal(request.ID.String(), requestID)
			return nil
		},
	)
	s.env.ExecuteWorkflow(MintWorkflow, request.Amount, request.Recipient, request.ID.String())

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_MintWorkflow_MintActivityFails_UpdatesRequestToFailed() {
	InitTestDB()
	request := models.Request{
		ID:        uuid.New(),
		Type:      "mint",
		Amount:    10.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	db.Db.Create(&request)

	s.env.OnActivity(activities.MintActivity, mock.Anything, request.Amount, request.Recipient, request.ID.String()).Return(
		func(ctx context.Context, amount float64, recepient, requestID string) (activities.MintActivityResponse, error) {
			s.Equal(request.Recipient, recepient)
			s.Equal(request.ID.String(), requestID)
			return activities.MintActivityResponse{RequestId: requestID}, errors.New("test error")
		},
	)

	s.env.ExecuteWorkflow(MintWorkflow, request.Amount, request.Recipient, request.ID.String())

	s.True(s.env.IsWorkflowCompleted())

	s.NotNil(s.env.GetWorkflowError())
	s.True(cadence.IsGenericError(s.env.GetWorkflowError()))
	s.Equal("test error", s.env.GetWorkflowError().Error())

	var req models.Request
	db.Db.First(&req, "id = ?", request.ID)
	s.Equal("failed", req.Status)
}

func (s *UnitTestSuite) Test_RedeemWorkflow_Success_RequestIsMarkedCompleted() {
	InitTestDB()
	request := models.Request{
		ID:        uuid.New(),
		Type:      "redeem",
		Amount:    10.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	db.Db.Create(&request)

	s.env.ExecuteWorkflow(RedeemWorkflow, request.Amount, request.Recipient, request.ID.String())

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var req models.Request
	db.Db.First(&req, "id = ?", request.ID)
	s.Equal("completed", req.Status)
}

func (s *UnitTestSuite) Test_RedeemWorkflow_ActivityParamPassedCorrectly() {
	InitTestDB()
	request := models.Request{
		ID:        uuid.New(),
		Type:      "redeem",
		Amount:    10.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	db.Db.Create(&request)

	s.env.OnActivity(activities.RedeemActivity, mock.Anything, request.Amount, request.Recipient, request.ID.String()).Return(
		func(ctx context.Context, amount float64, recipient, requestID string) (activities.RedeemActivityResponse, error) {
			s.Equal(request.Amount, amount)
			s.Equal(request.Recipient, recipient)
			s.Equal(request.ID.String(), requestID)
			return activities.RedeemActivityResponse{RequestId: requestID}, nil
		},
	)

	s.env.OnActivity(activities.UpdateStatusActivity, mock.Anything, request.ID.String(), "completed").Return(
		func(ctx context.Context, requestID, status string) error {
			s.Equal(request.ID.String(), requestID)
			s.Equal("completed", status)
			return nil
		},
	)
	s.env.ExecuteWorkflow(RedeemWorkflow, request.Amount, request.Recipient, request.ID.String())

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_RedeemWorkflow_RedeemActivityFails_UpdatesRequestToFailed() {
	InitTestDB()
	request := models.Request{
		ID:        uuid.New(),
		Type:      "redeem",
		Amount:    10.50,
		Recipient: "0xnotdeadbeef",
		Status:    "pending",
	}

	db.Db.Create(&request)

	s.env.OnActivity(activities.RedeemActivity, mock.Anything, request.Amount, request.Recipient, request.ID.String()).Return(
		func(ctx context.Context, amount float64, recipient, requestID string) (activities.RedeemActivityResponse, error) {
			s.Equal(request.Recipient, recipient)
			s.Equal(request.ID.String(), requestID)
			return activities.RedeemActivityResponse{RequestId: requestID}, errors.New("test error")
		},
	)

	s.env.ExecuteWorkflow(RedeemWorkflow, request.Amount, request.Recipient, request.ID.String())

	s.True(s.env.IsWorkflowCompleted())
	s.NotNil(s.env.GetWorkflowError())
	s.Equal("test error", s.env.GetWorkflowError().Error())

	var req models.Request
	db.Db.First(&req, "id = ?", request.ID)
	s.Equal("failed", req.Status)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
