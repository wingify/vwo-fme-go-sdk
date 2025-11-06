# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2025-11-06

### Added

- First release of VWO Feature Management and Experimentation capabilities

	```go
	package main

	import (
		"fmt"
		"log"

		vwo "github.com/wingify/vwo-fme-go-sdk"
	)

	func main() {
		// Initialize VWO SDK with your account details
		options := map[string]interface{}{
			"sdkKey":    "32-alpha-numeric-sdk-key", // Replace with your SDK key
			"accountId": "123456",                   // Replace with your account ID
		}

		// Initialize VWO instance
		vwoInstance, err := vwo.Init(options)
		if err != nil {
			log.Fatalf("Failed to initialize VWO client: %v", err)
		}

		// Create user context
		context := map[string]interface{}{
			"id": "unique_user_id", // Set a unique user identifier
		}

		// Check if a feature flag is enabled
		getFlag, err := vwoInstance.GetFlag("feature_key", context)
		if err != nil {
			log.Printf("Error getting feature flag: %v", err)
		} else {
			isFeatureEnabled := getFlag.IsEnabled()
			fmt.Println("Is feature enabled?", isFeatureEnabled)

			// Get a variable value with a default fallback
			variableValue := getFlag.GetVariable("feature_variable", "default_value")
			fmt.Println("Variable value:", variableValue)
		}

		// Track a custom event
		trackResponse, err := vwoInstance.TrackEvent("event_name", context)
		if err != nil {
			log.Printf("Error tracking event: %v", err)
		} else {
			fmt.Println("Event tracked:", trackResponse)
		}

		// Set multiple custom attributes
		attributeMap := map[string]interface{}{
			"attribute-name": "attribute-value",
		}
		err = vwoInstance.SetAttribute(attributeMap, context)
		if err != nil {
			log.Printf("Error setting attributes: %v", err)
		}
	}
	```
