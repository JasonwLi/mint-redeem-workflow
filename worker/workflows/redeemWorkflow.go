package workflows

import (
	"mint-redeem-workflow/activities"
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

type RedeemInput struct {
	Amount    float64
	Recipient string
	RequestID string
}

var redeemActivityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func RedeemWorkflow(ctx workflow.Context, amount float64, recipient string, requestID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("RedeemWorkflow started")
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var redeemRes activities.RedeemActivityResponse
	redeemRes.RequestId = requestID

	if err := workflow.ExecuteActivity(ctx, activities.RedeemActivity, amount, recipient, requestID).Get(ctx, &redeemRes); err != nil {
		if err := workflow.ExecuteActivity(ctx, activities.UpdateStatusActivity, redeemRes.RequestId, "failed").Get(ctx, &redeemRes); err != nil {
			return err
		}
		return err
	} else {
		if err := workflow.ExecuteActivity(ctx, activities.UpdateStatusActivity, redeemRes.RequestId, "completed").Get(ctx, &redeemRes); err != nil {
			return err
		}
	}

	logger.Info("Workflow completed.", zap.String("Result", redeemRes.RequestId))

	return nil
}
