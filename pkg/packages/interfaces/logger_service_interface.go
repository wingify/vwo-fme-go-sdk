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

package interfaces

import "github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/enums"

// LoggerServiceInterface defines the contract for logger service
// This interface is used to break import cycles and provide abstraction over the logger implementation
type LoggerServiceInterface interface {
	// Basic logging methods
	Trace(message string)
	Debug(message string)
	Info(message string)
	Warn(message string)

	// Error method with template support and extra data
	Error(template string, templateData map[string]interface{}, extraData map[string]interface{}, shouldSendDebugEvent ...bool)

	// Generic log method with level
	Log(level enums.LogLevel, message string)

	// SetSettingsManager method
	SetSettingsManager(settingsManager SettingsManagerInterface)
}
