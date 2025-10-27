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
	"fmt"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// SendDebugEventToVWO sends a debug event to VWO.
// @param settingsManager The settings manager containing configuration
// @param eventProps The properties for the event
func SendDebugEventToVWO(settingsManager interfaces.SettingsManagerInterface, eventProps map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			settingsManager.GetLoggerService().Error("ERROR_SENDING_DEBUG_EVENT_TO_VWO", map[string]interface{}{
				"err": fmt.Sprintf("%v", r),
			}, nil)
		}
	}()

	// Create query parameters
	properties := GetEventsBaseProperties(
		settingsManager,
		enums.VWODebuggerEvent.GetValue(),
		"",
		"",
	)

	// Create payload
	payload := GetDebuggerEventPayload(
		settingsManager,
		eventProps,
	)

	// Send event
	SendEventDirectlyToDACDN(
		settingsManager,
		properties,
		payload,
		enums.VWODebuggerEvent.GetValue(),
	)
}
