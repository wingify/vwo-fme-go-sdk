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

package api

import (
	"fmt"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/utils"
)

// TrackEvent tracks a custom event and returns true if the event is tracked successfully
func TrackEvent(eventName string, context *user.VWOContext, eventProperties map[string]interface{}, serviceContainer interfaces.ServiceContainerInterface) bool {
	serviceContainer.GetDebuggerService().AddStandardDebugProp(enums.DebugPropAPI.GetValue(), string(enums.ApiTrackEvent))
	// Try-catch equivalent in Go using defer and recover
	defer func() {
		if r := recover(); r != nil {
			serviceContainer.GetLoggerService().Error("EXECUTION_FAILED", map[string]interface{}{"apiName": string(enums.ApiTrackEvent), "err": fmt.Sprintf("Error in Tracking event: %s : %v", eventName, r)}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	if utils.DoesEventBelongToAnyFeature(eventName, serviceContainer.GetSettings()) {
		createAndSendImpressionForTrack(eventName, context, eventProperties, serviceContainer)

		// Set hooks data
		objectToReturn := map[string]interface{}{
			"eventName": eventName,
			"api":       string(enums.ApiTrackEvent),
		}
		serviceContainer.GetHooksManager().Set(objectToReturn)
		serviceContainer.GetHooksManager().Execute(serviceContainer.GetHooksManager().Get())

		return true
	} else {
		serviceContainer.GetLoggerService().Error("EVENT_NOT_FOUND", map[string]interface{}{"eventName": eventName}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return false
	}
}

// createAndSendImpressionForTrack creates and sends an impression for a track event
func createAndSendImpressionForTrack(
	eventName string,
	context *user.VWOContext,
	eventProperties map[string]interface{},
	serviceContainer interfaces.ServiceContainerInterface,
) {
	// Get base properties for the event
	properties := utils.GetEventsBaseProperties(
		serviceContainer.GetSettingsManager(),
		eventName,
		context.GetUserAgent(),
		context.GetIPAddress(),
	)

	// Construct payload data for tracking the event
	payload := utils.GetTrackGoalPayloadData(
		serviceContainer,
		context.GetID(),
		eventName,
		context,
		eventProperties,
	)

	// Check if batch event queue is available
	if serviceContainer.GetBatchEventQueue().IsInitialized() {
		// Enqueue the event to the batch queue for future processing
		serviceContainer.GetBatchEventQueue().Enqueue(payload)
	} else {
		// Send the event immediately if batch event queue is not available
		utils.SendPostAPIRequest(
			serviceContainer,
			properties,
			payload,
			context,
			map[string]interface{}{},
		)
	}
}
