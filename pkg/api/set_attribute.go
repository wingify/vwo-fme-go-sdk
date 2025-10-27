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
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/utils"
)

// SetAttribute sets user attributes and sends an impression to vwo
func SetAttribute(attributeMap map[string]interface{}, context *user.VWOContext, serviceContainer interfaces.ServiceContainerInterface) {
	createAndSendImpressionForSetAttribute(attributeMap, context, serviceContainer)
}

// createAndSendImpressionForSetAttribute creates and sends an impression for setting attributes
func createAndSendImpressionForSetAttribute(
	attributeMap map[string]interface{},
	context *user.VWOContext,
	serviceContainer interfaces.ServiceContainerInterface,
) {
	serviceContainer.GetDebuggerService().AddStandardDebugProp(enums.DebugPropAPI.GetValue(), string(enums.ApiSetAttribute))
	// Get base properties for the event
	properties := utils.GetEventsBaseProperties(
		serviceContainer.GetSettingsManager(),
		enums.VWOSyncVisitorProp.GetValue(),
		context.GetUserAgent(),
		context.GetIPAddress(),
	)

	// Construct payload data for setting attributes
	payload := utils.GetAttributePayloadData(
		serviceContainer,
		context.GetID(),
		enums.VWOSyncVisitorProp.GetValue(),
		attributeMap,
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
