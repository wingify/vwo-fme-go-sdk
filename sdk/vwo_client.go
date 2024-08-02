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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wingify/vwo-fme-go-sdk/sdk/constants"
	"github.com/wingify/vwo-fme-go-sdk/sdk/httpclient"
	"github.com/wingify/vwo-fme-go-sdk/sdk/models"
)

// VWOClient represents the client that interacts with VWO services.
type VWOClient struct {
	settings string
	options  models.VWOInitOptions
}

// NewVWOClient initializes a new instance of VWOClient with the provided settings and options.
// It sets up the client with the specified settings and options.
func NewVWOClient(settings string, options models.VWOInitOptions) *VWOClient {
	return &VWOClient{settings: settings, options: options}
}

// GetSettings returns the settings of the VWOClient.
// It retrieves the configuration settings of the VWO client.
func (vc *VWOClient) GetSettings() string {
	return vc.settings
}

// GetFlag retrieves the feature flag for a given feature key and context.
// It sends a request to the VWO server to get the status and variables of a feature flag.
func (vc *VWOClient) GetFlag(featureKey string, context map[string]interface{}) (*GetFlagResponse, error) {
	client := httpclient.GetClient()

	if featureKey == "" {
		return nil, errors.New("featureKey is required to get the feature flag")
	}

	if context == nil {
		return nil, errors.New("context is required to get the feature flag")
	}

	context["featureKey"] = featureKey

	// Serialize the userContext map to JSON
	parsedBody, err := json.Marshal(context)
	if err != nil {
		return nil, fmt.Errorf("error serializing user context to JSON: %w", err)
	}

	endpoint := constants.EndPointsGetFlag + "?accountId=" + vc.options.AccountID + "&sdkKey=" + vc.options.SdkKey

	response, err := client.PostRequest(endpoint, parsedBody)
	if err != nil {
		return nil, fmt.Errorf("error making post request: %w", err)
	}

	return NewGetFlagResponse(response)
}

// TrackEvent tracks the event for the user.
// It sends a request to the VWO server to log a user event.
func (vc *VWOClient) TrackEvent(eventName string, context map[string]interface{}, eventProperties map[string]interface{}) (string, error) {
	if eventName == "" {
		return "", errors.New("eventName is required to track the event")
	}

	if context == nil || context["userId"] == nil || context["userId"] == "" {
		return "", errors.New("userId is required to track the event")
	}

	client := httpclient.GetClient()

	context["eventName"] = eventName
	context["eventProperties"] = eventProperties

	// Serialize the userContext map to JSON
	parsedBody, err := json.Marshal(context)
	if err != nil {
		return "", fmt.Errorf("error serializing user context to JSON: %w", err)
	}

	endpoint := constants.EndPointsTrackEvent + "?accountId=" + vc.options.AccountID + "&sdkKey=" + vc.options.SdkKey

	response, err := client.PostRequest(endpoint, parsedBody)
	if err != nil {
		return "", fmt.Errorf("error making post request: %w", err)
	}

	return string(response), nil
}

// SetAttribute sets the attribute for the user.
// It sends a request to the VWO server to set a user attribute.
func (vc *VWOClient) SetAttribute(attributeKey string, attributeValue interface{}, context map[string]interface{}) error {
	if attributeKey == "" {
		return errors.New("attributeKey is required for setAttribute")
	}

	if attributeValue == "" || attributeValue == nil {
		return errors.New("attributeValue is required for setAttribute")
	}

	if context == nil || context["userId"] == nil || context["userId"] == "" {
		return errors.New("userId is required for setAttribute")
	}

	client := httpclient.GetClient()

	context["attributeKey"] = attributeKey
	context["attributeValue"] = attributeValue

	// Serialize the userContext map to JSON
	parsedBody, err := json.Marshal(context)
	if err != nil {
		return fmt.Errorf("error serializing user context to JSON: %w", err)
	}

	endpoint := constants.EndPointsSetAttirbute + "?accountId=" + vc.options.AccountID + "&sdkKey=" + vc.options.SdkKey

	_, err = client.PostRequest(endpoint, parsedBody)
	if err != nil {
		return fmt.Errorf("error making post request: %w", err)
	}

	return nil
}
