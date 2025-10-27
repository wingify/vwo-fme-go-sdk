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

package handlers

import (
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"
)

// RequestHandler handles request creation and merging with global configuration
type RequestHandler struct{}

// NewRequestHandler creates a new RequestHandler instance
func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

// CreateRequest creates a new request by merging properties from a base request and a configuration model
// If both the request URL and the base URL from the configuration are missing, it returns nil
func (requestHandler *RequestHandler) CreateRequest(request *models.RequestModel, config *models.GlobalRequestModel) *models.RequestModel {
	// Check if both the request URL and the configuration base URL are missing
	if config.BaseURL == "" && request.URL == "" {
		return nil // Return nil if no URL is specified
	}

	// Set the request URL, defaulting to the configuration base URL if not set
	if request.URL == "" {
		request.URL = config.BaseURL
	}

	// Set the request timeout, defaulting to the configuration timeout if not set
	if request.Timeout == -1 {
		request.Timeout = config.Timeout
	}

	// Set the request body, defaulting to the configuration body if not set
	if request.Body == nil {
		request.Body = config.Body
	}

	// Set the request headers, defaulting to the configuration headers if not set
	if request.Headers == nil {
		request.Headers = config.Headers
	}

	// Initialize request query parameters, defaulting to an empty map if not set
	if request.Query == nil {
		request.Query = make(map[string]string)
	}

	// Initialize configuration query parameters, defaulting to an empty map if not set
	configQueryParams := config.Query
	if configQueryParams == nil {
		configQueryParams = make(map[string]interface{})
	}

	// Merge configuration query parameters into the request query parameters if they don't exist
	for key, value := range configQueryParams {
		if _, exists := request.Query[key]; !exists {
			if strValue, ok := value.(string); ok {
				request.Query[key] = strValue
			}
		}
	}

	return request
}
