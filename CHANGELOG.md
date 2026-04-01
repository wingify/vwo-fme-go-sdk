# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.7.1] - 2026-04-01

### Fixed

- Removed validation check for zero variation weight, allowing variations with a custom weight of `0` to be correctly processed without errors.

## [1.7.0] - 2026-03-11

### Added

- Added support for custom bucketing seed via `bucketingSeed` in the context. This allows users to bucket by a shared identifier instead of the individual user ID, ensuring all users within the same group receive the same variation.

	Example usage:

	```go
	vwoInstance, err := vwo.Init(map[string]interface{}{
		"accountId":         123456,
		"sdkKey":            "32-alpha-numeric-sdk-key",
	})
	if err != nil {
		log.Fatalf("Failed to initialize VWO client: %v", err)
	}

	// Use bucketing seed
	context := map[string]interface{}{
		"id":            "new_user_id",
		"bucketingSeed": "some-seed-xyz",
	}
	
	flag, err := vwoInstance.GetFlag("feature-key", context)
	```

## [1.6.0] - 2026-02-25

### Added

- Added support to use the context `id` as the visitor UUID instead of auto-generating one. You can read the visitor UUID from the flag result via `flag.GetUUID()` (e.g. to pass to the web client).

	Example usage:

	```go
	vwoInstance, err := vwo.Init(map[string]interface{}{
		"accountId": 123456,
		"sdkKey":   "32-alpha-numeric-sdk-key",
	})
	if err != nil {
		log.Fatalf("Failed to initialize VWO client: %v", err)
	}

	// Default: SDK generates a UUID from id and account
	contextWithGeneratedUuid := map[string]interface{}{"id": "user-123"}
	flag1, err := vwoInstance.GetFlag("feature-key", contextWithGeneratedUuid)
	// Get the UUID from the flag result (e.g. to pass to web client)
	uuid := flag1.GetUUID()
	fmt.Println("Visitor UUID:", uuid)

	// Use your own UUID (e.g. from web client) by passing a valid web UUID in context.id
	contextWithCustomUuid := map[string]interface{}{
		"id": "D7E2EAA667909A2DB8A6371FF0975C2A5", // your existing UUID
	}
	flag2, err := vwoInstance.GetFlag("feature-key", contextWithCustomUuid)
	```

## [1.5.0] - 2026-01-23

### Added

- Added support for redirecting all network calls through a custom proxy URL. This feature allows users to route all SDK network requests (settings, tracking, etc.) through their own proxy server.
	```go
	options := map[string]interface{}{
		"sdkKey":    "32-alpha-numeric-sdk-key", // Replace with your SDK key
		"accountId": "123456",                   // Replace with your account ID
		"proxyUrl": "https://custom.proxy.com"   // Replace with your custom proxy url
	}

	// Initialize VWO instance
	vwoInstance, err := vwo.Init(options)
	```
	**Note**: If both gateway_service and proxy_url are provided, the SDK will give preference to the gateway_service for all network requests.

## [1.4.0] - 2025-12-12

### Fixed

- Send `vwo_sdkUsageStats` event when `usageStatsAccountId` is present in settings.

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
