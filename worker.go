package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func BuildTemporalClient() client.Client {
	domainName := os.Getenv("TEMPORAL_NAMESPACE")
	hostNameAndPort := os.Getenv("TEMPORAL_HOST")

	if domainName == "" || hostNameAndPort == "" {
		log.Println("Couldn't override helper builder because of missing environment variables")
	}

	c, err := client.NewClient(client.Options{
		Namespace: domainName,
		HostPort:  hostNameAndPort,
	})
	if err != nil {
		log.Println("Error building temporal client:", err.Error())
	}

	return c
}

func mainWorker() {
	service := BuildTemporalClient()

	workerOptions := worker.Options{
		MaxConcurrentActivityTaskPollers: 2,
	}

	w := worker.New(
		service,
		MainQueue,
		workerOptions,
	)

	w.RegisterWorkflow(CommerceWorkflow)
	w.RegisterActivity(PaymentActivity)
	w.RegisterActivity(RemoveProductFromShelf)
	w.RegisterActivity(PutProductBackInShelf)
	w.RegisterActivity(DispatchActivity)

	if err := w.Start(); err != nil {
		panic(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	log.Fatal(http.ListenAndServe(":8001", nil))

}
