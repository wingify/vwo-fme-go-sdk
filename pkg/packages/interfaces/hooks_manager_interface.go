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

// HooksManagerInterface defines the contract for hooks manager
// This interface is used to break import cycles and provide abstraction for hooks operations
type HooksManagerInterface interface {
	// Execute executes the callback with the provided properties
	Execute(properties map[string]interface{})

	// Set sets properties to the decision object
	Set(properties map[string]interface{})

	// Get retrieves the decision object
	Get() map[string]interface{}
}
