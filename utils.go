package main

import (
	"fmt"
	"os"
)

var XApiKey = "unset"

func getXApiKey() string {
	if XApiKey == "unset" {
		XApiKey = os.Getenv("X_API_KEY")
	}
	if XApiKey == "" {
		fmt.Println("Warning: X_API_KEY environment variable is not set.")
		panic("X_API_KEY environment variable is required")
	}
	return XApiKey
}
