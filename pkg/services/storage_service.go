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

import (
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/storage"
)

// StorageService handles storage operations
type StorageService struct{}

// NewStorageService creates a new StorageService instance
func NewStorageService() *StorageService {
	return &StorageService{}
}

// GetDataInStorage retrieves data from storage
func (s *StorageService) GetDataInStorage(featureKey string, context *user.VWOContext) (map[string]interface{}, error) {
	// Get storage instance
	storageInstance := storage.GetInstance().GetConnector()
	if storageInstance == nil {
		return nil, nil
	}

	// Call connector's Get method
	result, err := storageInstance.Get(featureKey, context.ID)
	if err != nil {
		return nil, err
	}

	// Cast result to map
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return nil, nil
}

// SetDataInStorage stores data in storage
func (s *StorageService) SetDataInStorage(data map[string]interface{}) bool {
	// Get storage instance
	storageInstance := storage.GetInstance().GetConnector()
	if storageInstance == nil {
		return false
	}

	// Use defer with recover
	defer func() {
		if r := recover(); r != nil {
			// Exception occurred, return false
		}
	}()

	// Call connector's Set method
	err := storageInstance.Set(data)
	if err != nil {
		return false
	}

	return true
}
