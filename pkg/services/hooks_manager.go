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

package services

// HooksManager manages integration callbacks
type HooksManager struct {
	callback func(map[string]interface{})
	decision map[string]interface{}
}

// NewHooksManager creates a new HooksManager instance
func NewHooksManager(callback func(map[string]interface{})) *HooksManager {
	return &HooksManager{
		callback: callback,
		decision: make(map[string]interface{}),
	}
}

// Execute executes the callback with the provided properties
func (hm *HooksManager) Execute(properties map[string]interface{}) {
	if hm.callback != nil {
		hm.callback(properties)
	}
}

// Set sets properties to the decision object
func (hm *HooksManager) Set(properties map[string]interface{}) {
	if hm.callback != nil {
		hm.decision = properties
	}
}

// Get retrieves the decision object
func (hm *HooksManager) Get() map[string]interface{} {
	return hm.decision
}
