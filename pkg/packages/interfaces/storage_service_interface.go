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

import (
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
)

// StorageServiceInterface defines the contract for storage service
// This interface is used to break import cycles and provide abstraction for storage operations
type StorageServiceInterface interface {
	// GetDataInStorage retrieves data from storage based on feature key and user context
	// Returns the stored data as a map and any error that occurred during retrieval
	GetDataInStorage(featureKey string, context *user.VWOContext) (map[string]interface{}, error)

	// SetDataInStorage stores data in storage
	// Returns true if the data was successfully stored, false otherwise
	SetDataInStorage(data map[string]interface{}) bool
}
