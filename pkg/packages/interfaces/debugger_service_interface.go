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

// DebuggerServiceInterface defines the contract for debugger service
// This interface is used to break import cycles and provide abstraction for debug operations
type DebuggerServiceInterface interface {
	// Standard debug props methods
	AddStandardDebugProps(standardDebugProps map[string]interface{})
	AddStandardDebugProp(key string, value interface{})
	GetStandardDebugProps() map[string]interface{}

	// Category-specific debug props methods
	AddCategoryDebugProps(category string, eventProps map[string]interface{})
	AddCategoryDebugProp(category string, key string, value interface{})
	GetDebugEventProps(category string) map[string]interface{}

	// Utility methods
	Clear()
}
