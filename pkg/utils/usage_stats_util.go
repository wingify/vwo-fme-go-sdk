/**
 * Copyright 2025 Wingify Software Pvt. Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"runtime"

	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/enums"

	eventEnums "github.com/wingify/vwo-fme-go-sdk/pkg/enums"
)

// SetUsageStats sets usage statistics based on provided options.
// Maps various SDK features and configurations to boolean flags.
func GetUsageStats(options *models.VWOInitOptions) map[string]interface{} {
	data := make(map[string]interface{})

	// Map configuration options to usage stats flags
	if options.GetIntegrations() != nil {
		data["ig"] = 1
	}

	// Check if the logger has transports in it
	if options.Logger != nil {
		if _, hasTransport := options.Logger["transport"]; hasTransport {
			data["cl"] = 1
		}
		if _, hasTransports := options.Logger["transports"]; hasTransports {
			data["cl"] = 1
		}
	}

	// Check the logger level
	// If the level is not valid, push -1
	// If the level is valid, push the enum value
	if options.Logger != nil {
		if level, exists := options.Logger["level"]; exists {
			if levelStr, ok := level.(string); ok {
				logLevel := enums.ParseLogLevel(levelStr)
				data["ll"] = logLevel.GetLevel()
			} else {
				data["ll"] = -1
			}
		}
	}

	if options.Storage != nil {
		data["ss"] = 1
	}

	// Check if the gatewayService is not an empty map
	if len(options.GatewayService) > 0 {
		data["gs"] = 1
	}

	if options.PollInterval > 0 {
		data["pi"] = 1
	}

	// Handle _vwo_meta
	if options.VWOMeta != nil {
		if _, hasEA := options.VWOMeta["ea"]; hasEA {
			data["_ea"] = 1
		}
	}

	// Get Go version
	data["lv"] = runtime.Version()

	return data
}

// SendSDKInitEvent sends the SDK init event to the DACDN
func SendSDKInitEvent(
	serviceContainer interfaces.ServiceContainerInterface,
	context *user.VWOContext,
	settingsFetchTime int,
	sdkInitTime int,
) {
	// Get base properties for the event
	properties := GetEventsBaseProperties(
		serviceContainer.GetSettingsManager(),
		eventEnums.VWOSDKInitEvent.GetValue(),
		"",
		"",
	)

	// Construct payload data for tracking the user
	payload := GetSDKInitEventPayload(
		serviceContainer.GetSettingsManager(),
		context.GetID(),
		eventEnums.VWOSDKInitEvent.GetValue(),
		settingsFetchTime,
		sdkInitTime,
	)

	// Check if batch event queue is available
	if serviceContainer.GetBatchEventQueue().IsInitialized() {
		// Enqueue the event to the batch queue for future processing
		serviceContainer.GetBatchEventQueue().Enqueue(payload)
	} else {
		SendPostAPIRequest(
			serviceContainer,
			properties,
			payload,
			context,
			map[string]interface{}{},
		)
	}
}

// SendSDKUsageStatsEvent sends the SDK usage stats event to the DACDN
func SendSDKUsageStatsEvent(
	serviceContainer interfaces.ServiceContainerInterface,
	context *user.VWOContext,
	usageStatsAccountId int,
) {
	// Get base properties for the event
	properties := GetEventsBaseProperties(
		serviceContainer.GetSettingsManager(),
		eventEnums.VWOSDKUsageStats.GetValue(),
		"",
		"",
	)

	// Construct payload data for tracking the user
	payload := GetSDKUsageStatsEventPayload(
		serviceContainer.GetSettingsManager(),
		context.GetID(),
		eventEnums.VWOSDKUsageStats.GetValue(),
		usageStatsAccountId,
		GetUsageStats(serviceContainer.GetVWOInitOptions()),
	)

	// Check if batch event queue is available
	if serviceContainer.GetBatchEventQueue().IsInitialized() {
		// Enqueue the event to the batch queue for future processing
		serviceContainer.GetBatchEventQueue().Enqueue(payload)
	} else {
		SendPostAPIRequest(
			serviceContainer,
			properties,
			payload,
			context,
			map[string]interface{}{},
		)
	}
}
