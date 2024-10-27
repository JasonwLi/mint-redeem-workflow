package main

import (
	"fmt"
	"log"
	"net/http"

	"mint-redeem-workflow/activities"
	"mint-redeem-workflow/api/mint"
	"mint-redeem-workflow/api/redeem"
	"mint-redeem-workflow/db"
	"mint-redeem-workflow/deps"
	"mint-redeem-workflow/infra/cadence"
	"mint-redeem-workflow/worker/workflows"

	"github.com/gin-gonic/gin"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	db.InitDB()

	go startAPIServer()

	startCadenceWorker()
}

func startAPIServer() {
	r := gin.Default()

	r.POST("/mint", func(c *gin.Context) {
		mint.HandleMintRedeemRequest(c)
	})

	r.POST("/redeem", func(c *gin.Context) {
		redeem.HandleRedeemRequest(c)
	})

	if err := r.Run(":8090"); err != nil {
		log.Fatal("Failed to run API server:", err)
	}
}

func startCadenceWorker() {
	cadenceClient, err := deps.BuildCadenceServiceClient()
	if err != nil {
		fmt.Printf("Error creating Cadence client: %v\n", err)
		return
	}

	cadence.StartWorker("test-worker", "test-domain2", buildLogger(), cadenceClient)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start worker server:", err)
	}
}

func buildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		panic("Failed to set up logger")
	}

	return logger
}

func init() {
	workflow.Register(workflows.MintWorkflow)
	workflow.Register(workflows.RedeemWorkflow)
	activity.Register(activities.MintActivity)
	activity.Register(activities.RedeemActivity)
	activity.Register(activities.UpdateStatusActivity)
}
