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
package models

// VWOInitOptions represents initialization options for VWO.
type VWOInitOptions struct {
	SdkKey            string
	AccountID         string
	GatewayServiceURL string
}

// NewVWOInitOptions initializes a new instance of VWOInitOptions with the provided options.
func NewVWOInitOptions(sdkKey string, accountID string, gatewayServiceURL string) *VWOInitOptions {
	return &VWOInitOptions{
		SdkKey:            sdkKey,
		AccountID:         accountID,
		GatewayServiceURL: gatewayServiceURL,
	}
}
