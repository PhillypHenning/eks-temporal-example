package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ApplicationName is the identifier for a specific application's task list
const ApplicationName = "sampleGroup"

// SimpleActivity is a very basic activity that just returns a message
func SimpleActivity(ctx context.Context, message string) (string, error) {
	return fmt.Sprintf("Activity processed message: %s", message), nil
}

// SampleWorkflow is a basic workflow that calls the SimpleActivity
func SampleWorkflow(ctx workflow.Context, message string) (string, error) {
	// Activity options setting
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Duration(10 * time.Second),
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err := workflow.ExecuteActivity(ctx, SimpleActivity, message).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func main() {
	// Create a Temporal client
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// Create a worker for the given task list
	w := worker.New(c, ApplicationName, worker.Options{})

	// Register the workflow and activities with the worker
	w.RegisterWorkflow(SampleWorkflow)
	w.RegisterActivity(SimpleActivity)

	// Start the worker to poll and process tasks
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
