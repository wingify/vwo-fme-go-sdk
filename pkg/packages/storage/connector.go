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

package storage

// Connector defines the contract for storage implementations (Go interface approach)
type Connector interface {
	// Set stores data in the storage
	Set(data map[string]interface{}) error

	// Get retrieves data from storage based on feature key and user ID
	Get(featureKey string, userID string) (interface{}, error)
}
