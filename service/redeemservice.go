package service

import (
	"context"
	"mint-redeem-workflow/infra/cadence"
	"mint-redeem-workflow/models"
	"mint-redeem-workflow/worker/workflows"
	"time"

	"go.uber.org/cadence/client"
	"gorm.io/gorm"
)

func ProcessRedeem(db *gorm.DB, request *models.Request, workflowParam workflows.RedeemInput, cadenceClient cadence.WorkflowClient) error {
	request.Status = "pending"
	if err := db.Create(request).Error; err != nil {
		return err
	}
	
	workflowOptions := client.StartWorkflowOptions{
		ID:                           request.ID.String(),
		TaskList:                     "test-worker",
		ExecutionStartToCloseTimeout: time.Minute * 5,
	}

	workflowRun, err := cadenceClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.RedeemWorkflow, workflowParam.Amount, workflowParam.Recipient, request.ID.String())
	if err != nil {
		return err
	}

	request.Status = "started"
	request.RunID = workflowRun.GetRunID()

	if err := db.Save(request).Error; err != nil {
		return err
	}

	return nil
}
