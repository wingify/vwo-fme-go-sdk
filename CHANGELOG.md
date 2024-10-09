# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2024-10-09

### Changed

- Changed the key from `userId` to `id` in the user context map.

## [1.0.0] - 2024-08-02

### Added

- First release of VWO Feature Management and Experimentation capabilities

	```go
	import vwo "github.com/wingify/vwo-fme-go-sdk"

	options := map[string]interface{}{
		"sdkKey": "your_sdk_key",
		"accountId": "your_account_id",
		"gatewayServiceURL": "your_gateway_sercice_url", // http://localhost:3000
	}

	instance, err := vwo.Init(options)

		// Correct JSON string with double quotes
	customVars := `{"key": "value"}`

	// Parse the JSON string into a Go map
	var customVariables map[string]interface{}
	json.Unmarshal([]byte(customVars), &customVariables)

		// Create the user context map
	userContext := map[string]interface{}{
		"id": "user_id",
		"customVariables": customVariables, // pass customVariables if using customVariables pre-segmentation
		"userAgent": "visitor_user_agent",
		"ipAddress": "visitor_ip_address",
	}

		// get flag to check if feature is Enabled for the user
		getFlag, err := instance.GetFlag("feature_key", userContext)

		isFeatureEnabled := getFlag.IsEnabled()
		getVariableValue := getFlag.GetVariable("variable_key", "default_value")

		// trackEvent to track the conversion for the user
		trackEventResponse, err := instance.TrackEvent("event_name", userContext, nil)

		// setAttribute to send attribute data for the user
		instance.SetAttribute("attribute_key", "attribute_value", userContext)
  ```
