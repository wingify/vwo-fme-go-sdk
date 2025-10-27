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

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	settingsModel "github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/core"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/manager"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/storage"
	"github.com/wingify/vwo-fme-go-sdk/pkg/services"
)

// vwoBuilder handles the construction of VWO client instances
type vwoBuilder struct {
	options                           *models.VWOInitOptions
	logManager                        interfaces.LoggerServiceInterface
	settingsManager                   *services.SettingsManager
	settings                          *settingsModel.Settings
	batchEventQueue                   *services.BatchEventQueue
	vwoClient                         *VWOClient
	originalSettings                  string
	isSettingsFetchInProgress         bool
	isValidPollIntervalPassedFromInit bool
	pollingStopChan                   chan bool
	networkManager                    *manager.NetworkManager
}

// SetLogger sets up the logger service
func (vwoBuilder *vwoBuilder) SetLogger() *vwoBuilder {
	if vwoBuilder.options != nil && vwoBuilder.options.Logger != nil {
		// Pass the logger config directly to NewLogManager
		// NewLogManager will handle defaults if config is nil or missing values
		vwoBuilder.logManager = core.NewLogManager(vwoBuilder.options.Logger)
		vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Logger",
		}))
		return vwoBuilder
	} else {
		vwoBuilder.logManager = core.NewLogManager(nil)
		vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Logger",
		}))
		return vwoBuilder
	}
}

// SetSettingsManager sets up the settings manager
func (vwoBuilder *vwoBuilder) SetSettingsManager() *vwoBuilder {
	vwoBuilder.settingsManager = services.NewSettingsManager(vwoBuilder.options, vwoBuilder.logManager)
	vwoBuilder.logManager.SetSettingsManager(vwoBuilder.settingsManager)
	return vwoBuilder
}

// SetNetworkManager sets up the network manager
func (vwoBuilder *vwoBuilder) SetNetworkManager() *vwoBuilder {
	// Network manager is a singleton, just attach default client
	vwoBuilder.networkManager = &manager.NetworkManager{}

	// Use retry configuration if provided
	if vwoBuilder.options != nil && vwoBuilder.options.RetryConfig != nil {
		vwoBuilder.networkManager.AttachDefaultClientWithRetry(vwoBuilder.options.RetryConfig)
		vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Network Manager with Retry",
			"retryConfig": map[string]interface{}{
				"shouldRetry":       vwoBuilder.options.RetryConfig.ShouldRetry,
				"maxRetries":        vwoBuilder.options.RetryConfig.MaxRetries,
				"initialDelay":      vwoBuilder.options.RetryConfig.InitialDelay,
				"backoffMultiplier": vwoBuilder.options.RetryConfig.BackoffMultiplier,
			},
		}))
	} else {
		vwoBuilder.networkManager.AttachDefaultClient()
		vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Network Manager",
		}))
	}

	vwoBuilder.settingsManager.SetNetworkManager(vwoBuilder.networkManager)
	return vwoBuilder
}

func (vwoBuilder *vwoBuilder) InitBatching() *vwoBuilder {
	// Check if batch event data is provided in options
	if vwoBuilder.options.BatchEventData != nil {
		// Check if gatewayService is provided and skip SDK batching if so
		if vwoBuilder.settingsManager != nil && vwoBuilder.settingsManager.GetIsGatewayServiceProvided() {
			vwoBuilder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["GATEWAY_AND_BATCH_EVENTS_CONFIG_MISMATCH"], nil))
			return vwoBuilder
		}
		batchEventData := models.NewBatchEventData(vwoBuilder.options.BatchEventData)
		eventsPerRequest := batchEventData.GetEventsPerRequest()
		requestTimeInterval := batchEventData.GetRequestTimeInterval()

		isEventsPerRequestValid := eventsPerRequest > 0 && eventsPerRequest <= constants.MaxEventsPerRequest
		isRequestTimeIntervalValid := requestTimeInterval > 0

		// Handle invalid data types for individual parameters
		if !isEventsPerRequestValid {
			vwoBuilder.logManager.Error("INVALID_EVENTS_PER_REQUEST_VALUE", nil, map[string]interface{}{"an": enums.ApiInit})
			eventsPerRequest = constants.DefaultEventsPerRequest // Use default if invalid
		}

		if !isRequestTimeIntervalValid {
			vwoBuilder.logManager.Error("INVALID_REQUEST_TIME_INTERVAL_VALUE", nil, map[string]interface{}{"an": enums.ApiInit})
			requestTimeInterval = constants.DefaultRequestTimeInterval // Use default if invalid
		}

		// Initialize BatchEventQueue for batching
		// Convert BatchFlushCallback to FlushCallback
		var flushCallback models.FlushCallback
		if batchEventData.GetFlushCallback() != nil {
			flushCallback = func(err string, events string) {
				// Convert string error to error type for BatchFlushCallback
				var errorObj error
				if err != "" {
					errorObj = fmt.Errorf("%s", err)
				}
				batchEventData.GetFlushCallback()(errorObj, events)
			}
		}

		vwoBuilder.batchEventQueue = services.NewBatchEventQueue(
			eventsPerRequest,
			requestTimeInterval,
			flushCallback,
			vwoBuilder.options.AccountID,
			vwoBuilder.options.SDKKey,
			vwoBuilder.logManager,
		)
		vwoBuilder.batchEventQueue.SetNetworkManager(vwoBuilder.networkManager)

		vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Batching",
		}))
	}

	return vwoBuilder
}

// SetStorage sets up the storage service
func (vwoBuilder *vwoBuilder) SetStorage() *vwoBuilder {
	if vwoBuilder.options != nil && vwoBuilder.options.Storage != nil {
		// Attach the storage connector to the singleton
		storageInstance := storage.GetInstance()
		storageInstance.AttachConnector(vwoBuilder.options.Storage)

		vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
			"service": "Storage",
		}))
	}
	return vwoBuilder
}

// InitPolling initializes the polling mechanism
func (vwoBuilder *vwoBuilder) InitPolling() *vwoBuilder {
	if vwoBuilder.options.PollInterval >= 1000 && vwoBuilder.options.PollInterval != 0 {
		// This is to check if the poll_interval passed in options is valid
		vwoBuilder.isValidPollIntervalPassedFromInit = true
		vwoBuilder.pollingStopChan = make(chan bool)
		go vwoBuilder.checkAndPoll()
		return vwoBuilder
	} else if vwoBuilder.options.PollInterval > 0 {
		// Only log error if poll_interval is present in options but invalid
		vwoBuilder.logManager.Error("INVALID_POLLING_CONFIGURATION", map[string]interface{}{
			"key":         "pollInterval",
			"correctType": "number",
		}, map[string]interface{}{"an": enums.ApiInit})
	}
	return vwoBuilder
}

// InitUsageStats initializes usage statistics (placeholder for future implementation)
func (vwoBuilder *vwoBuilder) InitUsageStats() *vwoBuilder {
	vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["SERVICE_INITIALIZED"], map[string]interface{}{
		"service": "Usage Stats",
	}))
	return vwoBuilder
}

// GetSettings fetches settings from VWO servers
func (vwoBuilder *vwoBuilder) GetSettings(forceFetch bool) string {
	settingsString := vwoBuilder.settingsManager.GetSettings(forceFetch)
	vwoBuilder.originalSettings = settingsString
	vwoBuilder.settings = vwoBuilder.settingsManager.GetSettingsObject()
	return settingsString
}

// Build creates and returns a VWOClient instance
func (vwoBuilder *vwoBuilder) Build(settingsData *settingsModel.Settings) *VWOClient {
	// Create VWO client using the newVWOClient function
	// This will process settings before creating the client
	vwoClient := newVWOClient(settingsData, vwoBuilder)

	// Set VWO client reference in builder
	vwoBuilder.vwoClient = vwoClient
	// If poll_interval is not present in options, set it to the pollInterval from settings
	vwoBuilder.updatePollIntervalAndCheckAndPoll(vwoBuilder.originalSettings, true)
	return vwoClient
}

// updatePollIntervalAndCheckAndPoll updates the poll interval from settings and starts polling if needed
func (vwoBuilder *vwoBuilder) updatePollIntervalAndCheckAndPoll(settingsJSON string, shouldCheckAndPoll bool) {
	// Only update the poll_interval if poll_interval is not valid or not present in options
	var processedSettings *settingsModel.Settings
	if settingsJSON != "" {
		err := json.Unmarshal([]byte(settingsJSON), &processedSettings)
		if err != nil {
			// Ignore error, processedSettings will be nil
			processedSettings = nil
		}
	}

	if !vwoBuilder.isValidPollIntervalPassedFromInit && processedSettings != nil {
		vwoBuilder.options.PollInterval = processedSettings.GetPollInterval()

		if processedSettings.GetPollInterval() == constants.DefaultPollInterval {
			vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["USING_POLL_INTERVAL_FROM_SETTINGS"], map[string]interface{}{
				"source":       "default",
				"pollInterval": strconv.Itoa(constants.DefaultPollInterval),
			}))
		} else {
			vwoBuilder.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["USING_POLL_INTERVAL_FROM_SETTINGS"], map[string]interface{}{
				"source":       "settings",
				"pollInterval": strconv.Itoa(vwoBuilder.options.PollInterval),
			}))
		}
	}

	if shouldCheckAndPoll && !vwoBuilder.isValidPollIntervalPassedFromInit && processedSettings != nil && vwoBuilder.options.PollInterval >= 1000 {
		vwoBuilder.pollingStopChan = make(chan bool)
		go vwoBuilder.checkAndPoll()
	}
}

// checkAndPoll checks for settings updates at the configured interval
func (vwoBuilder *vwoBuilder) checkAndPoll() {
	ticker := time.NewTicker(time.Duration(vwoBuilder.options.PollInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			latestSettings := vwoBuilder.fetchSettings(true)
			if vwoBuilder.originalSettings != "" && latestSettings != "" {
				// Compare settings
				if !vwoBuilder.areSettingsEqual(vwoBuilder.originalSettings, latestSettings) {
					vwoBuilder.updateSettingsOnBuilder(latestSettings)
				} else {
					vwoBuilder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["POLLING_NO_CHANGE_IN_SETTINGS"], map[string]interface{}{}))
				}
			} else if vwoBuilder.originalSettings == "" && latestSettings != "" {
				vwoBuilder.updateSettingsOnBuilder(latestSettings)
			}

		case <-vwoBuilder.pollingStopChan:
			return
		}
	}
}

// fetchSettings fetches settings from VWO servers
func (vwoBuilder *vwoBuilder) fetchSettings(forceFetch bool) string {
	// Check if a fetch operation is already in progress
	if vwoBuilder.isSettingsFetchInProgress || vwoBuilder.settingsManager == nil {
		return ""
	}

	apiName := string(enums.ApiInit)
	if forceFetch {
		apiName = constants.POLLING
	}

	// Set the flag to indicate that a fetch operation is in progress
	vwoBuilder.isSettingsFetchInProgress = true
	defer func() {
		if r := recover(); r != nil {
			vwoBuilder.logManager.Error("ERROR_FETCHING_SETTINGS", map[string]interface{}{
				"err": r,
			}, map[string]interface{}{"an": apiName})
			vwoBuilder.isSettingsFetchInProgress = false
		}
	}()

	// Retrieve the settings
	settingsString := vwoBuilder.settingsManager.GetSettings(forceFetch)

	if !forceFetch {
		// Store the original settings
		vwoBuilder.originalSettings = settingsString
	}
	vwoBuilder.isSettingsFetchInProgress = false
	return settingsString
}

// areSettingsEqual compares two settings JSON strings
func (vwoBuilder *vwoBuilder) areSettingsEqual(settings1, settings2 string) bool {
	var obj1, obj2 interface{}

	if err := json.Unmarshal([]byte(settings1), &obj1); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(settings2), &obj2); err != nil {
		return false
	}

	// Deep comparison using JSON marshaling
	json1, _ := json.Marshal(obj1)
	json2, _ := json.Marshal(obj2)

	return string(json1) == string(json2)
}

// updateSettingsOnBuilder updates the settings on the VWOBuilder instance
func (vwoBuilder *vwoBuilder) updateSettingsOnBuilder(latestSettings string) {
	if vwoBuilder.vwoClient != nil {
		err := vwoBuilder.vwoClient.updateSettingsInternal(latestSettings)
		if err == nil {
			vwoBuilder.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["POLLING_SET_SETTINGS"], map[string]interface{}{}))
			vwoBuilder.originalSettings = latestSettings
			vwoBuilder.updatePollIntervalAndCheckAndPoll(vwoBuilder.originalSettings, false)
		} else {
			vwoBuilder.logManager.Error("ERROR_UPDATING_SETTINGS", map[string]interface{}{
				"err":              err.Error(),
				"originalSettings": vwoBuilder.originalSettings,
				"latestSettings":   latestSettings,
			}, map[string]interface{}{"an": constants.POLLING})
		}
	}
}

// StopPolling stops the polling goroutine
func (vwoBuilder *vwoBuilder) StopPolling() {
	if vwoBuilder.pollingStopChan != nil {
		close(vwoBuilder.pollingStopChan)
	}
}
