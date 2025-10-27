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

package enums

// OptionsEnum represents initialization option keys
type OptionsEnum string

const (
	// Core initialization options
	OptionSDKKey               OptionsEnum = "sdkKey"
	OptionAccountID            OptionsEnum = "accountId"
	OptionStorage              OptionsEnum = "storage"
	OptionGatewayService       OptionsEnum = "gatewayService"
	OptionPollInterval         OptionsEnum = "pollInterval"
	OptionLogger               OptionsEnum = "logger"
	OptionIntegrations         OptionsEnum = "integrations"
	OptionSettings             OptionsEnum = "settings"
	OptionIsUsageStatsDisabled OptionsEnum = "isUsageStatsDisabled"
	OptionVWOMeta              OptionsEnum = "_vwo_meta"
	OptionRetryConfig          OptionsEnum = "retryConfig"
	OptionIsAliasingEnabled    OptionsEnum = "isAliasingEnabled"
	OptionBatchEventData       OptionsEnum = "batchEventData"
)

// GetValue returns the string value of the options enum
func (o OptionsEnum) GetValue() string {
	return string(o)
}
