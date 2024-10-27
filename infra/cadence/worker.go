package cadence

import (
	"github.com/uber-go/tally"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

func StartWorker(taskName string, domain string, logger *zap.Logger, service workflowserviceclient.Interface) {
	workerOptions := worker.Options{
		Logger:       logger,
		MetricsScope: tally.NewTestScope(taskName, map[string]string{}),
	}

	worker := worker.New(
		service,
		domain,
		taskName,
		workerOptions)
	err := worker.Start()
	if err != nil {
		panic("Failed to start worker")
	}

	logger.Info("Started Worker.", zap.String("worker", taskName))
}
