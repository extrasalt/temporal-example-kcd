package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
		fmt.Println("worker")

	case "server":
		fmt.Println("server")
	}
}
