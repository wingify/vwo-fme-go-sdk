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

package vwo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/api"
	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/core"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/schemas"
	settingsModel "github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/utils"
)

// VWOClient represents the VWO client with API methods
type VWOClient struct {
	originalSettings      string
	vwoBuilder            *vwoBuilder
	isSettingsValid       bool
	settingsInvalidReason string
	settings              *settingsModel.Settings
}

// newVWOClient creates a new VWOClient instance
func newVWOClient(settings *settingsModel.Settings, vwoBuilder *vwoBuilder) *VWOClient {
	// Create service container
	serviceContainer := core.NewServiceContainer(
		fmt.Sprintf("%d_%s", vwoBuilder.options.AccountID, vwoBuilder.options.SDKKey),
		vwoBuilder.logManager,
		vwoBuilder.settingsManager,
		vwoBuilder.options,
		vwoBuilder.batchEventQueue,
		settings,
		vwoBuilder.networkManager,
	)

	// Validate settings once during initialization
	isSettingsValid, settingsInvalidReason := validateSettingsInternal(vwoBuilder.originalSettings, settings, vwoBuilder)
	if isSettingsValid {
		serviceContainer.GetSettingsManager().SetSettingsValidOnInit(true)
	}
	vwoClient := &VWOClient{
		vwoBuilder:            vwoBuilder,
		originalSettings:      vwoBuilder.originalSettings,
		isSettingsValid:       isSettingsValid,
		settingsInvalidReason: settingsInvalidReason,
		settings:              settings,
	}

	if !isSettingsValid {
		return vwoClient
	}

	// Set settings for batch event queue
	if vwoClient.vwoBuilder.batchEventQueue.IsInitialized() {
		vwoClient.vwoBuilder.batchEventQueue.SetSettings(settings)
	}
	// send init and usage stats events
	sendInitAndUsageStatsEvents(serviceContainer)

	// Process settings before creating the client
	// This sets variation allocation, adds linked campaigns, and gateway service flags
	utils.ProcessSettings(settings, vwoBuilder.logManager)

	vwoBuilder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["CLIENT_INITIALIZED"], map[string]interface{}{
		"sdkName":    constants.SDKName,
		"sdkVersion": constants.SDKVersion,
	}))
	return vwoClient
}

// sendInitAndUsageStatsEvents sends the SDK init and usage stats events to vwo
func sendInitAndUsageStatsEvents(serviceContainer interfaces.ServiceContainerInterface) {
	// Convert context to VWOContext
	contextModel := user.NewVWOContext(map[string]interface{}{
		enums.ContextID.GetValue(): fmt.Sprintf("%s_%s", serviceContainer.GetSettingsManager().GetAccountID(), serviceContainer.GetSettingsManager().GetSDKKey()),
	})

	// generate uuid for the user
	uuid := utils.GetUUID(contextModel.ID, serviceContainer.GetSettingsManager().GetAccountID())
	contextModel.SetUUID(uuid)
	settingsFetchTime := serviceContainer.GetSettingsManager().GetSettingsFetchTime()
	if serviceContainer.GetSettingsManager().GetIsSettingsProvidedInInit() {
		settingsFetchTime = 0
	}

	sdkInitTime := time.Now().UnixNano()/1e6 - serviceContainer.GetSettingsManager().GetStartTimeForInit()

	// if settings are valid and was initialized earlier is false, then send sdk init event
	sdkMetaInfo := serviceContainer.GetSettingsManager().GetSettingsObject().SDKMetaInfo
	var wasInitializedEarlier bool
	if sdkMetaInfo != nil {
		if val, ok := sdkMetaInfo["wasInitializedEarlier"].(bool); ok {
			wasInitializedEarlier = val
		}
	}
	if serviceContainer.GetSettingsManager().GetIsSettingsValidOnInit() && !wasInitializedEarlier {
		// Send SDK init event
		utils.SendSDKInitEvent(serviceContainer, contextModel, int(settingsFetchTime), int(sdkInitTime))

		// get usage stats account id
		usageStatsAccountId := serviceContainer.GetSettings().GetUsageStatsAccountID()
		if usageStatsAccountId != 0 {
			// Send SDK usage stats event
			utils.SendSDKUsageStatsEvent(serviceContainer, contextModel, usageStatsAccountId)
		}
	}
}

// GetFlag retrieves a feature flag for a given feature key and context
func (client *VWOClient) GetFlag(featureKey string, context map[string]interface{}) (flag models.GetFlagResponse, err error) {
	apiName := enums.ApiGetFlag

	// handle panic and return default fallback values
	defer func() {
		if r := recover(); r != nil {
			client.vwoBuilder.logManager.Error("EXECUTION_FAILED", map[string]interface{}{
				"apiName": apiName,
				"err":     fmt.Sprintf("Error in GetFlag: %v", r),
			}, map[string]interface{}{"an": apiName})
			// Return default fallback values
			flag = &models.GetFlag{Enabled: false}
			err = fmt.Errorf("panic recovered in GetFlag: %v", r)
		}
	}()

	client.vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["API_CALLED"], map[string]interface{}{
		"apiName": apiName,
	}))

	// Validate featureKey
	if featureKey == "" {
		client.vwoBuilder.logManager.Error("INVALID_PARAM", map[string]interface{}{
			"apiName":     apiName,
			"key":         "featureKey",
			"type":        "empty string",
			"correctType": "non-empty string",
		}, map[string]interface{}{"an": apiName})
		return &models.GetFlag{Enabled: false}, fmt.Errorf("featureKey should be a non-empty string")
	}

	// Validate context
	if context == nil || context[enums.ContextID.GetValue()] == nil || context[enums.ContextID.GetValue()] == "" {
		client.vwoBuilder.logManager.Error("INVALID_CONTEXT", nil, map[string]interface{}{"an": apiName})
		return &models.GetFlag{Enabled: false}, fmt.Errorf("invalid context")
	}

	// Validate settings
	if !client.isSettingsValid {
		client.vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
			"errors":    client.settingsInvalidReason,
			"accountId": strconv.Itoa(client.vwoBuilder.options.AccountID),
			"sdkKey":    client.vwoBuilder.options.SDKKey,
			"settings":  client.originalSettings,
		}, map[string]interface{}{"an": apiName})
		return &models.GetFlag{Enabled: false}, fmt.Errorf(client.settingsInvalidReason)
	}

	// Convert context to VWOContext
	contextModel := user.NewVWOContext(context)

	// generate uuid for the user
	uuid := utils.GetUUID(contextModel.ID, strconv.Itoa(client.vwoBuilder.options.AccountID))
	contextModel.SetUUID(uuid)

	// Create service container
	serviceContainer := core.NewServiceContainer(
		contextModel.ID,
		client.vwoBuilder.logManager,
		client.vwoBuilder.settingsManager,
		client.vwoBuilder.options,
		client.vwoBuilder.batchEventQueue,
		client.settings,
		client.vwoBuilder.networkManager,
	)

	// Get flag using API
	flag = api.GetFlag(featureKey, contextModel, serviceContainer)
	if flag == nil {
		return &models.GetFlag{Enabled: false}, fmt.Errorf("failed to get flag")
	}

	return flag, nil
}

// TrackEvent tracks an event with specified properties and context and returns true if the event is tracked successfully
func (client *VWOClient) TrackEvent(eventName string, context map[string]interface{}, eventProperties ...map[string]interface{}) (result map[string]bool, err error) {
	apiName := enums.ApiTrackEvent

	// handle panic and return default fallback values
	defer func() {
		if r := recover(); r != nil {
			client.vwoBuilder.logManager.Error("EXECUTION_FAILED", map[string]interface{}{
				"apiName": apiName,
				"err":     fmt.Sprintf("Error in TrackEvent: %v", r),
			}, map[string]interface{}{"an": apiName})
			// Return default fallback values
			result = map[string]bool{eventName: false}
			err = fmt.Errorf("panic recovered in TrackEvent: %v", r)
		}
	}()

	client.vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["API_CALLED"], map[string]interface{}{
		"apiName": apiName,
	}))

	// Validate eventName
	if eventName == "" {
		client.vwoBuilder.logManager.Error("INVALID_PARAM", map[string]interface{}{
			"apiName":     apiName,
			"key":         "eventName",
			"type":        "empty string",
			"correctType": "non-empty string",
		}, map[string]interface{}{"an": apiName})
		return map[string]bool{eventName: false}, fmt.Errorf("eventName should be a non-empty string")
	}

	// Validate context
	if context == nil || context[enums.ContextID.GetValue()] == nil || context[enums.ContextID.GetValue()] == "" {
		client.vwoBuilder.logManager.Error("INVALID_CONTEXT", nil, map[string]interface{}{"an": apiName})
		return map[string]bool{eventName: false}, fmt.Errorf("invalid context")
	}

	// Validate settings
	if !client.isSettingsValid {
		client.vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
			"errors":    client.settingsInvalidReason,
			"accountId": strconv.Itoa(client.vwoBuilder.options.AccountID),
			"sdkKey":    client.vwoBuilder.options.SDKKey,
			"settings":  client.originalSettings,
		}, map[string]interface{}{"an": apiName})
		return map[string]bool{eventName: false}, fmt.Errorf(client.settingsInvalidReason)
	}

	// Convert context to VWOContext
	contextModel := user.NewVWOContext(context)

	// generate uuid for the user
	uuid := utils.GetUUID(contextModel.ID, strconv.Itoa(client.vwoBuilder.options.AccountID))
	contextModel.SetUUID(uuid)

	// Create service container
	serviceContainer := core.NewServiceContainer(
		contextModel.ID,
		client.vwoBuilder.logManager,
		client.vwoBuilder.settingsManager,
		client.vwoBuilder.options,
		client.vwoBuilder.batchEventQueue,
		client.settings,
		client.vwoBuilder.networkManager,
	)

	// Track event using API
	var eventPropertiesMap map[string]interface{}
	if len(eventProperties) > 0 {
		eventPropertiesMap = eventProperties[0]
	}
	success := api.TrackEvent(eventName, contextModel, eventPropertiesMap, serviceContainer)
	result = map[string]bool{eventName: success}

	return result, nil
}

// SetAttribute sets multiple attributes for a user and sends an impression to vwo
func (client *VWOClient) SetAttribute(attributes map[string]interface{}, context map[string]interface{}) error {
	apiName := enums.ApiSetAttribute

	// handle panic and return error
	defer func() {
		if r := recover(); r != nil {
			client.vwoBuilder.logManager.Error("EXECUTION_FAILED", map[string]interface{}{
				"apiName": apiName,
				"err":     fmt.Sprintf("Error in SetAttributes: %v", r),
			}, map[string]interface{}{"an": apiName})
		}
	}()

	client.vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["API_CALLED"], map[string]interface{}{
		"apiName": apiName,
	}))

	// Validate attributes
	if len(attributes) == 0 {
		client.vwoBuilder.logManager.Error("ATTRIBUTES_MAP_ERROR", nil, map[string]interface{}{"an": apiName})
		return fmt.Errorf("attributes map should contain at least 1 key-value pair")
	}

	// Validate each attribute value type
	for key, value := range attributes {
		switch value.(type) {
		case bool, string, int, int64, float64:
			// Valid types
		default:
			client.vwoBuilder.logManager.Error("INVALID_PARAM", map[string]interface{}{
				"apiName":     apiName,
				"key":         key,
				"type":        fmt.Sprintf("%T", value),
				"correctType": "boolean, string or number",
			}, map[string]interface{}{"an": apiName})
			return fmt.Errorf("invalid attribute type for key %s", key)
		}
	}

	// Validate context
	if context == nil || context[enums.ContextID.GetValue()] == nil || context[enums.ContextID.GetValue()] == "" {
		client.vwoBuilder.logManager.Error("INVALID_CONTEXT", nil, map[string]interface{}{"an": apiName})
		return fmt.Errorf("invalid context")
	}

	// Validate settings
	if !client.isSettingsValid {
		client.vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
			"errors":    client.settingsInvalidReason,
			"accountId": strconv.Itoa(client.vwoBuilder.options.AccountID),
			"sdkKey":    client.vwoBuilder.options.SDKKey,
			"settings":  client.originalSettings,
		}, map[string]interface{}{"an": apiName})
		return fmt.Errorf(client.settingsInvalidReason)
	}

	// Convert context to VWOContext
	contextModel := user.NewVWOContext(context)

	// generate uuid for the user
	uuid := utils.GetUUID(contextModel.ID, strconv.Itoa(client.vwoBuilder.options.AccountID))
	contextModel.SetUUID(uuid)

	// Create service container
	serviceContainer := core.NewServiceContainer(
		contextModel.ID,
		client.vwoBuilder.logManager,
		client.vwoBuilder.settingsManager,
		client.vwoBuilder.options,
		client.vwoBuilder.batchEventQueue,
		client.settings,
		client.vwoBuilder.networkManager,
	)
	// Set attributes using API
	api.SetAttribute(attributes, contextModel, serviceContainer)

	return nil
}

// UpdateSettings updates the settings by fetching from the VWO server
func (client *VWOClient) UpdateSettings(options ...interface{}) (err error) {
	apiName := enums.ApiUpdateSettings

	// Defaults
	settings := ""
	isViaWebhook := true

	// Flexible options handling
	if len(options) > 0 {
		switch v := options[0].(type) {
		case string:
			settings = v
		case bool:
			isViaWebhook = v
		}
	}
	// Optional second argument as bool for isViaWebhook
	if len(options) > 1 {
		if b, ok := options[1].(bool); ok {
			isViaWebhook = b
		}
	}

	// handle panic and return error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error updating settings: %v", r)
			client.vwoBuilder.logManager.Error("UPDATING_CLIENT_INSTANCE_FAILED_WHEN_WEBHOOK_TRIGGERED", map[string]interface{}{
				"apiName":      apiName,
				"isViaWebhook": isViaWebhook,
				"err":          fmt.Sprintf("%v", r),
			}, map[string]interface{}{"an": apiName})
		}
	}()

	client.vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["API_CALLED"], map[string]interface{}{
		"apiName": apiName,
	}))

	// Fetch settings from server if not provided
	var newSettings string
	if settings == "" {
		// Fetch settings from server
		newSettingsData, err := client.vwoBuilder.settingsManager.FetchSettings(isViaWebhook)
		if err != nil {
			client.vwoBuilder.logManager.Error("UPDATING_CLIENT_INSTANCE_FAILED_WHEN_WEBHOOK_TRIGGERED", map[string]interface{}{
				"apiName":      apiName,
				"isViaWebhook": isViaWebhook,
				"err":          err.Error(),
			}, map[string]interface{}{"an": apiName})
			return err
		}
		newSettings = newSettingsData
	} else {
		newSettings = settings
	}

	err = client.updateSettingsInternal(newSettings)
	if err != nil {
		return err
	}
	client.vwoBuilder.originalSettings = newSettings
	client.vwoBuilder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["SETTINGS_UPDATED"], map[string]interface{}{
		"apiName":      apiName,
		"isViaWebhook": isViaWebhook,
	}))

	return nil
}

// GetOriginalSettings returns the original settings as map
func (client *VWOClient) GetOriginalSettings() string {
	return client.originalSettings
}

// updateSettingsInternal updates the settings internally (called from polling)
func (client *VWOClient) updateSettingsInternal(settingsJSON string) (err error) {
	if settingsJSON == "" {
		return fmt.Errorf("settings string is empty")
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error unmarshalling settings: %v", r)
			client.vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
				"errors":    err.Error(),
				"accountId": strconv.Itoa(client.vwoBuilder.options.AccountID),
				"sdkKey":    client.vwoBuilder.options.SDKKey,
				"settings":  settingsJSON,
			}, map[string]interface{}{"an": enums.ApiUpdateSettings})
		}
	}()

	// Parse the settings JSON
	var newSettings settingsModel.Settings
	err = json.Unmarshal([]byte(settingsJSON), &newSettings)
	if err != nil {
		client.vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
			"errors":    err.Error(),
			"accountId": strconv.Itoa(client.vwoBuilder.options.AccountID),
			"sdkKey":    client.vwoBuilder.options.SDKKey,
			"settings":  settingsJSON,
		}, map[string]interface{}{"an": enums.ApiUpdateSettings})
		return err
	}

	// Re-validate settings after update
	isSettingsValid, settingsInvalidReason := validateSettingsInternal(settingsJSON, &newSettings, client.vwoBuilder)
	if !isSettingsValid {
		return fmt.Errorf("settings are invalid: %v", settingsInvalidReason)
	}
	client.isSettingsValid = isSettingsValid
	client.settingsInvalidReason = settingsInvalidReason
	client.settings = &newSettings
	client.originalSettings = settingsJSON

	// Process the new settings
	utils.ProcessSettings(&newSettings, client.vwoBuilder.logManager)
	return nil
}

// validateSettingsInternal validates the settings using the settings schema (internal function)
func validateSettingsInternal(settingJSON string, processedSettings *settingsModel.Settings, vwoBuilder *vwoBuilder) (bool, string) {
	defer func() {
		if r := recover(); r != nil {
			vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
				"errors":    fmt.Sprintf("Error validating settings: %v", r),
				"accountId": strconv.Itoa(vwoBuilder.options.AccountID),
				"sdkKey":    vwoBuilder.options.SDKKey,
				"settings":  "null",
			}, map[string]interface{}{"an": enums.ApiUpdateSettings})
		}
	}()
	if processedSettings == nil {
		vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
			"errors":    "Settings object is null",
			"accountId": strconv.Itoa(vwoBuilder.options.AccountID),
			"sdkKey":    vwoBuilder.options.SDKKey,
			"settings":  "null",
		}, map[string]interface{}{"an": enums.ApiUpdateSettings})
		return false, "Settings object is null"
	}

	// Create a new settings schema validator
	settingsSchema := schemas.NewSettingsSchema()
	validationResult := settingsSchema.ValidateSettings(processedSettings)

	if !validationResult.IsValid() {
		vwoBuilder.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
			"errors":    validationResult.GetErrorsAsString(),
			"accountId": strconv.Itoa(vwoBuilder.options.AccountID),
			"sdkKey":    vwoBuilder.options.SDKKey,
			"settings":  settingJSON,
		}, map[string]interface{}{"an": enums.ApiUpdateSettings})
		return false, validationResult.GetErrorsAsString()
	}

	return true, ""
}

// FlushEvents flushes the events in the batch event queue
func (client *VWOClient) FlushEvents() (err error) {
	apiName := enums.ApiFlushEvents

	client.vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["API_CALLED"], map[string]interface{}{
		"apiName": apiName,
	}))

	// handle panic and return error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error flushing events: %v", r)
		}
	}()

	if client.vwoBuilder.batchEventQueue.IsInitialized() {
		client.vwoBuilder.batchEventQueue.FlushAndClearInterval()
	} else {
		client.vwoBuilder.logManager.Error("BATCHING_NOT_ENABLED", nil, map[string]interface{}{"an": apiName})
		err = fmt.Errorf("batching is not enabled")
		return err
	}

	return nil
}
