package main

import "go.temporal.io/sdk/workflow"

type Order struct {
	Status string
}

func CommerceWorkflow(ctx workflow.Context, customerID string, productID string) (Order, error) {
	// Invoke inventory service and payment service simultaneously
	payment := workflow.ExecuteActivity(ctx, PaymentActivity, customerID)
	inv := workflow.ExecuteActivity(ctx, RemoveProductFromShelf, productID)

	var paymentInfo string

	err := payment.Get(ctx, &paymentInfo)
	if err != nil {
		return Order{}, err
	}

	err = inv.Get(ctx, inv)
	if err != nil {
		return Order{}, err
	}

	return Order{
		Status: "CONFIRMED",
	}, nil

}
