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

package gateway

import (
	"fmt"
	"net/url"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	networkModels "github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"
)

// GetFromGatewayService fetches data from the gateway service with proper error handling
func GetFromGatewayService(
	serviceContainer interfaces.ServiceContainerInterface,
	queryParams map[string]string,
	endpoint string,
) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			serviceContainer.GetLoggerService().Error("ERROR_FETCHING_DATA_FROM_GATEWAY", map[string]interface{}{
				"err": fmt.Sprintf("%v", r),
			}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	if serviceContainer.GetBaseUrl() == constants.HostName {
		serviceContainer.GetLoggerService().Error("INVALID_GATEWAY_URL", nil, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return "", fmt.Errorf("gateway URL error: base URL is default hostname")
	}

	// Create request model
	request := networkModels.NewRequestModel(
		serviceContainer.GetBaseUrl(),
		enums.ApiMethodGet.GetValue(),
		endpoint,
		queryParams,
		nil,
		nil,
		serviceContainer.GetSettingsManager().GetProtocol(),
		serviceContainer.GetSettingsManager().GetPort(),
		"",
	)

	// Execute GET request
	response := serviceContainer.GetNetworkManager().Get(request)
	if response == nil || response.Error != nil {
		errMsg := "network call failed"
		if response != nil && response.Error != nil {
			errMsg = response.Error.Error()
		}
		serviceContainer.GetLoggerService().Error("ERROR_FETCHING_DATA_FROM_GATEWAY", map[string]interface{}{
			"err": errMsg,
		}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return "", fmt.Errorf("network call failed: %s", errMsg)
	}

	return response.Data, nil
}

// GetQueryParams encodes the query parameters to ensure they are URL-safe
func GetQueryParams(queryParams map[string]string) map[string]string {
	encodedParams := make(map[string]string)

	for key, value := range queryParams {
		// Convert the value to a string
		valStr := fmt.Sprintf("%v", value)

		// Encode the parameter value to ensure it is URL-safe
		encodedValue := url.QueryEscape(valStr)
		// Add the encoded parameter to the result map
		encodedParams[key] = encodedValue
	}

	return encodedParams
}
