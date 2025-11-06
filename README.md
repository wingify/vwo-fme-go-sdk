# VWO Feature Management and Experimentation SDK for Go

[![CI](https://github.com/wingify/vwo-fme-go-sdk/workflows/CI/badge.svg?branch=master)](https://github.com/wingify/vwo-fme-go-sdk/actions?query=workflow%3ACI)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

## Overview

The **VWO Feature Management and Experimentation SDK** (VWO FME Go SDK) enables Go developers to integrate feature flagging and experimentation into their applications. This SDK provides full control over feature rollout, A/B testing, and event tracking, allowing teams to manage features dynamically and gain insights into user behavior.

## Requirements

The Go SDK supports:

* Go 1.16 or higher

Our [Build](https://github.com/wingify/vwo-fme-go-sdk/actions) is successful on these Go Versions -

## Installation

```bash
go get github.com/wingify/vwo-fme-go-sdk
```

## Basic Usage Example

The following example demonstrates initializing the SDK with a VWO account ID and SDK key, setting a user context, checking if a feature flag is enabled, and tracking a custom event.

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
    trackResponse, err := vwoInstance.TrackEvent("event_name", context, nil)
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

## Advanced Configuration Options

To customize the SDK further, additional parameters can be passed to the `Init()` API using the options map. Here's a table describing each option:

| **Parameter**                | **Description**                                                                                                                                             | **Required** | **Type** | **Example**                     |
| ---------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------ | -------- | ------------------------------- |
| `accountId`                  | VWO Account ID for authentication.                                                                                                                          | Yes          | String/Int   | `"123456"` or `123456`                      |
| `sdkKey`                     | SDK key corresponding to the specific environment to initialize the VWO SDK Client. You can get this key from VWO Application.                              | Yes          | String   | `"32-alpha-numeric-sdk-key"`    |
| `pollInterval`               | Time interval for fetching updates from VWO servers (in milliseconds).                                                                                      | No           | Number   | `60000`                         |
| `gatewayService`             | Configuration for integrating VWO Gateway Service. Service.                                                                                   | No           | Object   | see [Gateway](#gateway) section |
| `storage`                    | Custom storage connector for persisting user decisions and campaign data. data.                                                                                   | No           | Object   | See [Storage](#storage) section |
| `logger`                     | Toggle log levels for more insights or for debugging purposes. You can also customize your own transport in order to have better control over log messages. | No           | Object   | See [Logger](#logger) section   |
| `integrations`               | Callback function for integrating with third-party analytics services.                                                                                      | No           | Function | See [Integrations](#integrations) section |
| `retryConfig`                | Configuration for network request retry behavior and exponential backoff strategy                                                                           | No           | Object   | See [Retry Config](#retry-config) section |

Refer to the [official VWO documentation](https://developers.vwo.com/v2/docs/fme-go-install) for additional parameter details.

### User Context

The user context is a `map[string]interface{}` that uniquely identifies users and is crucial for consistent feature rollouts. A typical context includes an `id` for identifying the user. It can also include other attributes that can be used for targeting and segmentation, such as custom variables, user agent, and IP address.

#### Parameters Table

The following table explains all the parameters in the context map:

| **Parameter** | **Description** | **Required** | **Type** |
|---------------|-----------------|--------------|----------|
| `id` | Unique identifier for the user. | Yes | string |
| `customVariables` | Custom attributes for targeting. | No | map[string]interface{} |
| `userAgent` | User agent string for identifying the user's browser and operating system. | No | string |
| `ipAddress` | IP address of the user. | No | string |

#### Example

```go
context := map[string]interface{}{
    "id": "unique_user_id", // Set a unique user identifier

    // Create custom variables
    "customVariables": map[string]interface{}{
        "age":      25,
        "location": "US",
    },

    "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
    "ipAddress": "1.1.1.1",
}
```

### Basic Feature Flagging

Feature Flags serve as the foundation for all testing, personalization, and rollout rules within FME.
To implement a feature flag, first use the `GetFlag()` method to retrieve the flag configuration.
The `GetFlag()` method provides a simple way to check if a feature is enabled for a specific user and access its variables. It returns a `GetFlag` object that contains methods like `IsEnabled()` for checking the feature's status and `GetVariable()` for retrieving any associated variables.

| Parameter | Description | Required | Type |
|-----------|-------------|----------|------|
| `featureKey` | Unique identifier of the feature flag | Yes | string |
| `context` | Map containing user identification and contextual information | Yes | map[string]interface{} |

Example usage:

```go
featureFlag, err := vwoInstance.GetFlag("feature_key", context)
if err != nil {
    log.Printf("Error getting feature flag: %v", err)
    return
}

isEnabled := featureFlag.IsEnabled()

if isEnabled {
    fmt.Println("Feature is enabled!")

    // Get and use feature variable with type safety
    variableValue := featureFlag.GetVariable("feature_variable", "default_value")
    fmt.Println("Variable value:", variableValue)
} else {
    fmt.Println("Feature is not enabled!")
}
```

### Custom Event Tracking

Feature flags can be enhanced with connected metrics to track key performance indicators (KPIs) for your features. These metrics help measure the effectiveness of your testing rules by comparing control versus variation performance, and evaluate the impact of personalization and rollout campaigns. Use the `TrackEvent()` method to track custom events like conversions, user interactions, and other important metrics:

| Parameter | Description | Required | Type |
|-----------|-------------|----------|------|
| `eventName` | Name of the event you want to track | Yes | string |
| `context` | Map containing user identification and contextual information | Yes | map[string]interface{} |
| `eventProperties` | Additional properties/metadata associated with the event | No | map[string]interface{} |

Example usage:

```go
eventProperties := map[string]interface{}{
    "revenue":  100.50,
    "currency": "USD",
}

trackResponse, err := vwoInstance.TrackEvent("purchase", context, eventProperties)
if err != nil {
    log.Printf("Error tracking event: %v", err)
} else {
    fmt.Println("Event tracked:", trackResponse)
}
```

See [Tracking Conversions](https://developers.vwo.com/v2/docs/fme-go-metrics#usage) documentation for more information.

### Pushing Attributes

User attributes provide rich contextual information about users, enabling powerful personalization. The `SetAttribute()` method in VWO provides a simple way to associate these attributes with users in VWO for advanced segmentation. The method accepts an attribute map and context map containing the user information.

| Parameter | Description | Required | Type |
|-----------|-------------|----------|------|
| `attributeMap` | Multiple attributes you want to set for a user. | Yes | map[string]interface{} |
| `context` | Map containing user identification and other contextual information | Yes | map[string]interface{} |

Example usage:

```go
attributeMap := map[string]interface{}{
    "plan":        "premium",
    "trial_used":  true,
    "signup_date": "2025-01-01",
}

err := vwoInstance.SetAttribute(attributeMap, context)
if err != nil {
    log.Printf("Error setting attributes: %v", err)
}
```

See [Pushing Attributes](https://developers.vwo.com/v2/docs/fme-go-attributes#usage) documentation for additional information.

### Polling Interval Adjustment

The `pollInterval` is an optional parameter that allows the SDK to automatically fetch and update settings from the VWO server at specified intervals. The polling interval can be configured in three ways:

1. Set via SDK options: If `pollInterval` is specified in the initialization options (must be >= 1000 milliseconds), that interval will be used
2. VWO Application Settings: If configured in your VWO application settings, that interval will be used
3. Default Fallback: If neither of the above is set, a 10 minute (600,000 milliseconds) polling interval is used

Setting this parameter ensures your application always uses the latest configuration by periodically checking for and applying any updates.

```go
options := map[string]interface{}{
    "sdkKey":       "32-alpha-numeric-sdk-key",
    "accountId":    "123456",
    "pollInterval": 60000, // Set the poll interval to 60 seconds
}

vwoInstance, err := vwo.Init(options)
```

### Logger

VWO by default logs all `ERROR` level messages to your server console.
To gain more control over VWO's logging behaviour, you can use the `logger` parameter in the `Init` configuration.

| **Parameter** | **Description** | **Required** | **Type** | **Default Value** |
|---------------|-----------------|--------------|----------|-------------------|
| `level` | Log level to control verbosity of logs | Yes | string | `"ERROR"` |
| `prefix` | Custom prefix for log messages | No | string | `"VWO-SDK"` |

#### Example 1: Set log level to control verbosity of logs

```go
options := map[string]interface{}{
    "sdkKey":    "32-alpha-numeric-sdk-key",
    "accountId": "123456",
    "logger": map[string]interface{}{
        "level": "DEBUG",
    },
}

vwoInstance, err := vwo.Init(options)
```

#### Example 2: Add custom prefix to log messages for easier identification

```go
options := map[string]interface{}{
    "sdkKey":    "32-alpha-numeric-sdk-key",
    "accountId": "123456",
    "logger": map[string]interface{}{
        "level":  "DEBUG",
        "prefix": "CUSTOM LOG PREFIX",
    },
}

vwoInstance, err := vwo.Init(options)
```

### Gateway

The VWO FME Gateway Service is an optional but powerful component that enhances VWO's Feature Management and Experimentation (FME) SDKs. It acts as a critical intermediary for pre-segmentation capabilities based on user location and user agent (UA). By deploying this service within your infrastructure, you benefit from minimal latency and strengthened security for all FME operations.

#### Why Use a Gateway?

The Gateway Service is required in the following scenarios:

- When using pre-segmentation features based on user location or user agent.
- For applications requiring advanced targeting capabilities.
- It's mandatory when using any thin-client SDK (e.g., Go).

#### How to Use the Gateway

The gateway can be customized by passing the `gatewayService` parameter in the `Init` configuration.

```go
options := map[string]interface{}{
    "sdkKey":    "32-alpha-numeric-sdk-key",
    "accountId": "123456",
    "gatewayService": map[string]interface{}{
        "url": "http://custom.gateway.com",
    },
}

vwoInstance, err := vwo.Init(options)
```

Refer to the [Gateway Documentation](https://developers.vwo.com/v2/docs/gateway-service) for further details.

### Storage

The SDK operates in a stateless mode by default, meaning each `GetFlag` call triggers a fresh evaluation of the flag against the current user context.

To optimize performance and maintain consistency, you can implement a custom storage mechanism by passing a `storage` parameter during initialization. This allows you to persist feature flag decisions in your preferred database system (like Redis, MongoDB, or any other data store).

Key benefits of implementing storage:

- Improved performance by caching decisions
- Consistent user experience across sessions
- Reduced load on your application

The storage mechanism ensures that once a decision is made for a user, it remains consistent even if campaign settings are modified in the VWO Application. This is particularly useful for maintaining a stable user experience during A/B tests and feature rollouts.

```go
// CustomStorageConnector implements the storage.Connector interface
type CustomStorageConnector struct {
	data map[string]map[string]interface{}
}

// NewCustomStorageConnector creates a new custom storage connector
func NewCustomStorageConnector() *CustomStorageConnector {
	return &CustomStorageConnector{
		data: make(map[string]map[string]interface{}),
	}
}

// Set stores data in the custom storage
func (c *CustomStorageConnector) Set(data map[string]interface{}) error {
    // example implementation of SET
	featureKey, _ := data["featureKey"].(string)
	userID, _ := data["userId"].(string)

	key := featureKey + ":" + userID
	c.data[key] = data
	return nil
}

// Get retrieves data from the custom storage
func (c *CustomStorageConnector) Get(featureKey string, userID string) (interface{}, error) {
     // example implementation of GET
	key := featureKey + ":" + userID
	if data, exists := c.data[key]; exists {
		return data, nil
	}
	return nil, nil
}

// initialise the storage
customStorage := NewCustomStorageConnector()

// Use in initialization
options := map[string]interface{}{
    "sdkKey":    "32-alpha-numeric-sdk-key",
    "accountId": "123456",
    "storage":   customStorage
}

vwoInstance, err := vwo.Init(options)
```

### Integrations

VWO FME SDKs provide seamless integration with third-party tools like analytics platforms, monitoring services, customer data platforms (CDPs), and messaging systems. This is achieved through a simple yet powerful callback mechanism that receives VWO-specific properties and can forward them to any third-party tool of your choice.

```go
options := map[string]interface{}{
    "sdkKey":       "32-alpha-numeric-sdk-key",
    "accountId":    "123456",
    "integrations": map[string]interface{}{
        "Callback": func(properties map[string]interface{}) {
            // implement your custom logic here
            fmt.Printf("Integration callback called with properties: %+v\n", properties)
        },
    },
}

vwoInstance, err := vwo.Init(options)
```

Refer to the [Integrations](https://developers.vwo.com/v2/docs/fme-go-integrations) documentation for more information.

### Retry Config

The `retryConfig` parameter allows you to customize the retry behavior for network requests. This is particularly useful for applications that need to handle network failures gracefully with exponential backoff strategies.

| **Parameter**       | **Description**                                           | **Required** | **Type** | **Default** | **Validation**                      |
| ------------------- | --------------------------------------------------------- | ------------ | -------- | ----------- | ----------------------------------- |
| `shouldRetry`       | Whether to enable automatic retry on network failures     | No           | Boolean  | `true`      | Must be a boolean value             |
| `maxRetries`        | Maximum number of retry attempts before giving up         | No           | Number   | `3`         | Must be a non-negative integer >= 1 |
| `initialDelay`      | Initial delay (in seconds) before the first retry attempt | No           | Number   | `2`         | Must be a non-negative integer >= 1 |
| `backoffMultiplier` | Multiplier for exponential backoff between retry attempts | No           | Number   | `2`         | Must be a non-negative integer >= 2 |

#### How Retry Logic Works

The SDK implements an exponential backoff strategy for retrying failed network requests:

1. **Initial Request**: The SDK attempts the initial network request
2. **On Failure**: If the request fails and `shouldRetry` is `true`, the SDK waits for `initialDelay` seconds
3. **Exponential Backoff**: For subsequent retries, the delay is calculated as: `initialDelay × (backoffMultiplier ^ attempt)`
4. **Maximum Attempts**: The SDK will retry up to `maxRetries` times before giving up

#### Example Usage

```go
options := map[string]interface{}{
    "sdkKey":    "32-alpha-numeric-sdk-key",
    "accountId": "123456",
    "retryConfig": map[string]interface{}{
        "shouldRetry":       true,  // Enable retries
        "maxRetries":        5,     // Retry up to 5 times
        "initialDelay":      3,     // Wait 3 seconds before first retry
        "backoffMultiplier": 2,     // Double the delay for each subsequent retry
    },
}

vwoInstance, err := vwo.Init(options)
```

With this configuration, the retry delays would be:

- 1st retry: 3 seconds (3 × 2^0)
- 2nd retry: 6 seconds (3 × 2^1)
- 3rd retry: 12 seconds (3 × 2^2)
- 4th retry: 24 seconds (3 × 2^3)
- 5th retry: 48 seconds (3 × 2^4)

### Version History

The version history tracks changes, improvements, and bug fixes in each version. For a full history, see the [CHANGELOG.md](https://github.com/wingify/vwo-fme-go-sdk/blob/master/CHANGELOG.md).

## Contributing

We welcome contributions to improve this SDK! Please read our [contributing guidelines](https://github.com/wingify/vwo-fme-go-sdk/blob/master/CONTRIBUTING.md) before submitting a PR.

## Code of Conduct

Our [Code of Conduct](https://github.com/wingify/vwo-fme-go-sdk/blob/master/CODE_OF_CONDUCT.md) outlines expectations for all contributors and maintainers.

## License

[Apache License, Version 2.0](https://github.com/wingify/vwo-fme-go-sdk/blob/master/LICENSE)

Copyright 2025 Wingify Software Pvt. Ltd.
