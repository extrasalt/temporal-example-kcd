package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

func main() {
	var mode string

	flag.StringVar(&mode, "m", "server", "Mode is worker or trigger.")
	flag.Parse()

	temporalNS := os.Getenv("TEMPORAL_NAMESPACE")
	if temporalNS == "" {
		log.Println("Namespace missing")
	}

	switch mode {
	case "worker":
		mainWorker()

	case "trigger":
		temporalClient := BuildTemporalClient()

		customerID := "mohan"
		productID := "123"
		workflowOptions := client.StartWorkflowOptions{
			ID:                       fmt.Sprintf("%s-%s", customerID, productID),
			TaskQueue:                MainQueue,
			WorkflowExecutionTimeout: time.Minute * 180,
			WorkflowTaskTimeout:      time.Second * 60,
			WorkflowIDReusePolicy:    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			// CronSchedule:             "0 0 * * *",
		}

		temporalResponse, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, CommerceWorkflow, customerID, productID)
		if err != nil {
			log.Println(fmt.Sprintf("Couldn't Schedule Workflow. Error: %s", err))

			return
		}

		fmt.Printf("%+v", temporalResponse)
	}
}
