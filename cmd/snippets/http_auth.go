package main

import (
	"flag"
	"fmt"
	"time"
)

// The config struct holds all configuration settings for the application.
type config struct {
	port           int
	verboseLogging bool
	requestTimeout time.Duration
	basicAuth      struct {
		username string
		password string
	}
}

func main() {
	// Create a new config instance.
	var cfg config

	// Define the command-line flags. Notice that we define these so that the values
	// are read directly into the appropriate config struct field, and set sensible default
	// values for each of them.
	flag.IntVar(&cfg.port, "port", 4000, "The port number the web application listens on")
	flag.BoolVar(&cfg.verboseLogging, "verbose-logging", false, "Enables detailed request and error logging")
	flag.DurationVar(&cfg.requestTimeout, "request-timeout", 5*time.Second, "Maximum duration to wait for a request to complete")
	flag.StringVar(&cfg.basicAuth.username, "basic-auth-username", "", "Username required for HTTP Basic Authentication")
	flag.StringVar(&cfg.basicAuth.password, "basic-auth-password", "", "Password required for HTTP Basic Authentication")

	// Parse the flags with the flag.Parse function. This is important!
	flag.Parse()

	// Print all configuration settings.
	fmt.Printf("Port: %d\n", cfg.port)
	fmt.Printf("Verbose Logging: %t\n", cfg.verboseLogging)
	fmt.Printf("Request Timeout: %v\n", cfg.requestTimeout)
	fmt.Printf("Basic Auth Username: %s\n", cfg.basicAuth.username)
	fmt.Printf("Basic Auth Password: %s\n", cfg.basicAuth.password)
}