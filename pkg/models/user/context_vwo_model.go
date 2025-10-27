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

package user

// ContextVWO represents VWO-specific context information
type ContextVWO struct {
	Location  map[string]string `json:"location,omitempty"`
	UserAgent map[string]string `json:"userAgent,omitempty"`
}

// NewContextVWO creates a new ContextVWO from a map
func NewContextVWO(context map[string]interface{}) *ContextVWO {
	contextVWO := &ContextVWO{}

	if location, ok := context["location"]; ok {
		if locMap, ok := location.(map[string]string); ok {
			contextVWO.Location = locMap
		}
	}

	if userAgent, ok := context["userAgent"]; ok {
		if uaMap, ok := userAgent.(map[string]string); ok {
			contextVWO.UserAgent = uaMap
		}
	}

	return contextVWO
}

// GetLocation returns the location information
func (c *ContextVWO) GetLocation() map[string]string {
	return c.Location
}

// GetUaInfo returns the user agent information
func (c *ContextVWO) GetUaInfo() map[string]string {
	return c.UserAgent
}
