package main

import (
	"context"
	"net/http"
)

func PaymentActivity(ctx context.Context, customerID string) (string, error) {
	req, err := http.NewRequest("POST", "https://api.paymentgateway.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer yourkeyorsomething")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return "SUCCESS", nil
}

func RemoveProductFromShelf(ctx context.Context, productID string) (string, error) {
	req, err := http.NewRequest("POST", "https://inventory.service", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer yourkeyorsomething")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return "SUCCESS", nil
}
