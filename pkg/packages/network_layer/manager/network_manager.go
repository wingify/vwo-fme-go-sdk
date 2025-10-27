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

package manager

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/client"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/handlers"
	networkModels "github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"
)

// FlushCallback is a callback function type for flush operations
type FlushCallback func(error string, events string)

// NetworkManager manages network operations with singleton pattern
type NetworkManager struct {
	config      *networkModels.GlobalRequestModel
	client      client.NetworkClientInterface
	retryConfig *models.RetryConfig
	mu          sync.Mutex
}

// AttachClient attaches a custom network client
func (networkManager *NetworkManager) AttachClient(customClient client.NetworkClientInterface) {
	networkManager.mu.Lock()
	defer networkManager.mu.Unlock()

	networkManager.client = customClient
	networkManager.config = networkModels.NewGlobalRequestModel("", nil, nil, nil)
}

// AttachDefaultClient attaches the default network client
func (networkManager *NetworkManager) AttachDefaultClient() {
	networkManager.mu.Lock()
	defer networkManager.mu.Unlock()

	networkManager.client = client.NewNetworkClient()
	networkManager.config = networkModels.NewGlobalRequestModel("", nil, nil, nil)
}

// AttachDefaultClientWithRetry attaches the default network client with retry configuration
func (networkManager *NetworkManager) AttachDefaultClientWithRetry(retryConfig *models.RetryConfig) {
	networkManager.mu.Lock()
	defer networkManager.mu.Unlock()

	networkManager.retryConfig = retryConfig
	networkManager.client = client.NewNetworkClientWithRetry(retryConfig)
	networkManager.config = networkModels.NewGlobalRequestModel("", nil, nil, nil)
}

// SetRetryConfig sets the retry configuration for the network manager
func (networkManager *NetworkManager) SetRetryConfig(retryConfig *models.RetryConfig) {
	networkManager.mu.Lock()
	defer networkManager.mu.Unlock()

	networkManager.retryConfig = retryConfig
	if networkManager.client != nil {
		if networkClient, ok := networkManager.client.(*client.NetworkClient); ok {
			networkClient.SetRetryConfig(retryConfig)
		}
	}
}

// SetConfig sets the global request configuration
func (networkManager *NetworkManager) SetConfig(config *networkModels.GlobalRequestModel) {
	networkManager.mu.Lock()
	defer networkManager.mu.Unlock()

	networkManager.config = config
}

// GetConfig returns the global request configuration
func (networkManager *NetworkManager) GetConfig() *networkModels.GlobalRequestModel {
	networkManager.mu.Lock()
	defer networkManager.mu.Unlock()

	return networkManager.config
}

// CreateRequest creates a RequestModel from the given request
func (networkManager *NetworkManager) CreateRequest(request *networkModels.RequestModel) *networkModels.RequestModel {
	networkManager.mu.Lock()
	config := networkManager.config
	networkManager.mu.Unlock()

	handler := handlers.NewRequestHandler()
	return handler.CreateRequest(request, config)
}

// Get synchronously sends a GET request to the server
func (networkManager *NetworkManager) Get(request *networkModels.RequestModel) *networkModels.ResponseModel {
	networkOptions := networkManager.CreateRequest(request)
	if networkOptions == nil {
		return nil
	} else {
		networkManager.mu.Lock()
		currentClient := networkManager.client
		networkManager.mu.Unlock()

		return currentClient.GET(request)
	}
}

// Post synchronously sends a POST request to the server
func (networkManager *NetworkManager) Post(request *networkModels.RequestModel, flushCallback FlushCallback) *networkModels.ResponseModel {
	response := &networkModels.ResponseModel{}
	networkOptions := networkManager.CreateRequest(request)
	if networkOptions == nil {
		return nil
	}
	// Perform the actual POST request
	networkManager.mu.Lock()
	currentClient := networkManager.client
	networkManager.mu.Unlock()

	response = currentClient.POST(request)

	// Handle the response and trigger callback based on success or failure
	if response != nil && response.StatusCode >= 200 && response.StatusCode < 300 {
		if flushCallback != nil {
			body, err := json.Marshal(request.Body)
			if err == nil {
				flushCallback("", string(body)) // Success, pass response body to callback
			} else {
				flushCallback("", fmt.Sprintf("%v", request.Body)) // Success, pass response body to callback
			}
		}
	} else {
		if flushCallback != nil {
			body, err := json.Marshal(request.Body)
			if err == nil {
				flushCallback(fmt.Sprintf("Failed with status code: %d", response.StatusCode), string(body)) // Failure, pass error message
			} else {
				flushCallback(fmt.Sprintf("Failed with status code: %d", response.StatusCode), fmt.Sprintf("%v", request.Body)) // Failure, pass error message
			}
		}
	}
	return response
}

// PostAsync asynchronously sends a POST request to the server using a goroutine
func (networkManager *NetworkManager) PostAsync(request *networkModels.RequestModel, flushCallback FlushCallback) {
	go func() {
		networkManager.Post(request, flushCallback)
	}()
}
