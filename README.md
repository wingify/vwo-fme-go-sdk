# VWO Feature Management and Experimentation SDK for Golang

![Size in Bytes](https://img.shields.io/github/languages/code-size/wingify/vwo-fme-go-sdk)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

This open source library allows you to A/B Test your Website at server-side.

## Requirements

- Works with Go 1.16+

## Installation

```go
go get "github.com/wingify/vwo-fme-go-sdk"
```

## Basic usage

```go
import vwo "github.com/wingify/vwo-fme-go-sdk"

    options := map[string]interface{}{
		"sdkKey":            "your_sdk_key",
		"accountId":         "your_account_id",
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
		"userId":          "user_id",
		"customVariables": customVariables, // pass customVariables if using customVariables pre-segmentation
		"userAgent":       "visitor_user_agent",
		"ipAddress":       "visitor_ip_address",
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

## Contributing

Please go through our [contributing guidelines](CONTRIBUTING.md)

## Code of Conduct

[Code of Conduct](CODE_OF_CONDUCT.md)

## License

[Apache License, Version 2.0](LICENSE)