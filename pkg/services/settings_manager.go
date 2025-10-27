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

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/schemas"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/manager"
	networkModels "github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/utils"
)

// SettingsManager handles settings fetching and management
type SettingsManager struct {
	SDKKey                   string
	AccountID                int
	Expiry                   int
	NetworkTimeout           int
	Hostname                 string
	Port                     int
	Protocol                 string
	IsGatewayServiceProvided bool
	IsSettingsProvidedInInit bool
	IsSettingsValidOnInit    bool
	StartTimeForInit         int64 // Time in milliseconds
	SettingsFetchTime        int64 // Time in milliseconds
	logManager               interfaces.LoggerServiceInterface
	settings                 *settings.Settings
	settingsString           string
	networkManager           *manager.NetworkManager
}

// NewSettingsManager creates a new settings manager
func NewSettingsManager(options *models.VWOInitOptions, logManager interfaces.LoggerServiceInterface) *SettingsManager {
	settingsManager := &SettingsManager{
		logManager:     logManager,
		SDKKey:         options.SDKKey,
		AccountID:      options.AccountID,
		Expiry:         int(constants.SettingsExpiry),
		NetworkTimeout: int(constants.SettingsTimeout),
		Protocol:       constants.HTTPSProtocol,
	}

	// Check if gateway service is provided
	if options.GatewayService != nil {
		settingsManager.IsGatewayServiceProvided = true
		settingsManager.parseGatewayService(options.GatewayService)
	} else {
		settingsManager.Hostname = constants.HostName
	}

	return settingsManager
}

// GetAccountID returns the account ID as string
func (settingsManager *SettingsManager) GetAccountID() string {
	return fmt.Sprintf("%d", settingsManager.AccountID)
}

// GetSDKKey returns the SDK key
func (settingsManager *SettingsManager) GetSDKKey() string {
	return settingsManager.SDKKey
}

// GetProtocol return the protocol
func (settingsManager *SettingsManager) GetProtocol() string {
	return settingsManager.Protocol
}

// GetHostname return the hostname
func (settingsManager *SettingsManager) GetHostname() string {
	return settingsManager.Hostname
}

// GetPort return the port
func (settingsManager *SettingsManager) GetPort() int {
	return settingsManager.Port
}

// parseGatewayService parses the gateway service configuration
func (settingsManager *SettingsManager) parseGatewayService(gatewayService map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			settingsManager.logManager.Error("ERROR_PARSING_GATEWAY_URL", map[string]interface{}{"err": fmt.Sprintf("%v", r)}, map[string]interface{}{"an": enums.ApiInit})
			settingsManager.Hostname = constants.HostName
		}
	}()

	gatewayURL, ok := gatewayService[enums.NetworkURL.GetValue()].(string)
	if !ok || gatewayURL == "" {
		settingsManager.Hostname = constants.HostName
		return
	}

	gatewayProtocol, _ := gatewayService[enums.NetworkProtocol.GetValue()].(string)
	gatewayPort, _ := gatewayService[enums.NetworkPort.GetValue()].(int)

	// Parse URL
	var parsedURL *url.URL
	var err error

	if len(gatewayURL) >= 7 && (gatewayURL[:7] == constants.HTTPProtocol+"://" || (len(gatewayURL) >= 8 && gatewayURL[:8] == constants.HTTPSProtocol+"://")) {
		parsedURL, err = url.Parse(gatewayURL)
	} else if gatewayProtocol != "" {
		parsedURL, err = url.Parse(gatewayProtocol + "://" + gatewayURL)
	} else {
		parsedURL, err = url.Parse(constants.HTTPSProtocol + "://" + gatewayURL)
	}

	if err != nil {
		settingsManager.logManager.Error("ERROR_PARSING_GATEWAY_URL", map[string]interface{}{"err": err.Error()}, map[string]interface{}{"an": enums.ApiInit})
		settingsManager.Hostname = constants.HostName
		return
	}

	settingsManager.Hostname = parsedURL.Hostname()
	settingsManager.Protocol = parsedURL.Scheme

	if parsedURL.Port() != "" {
		fmt.Sscanf(parsedURL.Port(), "%d", &settingsManager.Port)
	} else if gatewayPort != 0 {
		settingsManager.Port = gatewayPort
	}
}

// GetSettingsFetchTime gets the settings fetch time
func (settingsManager *SettingsManager) GetSettingsFetchTime() int64 {
	return settingsManager.SettingsFetchTime
}

// fetchSettingsAndCacheInStorage fetches settings from the server
func (settingsManager *SettingsManager) fetchSettingsAndCacheInStorage() string {
	settingsData, err := settingsManager.FetchSettings(false)
	if err != nil {
		settingsManager.logManager.Error("ERROR_FETCHING_SETTINGS", map[string]interface{}{
			"err":       err.Error(),
			"accountId": strconv.Itoa(settingsManager.AccountID),
			"sdkKey":    settingsManager.SDKKey,
		}, map[string]interface{}{"an": enums.ApiInit}, false)
		return ""
	}
	return settingsData
}

// FetchSettings fetches settings from the server
func (settingsManager *SettingsManager) FetchSettings(isViaWebhook bool) (string, error) {
	if settingsManager.SDKKey == "" || settingsManager.AccountID == 0 {
		return "", fmt.Errorf("SDK Key and Account ID are required to fetch settings. Aborting")
	}

	// Build query parameters
	queryParams := settingsManager.getSettingsPath(settingsManager.SDKKey, settingsManager.AccountID)
	queryParams[enums.NetworkAPIVersion.GetValue()] = "3"
	queryParams[enums.NetworkSDKName.GetValue()] = constants.SDKName
	queryParams[enums.NetworkSDKVersion.GetValue()] = constants.SDKVersion

	// Check development mode
	config := settingsManager.networkManager.GetConfig()
	if config == nil || !config.IsDevelopmentMode {
		queryParams["s"] = "prod"
	}

	// Determine endpoint
	apiName := enums.ApiInit
	endpoint := constants.SettingsEndpoint
	if isViaWebhook {
		apiName = enums.ApiUpdateSettings
		endpoint = constants.WebhookSettingsEndpoint
	}

	// Create request
	request := networkModels.NewRequestModel(
		settingsManager.Hostname,
		enums.HTTPMethodGET.GetValue(),
		endpoint,
		queryParams,
		nil,
		nil,
		settingsManager.Protocol,
		settingsManager.Port,
		"",
	)
	request.Timeout = settingsManager.NetworkTimeout

	// Track fetch time
	startTime := time.Now().UnixNano() / 1e6

	// Make network request
	response := settingsManager.networkManager.Get(request)

	if response == nil {
		return "", fmt.Errorf("network request failed: response is nil")
	}

	if response.TotalAttempts > 0 {
		lt := enums.LogLevelEnumInfo.GetValue()
		category := enums.DebuggerCategoryRetry.GetValue()
		message_type := constants.NETWORK_CALL_SUCCESS_WITH_RETRIES
		msg := log.BuildMessage(log.InfoLogMessagesEnum[message_type], map[string]interface{}{
			"extraData": endpoint,
			"attempts":  response.TotalAttempts,
			"err":       response.Error.Error(),
		})

		if response.StatusCode != 200 {
			category = enums.DebuggerCategoryNetwork.GetValue()
			message_type = constants.NETWORK_CALL_FAILURE_AFTER_MAX_RETRIES
			msg = log.BuildMessage(log.ErrorLogMessagesEnum[message_type], map[string]interface{}{
				"extraData": endpoint,
				"attempts":  response.TotalAttempts,
				"err":       response.Error.Error(),
			})
			lt = enums.LogLevelEnumError.GetValue()
		}

		// create debug event props
		debugEventProps := map[string]interface{}{
			enums.DebugPropCategory.GetValue():    category,
			enums.DebugPropAPI.GetValue():         string(apiName),
			enums.DebugPropMessage.GetValue():     msg,
			enums.DebugPropLogLevel.GetValue():    lt,
			enums.DebugPropMessageType.GetValue(): message_type,
		}

		// send debug event to vwo
		utils.SendDebugEventToVWO(settingsManager, debugEventProps)
	}

	if response.StatusCode != 200 {
		errMsg := "request failed with status code: " + strconv.Itoa(response.StatusCode)
		if response.Error != nil {
			errMsg = response.Error.Error()
		}

		return "", errors.New(errMsg)
	}

	settingsManager.SettingsFetchTime = time.Now().UnixNano()/1e6 - startTime

	// normalise settings data: if features or campaigns are empty objects, convert to empty arrays
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(response.Data), &raw); err == nil {
		if v, ok := raw["features"].(map[string]interface{}); ok && len(v) == 0 {
			raw["features"] = []interface{}{}
		}
		if v, ok := raw["campaigns"].(map[string]interface{}); ok && len(v) == 0 {
			raw["campaigns"] = []interface{}{}
		}
		if normalized, err := json.Marshal(raw); err == nil {
			return string(normalized), nil
		}
	}
	return response.Data, nil
}

// getSettingsPath creates the query parameters for the settings API
func (settingsManager *SettingsManager) getSettingsPath(sdkKey string, accountID int) map[string]string {
	randomNum := rand.Float64()
	return map[string]string{
		"i": sdkKey,
		"r": fmt.Sprintf("%.16f", randomNum),
		"a": fmt.Sprintf("%d", accountID),
	}
}

// GetSettings fetches settings from the server with optional validation
func (settingsManager *SettingsManager) GetSettings(forceFetch bool) string {
	apiName := string(enums.ApiInit)
	if forceFetch {
		apiName = constants.POLLING
		return settingsManager.fetchSettingsAndCacheInStorage()
	}

	settingsData := settingsManager.fetchSettingsAndCacheInStorage()
	if settingsData == "" {
		fmt.Println("Settings is null")
		return ""
	}

	// Parse and validate settings
	var settingsObj settings.Settings
	err := json.Unmarshal([]byte(settingsData), &settingsObj)
	if err != nil {
		settingsManager.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
			"errors":    fmt.Sprintf("Exception during parsing: %v", err),
			"accountId": strconv.Itoa(settingsManager.AccountID),
			"sdkKey":    settingsManager.SDKKey,
			"settings":  "null",
		}, map[string]interface{}{"an": apiName}, false)
		return ""
	}
	settingsManager.settings = &settingsObj

	// Validate settings using SettingsSchema
	validationResult := schemas.NewSettingsSchema().ValidateSettings(&settingsObj)
	if validationResult.IsValid() {
		settingsManager.IsSettingsValidOnInit = true
		settingsManager.settingsString = settingsData
		return settingsData
	}

	settingsManager.logManager.Error("INVALID_SETTINGS_SCHEMA", map[string]interface{}{
		"errors":    validationResult.GetErrorsAsString(),
		"accountId": strconv.Itoa(settingsManager.AccountID),
		"sdkKey":    settingsManager.SDKKey,
		"settings":  settingsData,
	}, map[string]interface{}{"an": apiName}, false)

	return settingsData
}

// GetSettingsObject returns the parsed settings object
func (settingsManager *SettingsManager) GetSettingsObject() *settings.Settings {
	return settingsManager.settings
}

// GetIsGatewayServiceProvided returns whether gateway service is provided
func (settingsManager *SettingsManager) GetIsGatewayServiceProvided() bool {
	return settingsManager.IsGatewayServiceProvided
}

// GetLoggerService returns the logger service
func (settingsManager *SettingsManager) GetLoggerService() interfaces.LoggerServiceInterface {
	return settingsManager.logManager
}

// SetSettings sets the settings object and string
func (settingsManager *SettingsManager) SetSettings(settingsObj *settings.Settings, settingsString string) {
	settingsManager.settings = settingsObj
	settingsManager.settingsString = settingsString
}

// SetSettingsValidOnInit sets the IsSettingsValidOnInit flag
func (settingsManager *SettingsManager) SetSettingsValidOnInit(valid bool) {
	settingsManager.IsSettingsValidOnInit = valid
}

func (settingsManager *SettingsManager) SetNetworkManager(networkManager *manager.NetworkManager) {
	settingsManager.networkManager = networkManager
}

func (settingsManager *SettingsManager) GetIsSettingsProvidedInInit() bool {
	return settingsManager.IsSettingsProvidedInInit
}

func (settingsManager *SettingsManager) GetStartTimeForInit() int64 {
	return settingsManager.StartTimeForInit
}

func (settingsManager *SettingsManager) GetIsSettingsValidOnInit() bool {
	return settingsManager.IsSettingsValidOnInit
}

func (settingsManager *SettingsManager) GetNetworkManager() *manager.NetworkManager {
	return settingsManager.networkManager
}
