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

package data

import (
	"encoding/json"
	"fmt"
	"sync"

	storageModels "github.com/wingify/vwo-fme-go-sdk/pkg/models/storage"
)

// StorageTest implements the storage.Connector interface for testing
type StorageTest struct {
	data map[string]map[string]interface{}
	mu   sync.RWMutex
}

// NewStorageTest creates a new StorageTest instance
func NewStorageTest() *StorageTest {
	return &StorageTest{
		data: make(map[string]map[string]interface{}),
	}
}

// Set stores data in the test storage
func (s *StorageTest) Set(data map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	featureKey, ok := data["featureKey"].(string)
	if !ok {
		return fmt.Errorf("featureKey not found or not a string")
	}

	userID, ok := data["userId"].(string)
	if !ok {
		return fmt.Errorf("userId not found or not a string")
	}

	key := featureKey + ":" + userID
	s.data[key] = data

	return nil
}

// Get retrieves data from the test storage
func (s *StorageTest) Get(featureKey string, userID string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := featureKey + ":" + userID
	if data, exists := s.data[key]; exists {
		return data, nil
	}

	return nil, nil
}

// GetStorageData retrieves and parses storage data as StorageData struct
func (s *StorageTest) GetStorageData(featureKey string, userID string) (*storageModels.StorageData, error) {
	data, err := s.Get(featureKey, userID)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	// Convert to JSON and back to parse into StorageData struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var storageData storageModels.StorageData
	if err := json.Unmarshal(jsonData, &storageData); err != nil {
		return nil, err
	}

	return &storageData, nil
}
