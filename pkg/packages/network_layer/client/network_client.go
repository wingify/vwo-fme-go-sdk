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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	retryConfig "github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"
)

// NetworkClient implements NetworkClientInterface using Go's standard http package
type NetworkClient struct {
	retryConfig *retryConfig.RetryConfig
}

// NewNetworkClient creates a new NetworkClient instance
func NewNetworkClient() *NetworkClient {
	return &NetworkClient{
		retryConfig: retryConfig.NewRetryConfig(),
	}
}

// NewNetworkClientWithRetry creates a new NetworkClient instance with retry configuration
func NewNetworkClientWithRetry(config *retryConfig.RetryConfig) *NetworkClient {
	if config == nil {
		config = retryConfig.NewRetryConfig()
	}
	return &NetworkClient{
		retryConfig: config,
	}
}

// SetRetryConfig sets the retry configuration for the network client
func (networkClient *NetworkClient) SetRetryConfig(config *retryConfig.RetryConfig) {
	if config != nil {
		networkClient.retryConfig = config
	}
}

// executeWithRetry executes a network operation with retry logic
func (networkClient *NetworkClient) executeWithRetry(operation func() *models.ResponseModel, eventName string) *models.ResponseModel {
	var lastResponse *models.ResponseModel
	var lastError error
	for attempt := 0; attempt <= networkClient.retryConfig.MaxRetries; attempt++ {
		response := operation()
		response.TotalAttempts = attempt

		// If successful or not retryable, return immediately
		if response.Error == nil && ((response.StatusCode >= 200 && response.StatusCode < 300) || response.StatusCode == 400) {
			response.Error = lastError
			return response
		}

		// Check if we should retry
		if !networkClient.retryConfig.IsRetryable(attempt) || eventName == enums.VWODebuggerEvent.GetValue() {
			response.Error = lastError
			return response
		}

		lastError = response.Error
		lastResponse = response

		// Calculate delay for next retry
		if attempt <= networkClient.retryConfig.MaxRetries {
			delay := networkClient.retryConfig.GetRetryDelay(attempt)
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}

	return lastResponse
}

// constructURL constructs the full URL from network options
func constructURL(networkOptions map[string]interface{}) string {
	hostname, _ := networkOptions[enums.NetworkOptionsHostname.GetValue()].(string)
	path, _ := networkOptions[enums.NetworkOptionsPath.GetValue()].(string)
	scheme, _ := networkOptions[enums.NetworkOptionsScheme.GetValue()].(string)

	// Add port if specified and not default
	if port, ok := networkOptions[enums.NetworkOptionsPort.GetValue()].(int); ok && port != 0 {
		hostname = fmt.Sprintf("%s:%d", hostname, port)
	}

	return fmt.Sprintf("%s://%s%s", strings.ToLower(scheme), hostname, path)
}

// createHTTPClient creates an HTTP client with timeout configuration
func createHTTPClient() *http.Client {
	client := &http.Client{}
	client.Timeout = constants.NetworkTimeout
	return client
}

// GET performs a GET request using the provided RequestModel
func (networkClient *NetworkClient) GET(requestModel *models.RequestModel) *models.ResponseModel {
	return networkClient.executeWithRetry(func() *models.ResponseModel {
		responseModel := models.NewResponseModel()

		networkOptions := requestModel.GetOptions()
		url := constructURL(networkOptions)

		req, err := http.NewRequest(enums.HTTPMethodGET.GetValue(), url, nil)
		if err != nil {
			responseModel.Error = fmt.Errorf("%s: %w", enums.NetworkClientErrorFailedToCreateGETRequest.GetValue(), err)
			return responseModel
		}

		// Set headers if provided
		if headers, ok := networkOptions[enums.NetworkOptionsHeaders.GetValue()].(map[string]string); ok {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}

		client := createHTTPClient()
		resp, err := client.Do(req)
		if err != nil {
			responseModel.Error = fmt.Errorf("%s: %w", enums.NetworkClientErrorGETRequestFailed.GetValue(), err)
			return responseModel
		}
		defer resp.Body.Close()

		responseModel.StatusCode = resp.StatusCode

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			responseModel.Error = fmt.Errorf("%s: %w", enums.NetworkClientErrorFailedToReadResponse.GetValue(), err)
			return responseModel
		}

		contentType := resp.Header.Get(enums.HTTPHeaderContentType.GetValue())

		// Validate response
		if resp.StatusCode != 200 || !strings.Contains(contentType, enums.ContentTypeApplicationJSON.GetValue()) {
			responseModel.Error = fmt.Errorf(enums.NetworkClientErrorInvalidResponse.GetValue(), string(body), resp.StatusCode, resp.Status)
			return responseModel
		}

		responseModel.Data = string(body)

		return responseModel
	}, "")
}

// POST performs a POST request using the provided RequestModel
func (networkClient *NetworkClient) POST(request *models.RequestModel) *models.ResponseModel {
	return networkClient.executeWithRetry(func() *models.ResponseModel {
		responseModel := models.NewResponseModel()

		networkOptions := request.GetOptions()
		url := constructURL(networkOptions)

		// Prepare request body
		var reqBody io.Reader
		if body, ok := networkOptions[enums.NetworkOptionsBody.GetValue()].(map[string]interface{}); ok {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				responseModel.Error = fmt.Errorf("%s: %w", enums.NetworkClientErrorFailedToMarshalBody.GetValue(), err)
				return responseModel
			}
			reqBody = bytes.NewBuffer(jsonBody)
		}

		req, err := http.NewRequest(enums.HTTPMethodPOST.GetValue(), url, reqBody)
		if err != nil {
			responseModel.Error = fmt.Errorf("%s: %w", enums.NetworkClientErrorFailedToCreatePOSTRequest.GetValue(), err)
			return responseModel
		}

		// Set headers
		for key, value := range networkOptions {
			if key == enums.NetworkOptionsHeaders.GetValue() {
				if headers, ok := value.(map[string]string); ok {
					for headerKey, headerValue := range headers {
						req.Header.Set(headerKey, headerValue)
					}
				}
			}
		}

		req.Header.Set("User-Agent", constants.SDKName)
		client := createHTTPClient()
		resp, err := client.Do(req)
		if err != nil {
			responseModel.Error = fmt.Errorf("%s: %w", enums.NetworkClientErrorPOSTRequestFailed.GetValue(), err)
			return responseModel
		}
		defer resp.Body.Close()

		responseModel.StatusCode = resp.StatusCode

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			responseModel.Error = fmt.Errorf("%s: %w", enums.NetworkClientErrorFailedToReadResponse.GetValue(), err)
			return responseModel
		}

		responseModel.Data = string(body)

		if resp.StatusCode != 200 {
			responseModel.Error = fmt.Errorf(enums.NetworkClientErrorRequestFailed.GetValue(), resp.StatusCode, string(body))
		}

		return responseModel
	}, request.GetEventName())
}
