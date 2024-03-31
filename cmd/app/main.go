package main

import (
	"fmt"
	"simactive/internal/config"
)

func main() {

	// Initialize config object
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// Initialize logger

	// Initialize application

	// Run gRPC server
}
