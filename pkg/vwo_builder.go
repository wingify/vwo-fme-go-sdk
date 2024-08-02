/**
 * Copyright 2024 Wingify Software Pvt. Ltd.
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
package sdk

import (
	"sync"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/httpclient"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
)

// VWOBuilder is responsible for building and initializing the VWO client.
type VWOBuilder struct {
	options models.VWOInitOptions
	mu      sync.Mutex
}

// NewVWOBuilder creates a new instance of VWOBuilder with the provided options.
// It initializes the VWOBuilder with the given VWOInitOptions.
func NewVWOBuilder(options models.VWOInitOptions) *VWOBuilder {
	return &VWOBuilder{
		options: options,
	}
}

// InitClient initializes the HTTP client for VWOBuilder.
// It sets the base URL for the client based on the GatewayServiceURL option or a default value.
func (vb *VWOBuilder) InitClient() {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	// Check if GatewayServiceURL is provided
	if vb.options.GatewayServiceURL != "" {
		httpclient.InitializeClient(vb.options.GatewayServiceURL)
	} else {
		httpclient.InitializeClient(constants.HTTPSProtocol + constants.BaseURL)
	}
}

// GetSettings fetches the settings for the VWO account.
// It returns the settings as a string. If there's an error, it panics.
func (vb *VWOBuilder) GetSettings(isDevelopmentMode bool) string {
	vb.mu.Lock()
	defer vb.mu.Unlock()

	client := httpclient.GetClient()
	endpoint := constants.EndPointsAccountSettings + "?a=" + vb.options.AccountID + "&i=" + vb.options.SdkKey

	body, err := client.DoRequest("GET", endpoint, nil, nil)
	if err != nil {
		panic(err)
	}
	return string(body)
}
