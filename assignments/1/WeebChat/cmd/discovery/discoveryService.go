package main

import (
	"WeebChat/pkg/services/discovery"
	"fmt"
	"os"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	fmt.Printf("Setting up discovery service on %s:%s\n", host, port)

	discoveryService := discovery.NewDiscoveryServiceServer(host, port)

	discoveryService.Setup()
	err := discoveryService.Start()

	if err != nil {
		fmt.Println("Error starting up server")
		os.Exit(-1)
	}

	fmt.Printf("Running discovery service on %s:%s\n", host, port)

}
