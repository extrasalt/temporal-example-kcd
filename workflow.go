package main

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Order struct {
	Status string
}

var MainQueue = "COMMERCE_MAIN"

func CommerceWorkflow(ctx workflow.Context, customerID string, productID string) (Order, error) {

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    10 * time.Second,
		BackoffCoefficient: 2,
		MaximumInterval:    5 * time.Minute,
		MaximumAttempts:    5,
	}
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 180,
		StartToCloseTimeout:    time.Minute * 5,
		HeartbeatTimeout:       time.Minute * 5,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	ctx = workflow.WithTaskQueue(ctx, MainQueue)

	// Invoke inventory service and payment service simultaneously
	payment := workflow.ExecuteActivity(ctx, PaymentActivity, customerID)
	inv := workflow.ExecuteActivity(ctx, RemoveProductFromShelf, productID)

	var paymentInfo string

	err := payment.Get(ctx, &paymentInfo)
	if err != nil {
		var putRes string
		_ = workflow.ExecuteActivity(ctx, PutProductBackInShelf, productID).Get(ctx, &putRes)
		return Order{}, err
	}

	err = inv.Get(ctx, inv)
	if err != nil {
		return Order{}, err
	}

	var dispatchRes string
	err = workflow.ExecuteActivity(ctx, DispatchActivity, customerID, productID).Get(ctx, &dispatchRes)
	if err != nil {
		return Order{}, err
	}

	return Order{
		Status: "CONFIRMED",
	}, nil
}
