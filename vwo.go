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
package vwo

import (
	"errors"

	"github.com/wingify/vwo-fme-go-sdk/sdk"
	"github.com/wingify/vwo-fme-go-sdk/sdk/models"
)

var instance *sdk.VWOClient

// SetInstance configures and builds the VWO instance using the provided options.
// It initializes the VWO client with the given options and fetches the settings.
func SetInstance(options models.VWOInitOptions) (*sdk.VWOClient, error) {
	vwoBuilder := sdk.NewVWOBuilder(options)
	vwoBuilder.InitClient()
	settings := vwoBuilder.GetSettings(false)
	return sdk.NewVWOClient(settings, options), nil
}

// GetInstance gets the singleton instance of VWO.
// It returns the VWO client instance if it has been initialized.
func GetInstance() *sdk.VWOClient {
	return instance
}

// Init initializes the VWO instance with the provided options.
// It validates the options and sets up the VWO client instance.
func Init(options map[string]interface{}) (*sdk.VWOClient, error) {

	if options["sdkKey"] == nil || options["sdkKey"].(string) == "" {
		return nil, errors.New("sdkKey is required to initialize VWO. Please provide the sdkKey in the options")
	}
	if options["accountId"] == nil || options["accountId"].(string) == "" {
		return nil, errors.New("accountId is required to initialize VWO. Please provide the accountId in the options")
	}

	if options["gatewayServiceURL"] == nil || options["gatewayServiceURL"] == "" {
		return nil, errors.New("gatewayServiceURL is required to initialize VWO. Please provide the gatewayServiceURL in the options")
	}

	VWOInitOptions := models.NewVWOInitOptions(options["sdkKey"].(string), options["accountId"].(string), options["gatewayServiceURL"].(string))

	var err error
	instance, err = SetInstance(*VWOInitOptions)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
