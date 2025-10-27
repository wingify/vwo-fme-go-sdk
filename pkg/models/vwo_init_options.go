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

package models

import (
	"strconv"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/storage"
)

// VWOInitOptions represents the initialization options for VWO SDK
type VWOInitOptions struct {
	AccountID            int                    `json:"accountId"`
	SDKKey               string                 `json:"sdkKey"`
	Storage              storage.Connector      `json:"-"`
	GatewayService       map[string]interface{} `json:"gatewayService,omitempty"`
	PollInterval         int                    `json:"pollInterval,omitempty"`
	Logger               map[string]interface{} `json:"logger,omitempty"`
	Integrations         *IntegrationOptions    `json:"integrations,omitempty"`
	Settings             string                 `json:"settings,omitempty"`
	IsUsageStatsDisabled bool                   `json:"isUsageStatsDisabled,omitempty"`
	VWOMeta              map[string]interface{} `json:"_vwo_meta,omitempty"`
	RetryConfig          *RetryConfig           `json:"retryConfig,omitempty"`
	IsAliasingEnabled    bool                   `json:"isAliasingEnabled,omitempty"`
	BatchEventData       map[string]interface{} `json:"batchEventData,omitempty"`
}

func NewVWOInitOptions(options map[string]interface{}) *VWOInitOptions {
	vwoInitOptions := &VWOInitOptions{}

	if sdkKey, ok := options[enums.OptionSDKKey.GetValue()].(string); ok {
		vwoInitOptions.SDKKey = sdkKey
	}
	accountIDVal := options[enums.OptionAccountID.GetValue()]
	switch v := accountIDVal.(type) {
	case int:
		vwoInitOptions.AccountID = v
	case string:
		// try to convert string to int
		// ignore errors as zero is already default value, but you may want to log/warn if conversion fails
		if parsed, err := strconv.Atoi(v); err == nil {
			vwoInitOptions.AccountID = parsed
		}
	}
	if storage, ok := options[enums.OptionStorage.GetValue()].(storage.Connector); ok {
		vwoInitOptions.Storage = storage
	}
	if gatewayService, ok := options[enums.OptionGatewayService.GetValue()].(map[string]interface{}); ok {
		vwoInitOptions.GatewayService = gatewayService
	}
	if pollInterval, ok := options[enums.OptionPollInterval.GetValue()].(int); ok {
		vwoInitOptions.PollInterval = pollInterval
	}
	if logger, ok := options[enums.OptionLogger.GetValue()].(map[string]interface{}); ok {
		vwoInitOptions.Logger = logger
	}
	if integrations, ok := options[enums.OptionIntegrations.GetValue()].(map[string]interface{}); ok {
		// Create IntegrationOptions from map
		integrationOptions := &IntegrationOptions{}
		if callback, callbackOk := integrations[enums.IntegrationCallback.GetValue()].(func(map[string]interface{})); callbackOk {
			integrationOptions.Callback = callback
		}
		vwoInitOptions.Integrations = integrationOptions
	}
	if settings, ok := options[enums.OptionSettings.GetValue()].(string); ok {
		vwoInitOptions.Settings = settings
	}
	if isUsageStatsDisabled, ok := options[enums.OptionIsUsageStatsDisabled.GetValue()].(bool); ok {
		vwoInitOptions.IsUsageStatsDisabled = isUsageStatsDisabled
	}
	if vwoMeta, ok := options[enums.OptionVWOMeta.GetValue()].(map[string]interface{}); ok {
		vwoInitOptions.VWOMeta = vwoMeta
	}
	if retryConfig, ok := options[enums.OptionRetryConfig.GetValue()].(map[string]interface{}); ok {
		vwoInitOptions.RetryConfig = NewRetryConfigFromMap(retryConfig)
	} else {
		// Set default retry config if not provided
		vwoInitOptions.RetryConfig = NewRetryConfig()
	}
	if isAliasingEnabled, ok := options[enums.OptionIsAliasingEnabled.GetValue()].(bool); ok {
		vwoInitOptions.IsAliasingEnabled = isAliasingEnabled
	}
	if batchEventData, ok := options[enums.OptionBatchEventData.GetValue()].(map[string]interface{}); ok {
		vwoInitOptions.BatchEventData = batchEventData
	}
	return vwoInitOptions
}

// GetIntegrations returns the integrations options
func (v *VWOInitOptions) GetIntegrations() *IntegrationOptions {
	return v.Integrations
}

// IntegrationOptions represents integration callback options
type IntegrationOptions struct {
	Callback func(properties map[string]interface{}) `json:"-"`
}

// NetworkOptions represents network configuration options
type NetworkOptions struct {
	Client interface{} `json:"client,omitempty"`
}
