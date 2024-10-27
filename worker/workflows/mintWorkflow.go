package workflows

import (
	"mint-redeem-workflow/activities"
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

type MintInput struct {
	Amount    float64
	Recipient string
	RequestID string
}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func MintWorkflow(ctx workflow.Context, amount float64, recipient string, requestID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("MintWorkflow started")
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var mintRes activities.MintActivityResponse
	mintRes.RequestId = requestID

	err := workflow.ExecuteActivity(ctx, activities.MintActivity, amount, recipient, requestID).Get(ctx, &mintRes)
	if err != nil {
		if err := workflow.ExecuteActivity(ctx, activities.UpdateStatusActivity, mintRes.RequestId, "failed").Get(ctx, &mintRes); err != nil {
			return err
		}
		return err
	} else {
		if err := workflow.ExecuteActivity(ctx, activities.UpdateStatusActivity, mintRes.RequestId, "completed").Get(ctx, &mintRes); err != nil {
			return err
		}
	}

	logger.Info("Workflow completed.", zap.String("Result", mintRes.RequestId))

	return nil
}
