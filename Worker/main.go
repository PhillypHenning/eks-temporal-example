package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const ApplicationName = "sampleGroup"

func SimpleActivity(ctx context.Context, message string) (string, error) {
	return fmt.Sprintf("Activity processed message (B): %s", message), nil
}

// func SampleWorkflow(ctx workflow.Context, message string) (string, error) {
// 	options := workflow.ActivityOptions{
// 		StartToCloseTimeout: time.Duration(10 * time.Second),
// 	}
// 	ctx = workflow.WithActivityOptions(ctx, options)
// 	var result string
// 	err := workflow.ExecuteActivity(ctx, SimpleActivity, message).Get(ctx, &result)
// 	if err != nil {
// 		return "", err
// 	}
// 	return result, nil
// }

func SampleWorkflow(ctx workflow.Context, message string) (string, error) {
	// Set search attributes
	searchAttributes := map[string]interface{}{
		"CustomAttribute1": os.Getenv("CANARY_DEPLOYMENT_ATTRIBUTE"),
	}

	err := workflow.UpsertSearchAttributes(ctx, searchAttributes)
	if err != nil {
		return "", err
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Duration(10 * time.Second),
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err = workflow.ExecuteActivity(ctx, SimpleActivity, message).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func main() {
	temporalServer := os.Getenv("TEMPORAL_SERVER")
	if temporalServer == "" {
		log.Fatalln("TEMPORAL_SERVER environment variable not set")
	}

	// Create a Temporal client
	c, err := client.NewClient(client.Options{
		HostPort: temporalServer,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, ApplicationName, worker.Options{})

	w.RegisterWorkflow(SampleWorkflow)
	w.RegisterActivity(SimpleActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
