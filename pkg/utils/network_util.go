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

package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/request"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	networkModels "github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"
)

// GetSettingsPath creates the query parameters for the settings API
func GetSettingsPath(apiKey string, accountID int) map[string]string {
	settingsQueryParams := request.NewSettingsQueryParams(apiKey, generateRandom(), fmt.Sprintf("%d", accountID))
	return settingsQueryParams.GetQueryParams()
}

// GetEventsBaseProperties creates the base properties for the event arch APIs
func GetEventsBaseProperties(settingsManager interfaces.SettingsManagerInterface, eventName string, visitorUserAgent string, ipAddress string) map[string]string {
	requestQueryParams := request.NewRequestQueryParams(
		eventName,
		settingsManager.GetAccountID(),
		settingsManager.GetSDKKey(),
		visitorUserAgent,
		ipAddress,
	)
	return requestQueryParams.GetQueryParams()
}

// GetEventBasePayload creates the base payload for the event arch APIs
func GetEventBasePayload(settingsManager interfaces.SettingsManagerInterface, userID string, eventName string, visitorUserAgent string, ipAddress string, usageStatsAccountId int) *request.EventArchPayload {
	accountID := settingsManager.GetAccountID()
	if usageStatsAccountId != 0 {
		accountID = fmt.Sprintf("%d", usageStatsAccountId)
	}
	uuid := GetUUID(userID, accountID)

	eventArchData := &request.EventArchData{
		MsgID:     GenerateMsgID(uuid),
		VisID:     uuid,
		SessionID: GenerateSessionID(),
		Event:     createEvent(settingsManager.GetSDKKey(), eventName),
		Visitor:   createVisitor(settingsManager.GetSDKKey()),
	}

	setOptionalVisitorData(eventArchData, visitorUserAgent, ipAddress)

	eventArchPayload := &request.EventArchPayload{
		D: eventArchData,
	}

	return eventArchPayload
}

// GetTrackUserPayloadData constructs payload data for tracking user
func GetTrackUserPayloadData(serviceContainer interfaces.ServiceContainerInterface, eventName string, campaignID int, variationID int, context *user.VWOContext) map[string]interface{} {
	properties := GetEventBasePayload(serviceContainer.GetSettingsManager(), context.GetID(), eventName, context.GetUserAgent(), context.GetIPAddress(), 0)
	properties.D.Event.Props.ID = campaignID
	properties.D.Event.Props.Variation = fmt.Sprintf("%d", variationID)
	properties.D.Event.Props.IsFirst = 1

	postSegmentationVariables := context.GetPostSegmentationVariables()
	customVariables := context.GetCustomVariables()

	// Add post-segmentation variables if they exist in custom variables
	if postSegmentationVariables != nil && len(postSegmentationVariables) > 0 && customVariables != nil && len(customVariables) > 0 {
		for _, variable := range postSegmentationVariables {
			if value, ok := customVariables[variable]; ok {
				properties.D.Visitor.Props[variable] = value
			}
		}
	}

	// Log impression for track user
	serviceContainer.GetLoggerService().Debug(log.BuildMessage(log.DebugLogMessagesEnum["IMPRESSION_FOR_TRACK_USER"], map[string]interface{}{
		"accountId":  serviceContainer.GetSettingsManager().GetAccountID(),
		"userId":     context.GetID(),
		"campaignId": fmt.Sprintf("%d", campaignID),
	}))

	payload := convertToMap(properties)
	return removeNullValues(payload)
}

// GetTrackGoalPayloadData constructs payload data for tracking goals/events
func GetTrackGoalPayloadData(serviceContainer interfaces.ServiceContainerInterface, userID string, eventName string, context *user.VWOContext, eventProperties map[string]interface{}) map[string]interface{} {
	properties := GetEventBasePayload(serviceContainer.GetSettingsManager(), userID, eventName, context.UserAgent, context.IPAddress, 0)
	properties.D.Event.Props.IsCustomEvent = true
	addCustomEventProperties(properties, eventProperties)

	// Log impression for track goal
	serviceContainer.GetLoggerService().Debug(log.BuildMessage(log.DebugLogMessagesEnum["IMPRESSION_FOR_TRACK_GOAL"], map[string]interface{}{
		"accountId": serviceContainer.GetSettingsManager().GetAccountID(),
		"userId":    userID,
		"eventName": eventName,
	}))

	payload := convertToMap(properties)
	return removeNullValues(payload)
}

// GetAttributePayloadData constructs payload data for setting attributes
func GetAttributePayloadData(serviceContainer interfaces.ServiceContainerInterface, userID string, eventName string, attributeMap map[string]interface{}) map[string]interface{} {
	properties := GetEventBasePayload(serviceContainer.GetSettingsManager(), userID, eventName, "", "", 0)
	properties.D.Event.Props.IsCustomEvent = true
	properties.D.Visitor.Props = attributeMap

	// Log impression for sync visitor properties
	serviceContainer.GetLoggerService().Debug(log.BuildMessage(log.DebugLogMessagesEnum["IMPRESSION_FOR_SYNC_VISITOR_PROP"], map[string]interface{}{
		"accountId": serviceContainer.GetSettingsManager().GetAccountID(),
		"userId":    userID,
		"eventName": eventName,
	}))

	payload := convertToMap(properties)
	return removeNullValues(payload)
}

// SendPostAPIRequest sends visitor and conversion events to vwo
func SendPostAPIRequest(serviceContainer interfaces.ServiceContainerInterface, properties map[string]string, payload map[string]interface{}, context *user.VWOContext, campaignInfo map[string]interface{}) {
	headers := createHeaders(context.GetUserAgent(), context.GetIPAddress())
	eventName := properties["en"]
	request := networkModels.NewRequestModel(
		serviceContainer.GetBaseUrl(),
		enums.ApiMethodPost.GetValue(),
		enums.Events.GetURL(),
		properties,
		payload,
		headers,
		serviceContainer.GetSettingsManager().GetProtocol(),
		serviceContainer.GetSettingsManager().GetPort(),
		eventName,
	)

	var apiName string
	var extraDataForMessage string
	if eventName == enums.VWOVariationShown.GetValue() {
		apiName = string(enums.ApiGetFlag)
		if campaignInfo != nil {
			if campaignType, ok := campaignInfo["campaignType"].(string); ok {
				if campaignType == enums.CampaignTypeRollout.GetValue() || campaignType == enums.CampaignTypePersonalize.GetValue() {
					if featureKey, ok := campaignInfo["featureKey"].(string); ok {
						if variationName, ok := campaignInfo["variationName"].(string); ok {
							extraDataForMessage = fmt.Sprintf("feature: %s, rule: %s", featureKey, variationName)
						}
					}
				} else {
					if featureKey, ok := campaignInfo["featureKey"].(string); ok {
						if campaignKey, ok := campaignInfo["campaignKey"].(string); ok {
							if variationName, ok := campaignInfo["variationName"].(string); ok {
								extraDataForMessage = fmt.Sprintf("feature: %s, rule: %s and variation: %s", featureKey, campaignKey, variationName)
							}
						}
					}
				}
			}
		}
	} else if eventName == enums.VWOSyncVisitorProp.GetValue() {
		apiName = string(enums.ApiSetAttribute)
		extraDataForMessage = apiName
	} else if eventName == enums.VWOSDKInitEvent.GetValue() {
		apiName = eventName
		extraDataForMessage = apiName
	} else if eventName == enums.VWOSDKUsageStats.GetValue() {
		apiName = eventName
		extraDataForMessage = apiName
	} else {
		apiName = string(enums.ApiTrackEvent)
		extraDataForMessage = fmt.Sprintf("event: %s", eventName)
	}

	go func() {
		response := serviceContainer.GetNetworkManager().Post(request, nil)
		if response != nil && response.TotalAttempts > 0 {
			lt := enums.LogLevelEnumInfo.GetValue()
			category := enums.DebuggerCategoryRetry.GetValue()
			message_type := constants.NETWORK_CALL_SUCCESS_WITH_RETRIES
			msg := log.BuildMessage(log.InfoLogMessagesEnum[message_type], map[string]interface{}{
				"extraData": extraDataForMessage,
				"attempts":  response.TotalAttempts,
				"err":       response.Error.Error(),
			})

			if response.StatusCode != 200 {
				category = enums.DebuggerCategoryNetwork.GetValue()
				message_type = constants.NETWORK_CALL_FAILURE_AFTER_MAX_RETRIES
				msg = log.BuildMessage(log.ErrorLogMessagesEnum[message_type], map[string]interface{}{
					"extraData": extraDataForMessage,
					"attempts":  response.TotalAttempts,
					"err":       response.Error.Error(),
				})
				lt = enums.LogLevelEnumError.GetValue()
			}

			// create debug event props
			debugEventProps := map[string]interface{}{
				enums.DebugPropCategory.GetValue():    category,
				enums.DebugPropAPI.GetValue():         apiName,
				enums.DebugPropMessage.GetValue():     msg,
				enums.DebugPropLogLevel.GetValue():    lt,
				enums.DebugPropMessageType.GetValue(): message_type,
			}
			// send debug event to vwo
			SendDebugEventToVWO(serviceContainer.GetSettingsManager(), debugEventProps)
		}
		if response != nil && response.StatusCode >= 200 && response.StatusCode < 300 {
			serviceContainer.GetLoggerService().Debug(log.BuildMessage(log.DebugLogMessagesEnum["NETWORK_CALL_SUCCESS"], map[string]interface{}{
				"event":     eventName,
				"endPoint":  enums.Events.GetURL(),
				"accountId": serviceContainer.GetSettingsManager().GetAccountID(),
				"userId":    context.GetID(),
				"uuid":      payload["d"].(map[string]interface{})["visId"],
			}))
		} else {
			// Log network call exception
			errMsg := ""
			if response != nil {
				errMsg = response.Error.Error()
			}
			if eventName != enums.VWODebuggerEvent.GetValue() && eventName != enums.VWOSDKInitEvent.GetValue() && eventName != enums.VWOSDKUsageStats.GetValue() {
				serviceContainer.GetLoggerService().Error("NETWORK_CALL_EXCEPTION", map[string]interface{}{
					"extraData": extraDataForMessage,
					"accountId": serviceContainer.GetSettingsManager().GetAccountID(),
					"err":       errMsg,
				}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
			}
		}
	}()
}

// SendEventDirectlyToDACDN sends event directly to DACDN
func SendEventDirectlyToDACDN(settingsManager interfaces.SettingsManagerInterface, properties map[string]string, payload map[string]interface{}, eventName string) {
	headers := createHeaders("", "")

	request := networkModels.NewRequestModel(
		constants.HostName,
		enums.ApiMethodPost.GetValue(),
		enums.Events.GetURL(),
		properties,
		payload,
		headers,
		constants.HTTPSProtocol,
		0,
		eventName,
	)

	go func() {
		response := settingsManager.GetNetworkManager().Post(request, nil)
		if response == nil || (response.StatusCode < 200 || response.StatusCode >= 300) {
			if eventName != enums.VWODebuggerEvent.GetValue() {
				// Log network call exception
				settingsManager.GetLoggerService().Error("NETWORK_CALL_EXCEPTION", map[string]interface{}{
					"extraData": fmt.Sprintf("event: %s", eventName),
					"accountId": settingsManager.GetAccountID(),
					"err":       response.Error.Error(),
				}, nil)
			}
		}
	}()
}

// GetSDKInitEventPayload creates payload for SDK init event
func GetSDKInitEventPayload(settingsManager interfaces.SettingsManagerInterface, userID string, eventName string, settingsFetchTime int, sdkInitTime int) map[string]interface{} {
	properties := GetEventBasePayload(settingsManager, userID, eventName, "", "", 0)

	properties.D.Event.Props.EnvKey = settingsManager.GetSDKKey()
	properties.D.Event.Props.Product = constants.FME

	data := map[string]interface{}{
		enums.SDKInitPayloadIsSDKInitialized.GetValue():  true,
		enums.SDKInitPayloadSettingsFetchTime.GetValue(): settingsFetchTime,
		enums.SDKInitPayloadSDKInitTime.GetValue():       sdkInitTime,
	}

	properties.D.Event.Props.Data = data

	return removeNullValues(convertToMap(properties))
}

// GetSDKUsageStatsEventPayload creates payload for SDK usage stats event
func GetSDKUsageStatsEventPayload(settingsManager interfaces.SettingsManagerInterface, userID string, eventName string, usageStatsAccountId int, usageStatsData map[string]interface{}) map[string]interface{} {
	properties := GetEventBasePayload(settingsManager, userID, eventName, "", "", usageStatsAccountId)
	properties.D.Event.Props.Product = constants.FME
	properties.D.Event.Props.VWOMeta = usageStatsData

	return removeNullValues(convertToMap(properties))
}

// GetDebuggerEventPayload creates payload for debugger event API
func GetDebuggerEventPayload(settingsManager interfaces.SettingsManagerInterface, eventProps map[string]interface{}) map[string]interface{} {
	userID := fmt.Sprintf("%s_%s", settingsManager.GetAccountID(), settingsManager.GetSDKKey())
	properties := GetEventBasePayload(settingsManager, userID, enums.VWODebuggerEvent.GetValue(), "", "", 0)

	// check if eventProps has valid uuid
	if uuid, ok := eventProps[enums.DebugPropUUID.GetValue()].(string); ok && uuid != "" {
		properties.D.MsgID = GenerateMsgID(uuid)
		properties.D.VisID = uuid
	} else {
		// set uuid in eventProps
		eventProps[enums.DebugPropUUID.GetValue()] = properties.D.VisID
	}

	// check if eventProps has valid sessionId
	if sessionId, ok := eventProps[enums.DebugPropSessionID.GetValue()].(int64); ok && sessionId != 0 {
		properties.D.SessionID = sessionId
	} else {
		// set sessionId in eventProps
		eventProps[enums.DebugPropSessionID.GetValue()] = properties.D.SessionID
	}

	eventProps[enums.DebugPropAccountID.GetValue()] = settingsManager.GetAccountID()
	eventProps[enums.DebugPropProduct.GetValue()] = constants.FME
	eventProps[enums.DebugPropSDKName.GetValue()] = constants.SDKName
	eventProps[enums.DebugPropSDKVersion.GetValue()] = constants.SDKVersion
	eventProps[enums.DebugPropEventID.GetValue()] = GetRandomUUID(settingsManager.GetSDKKey())

	properties.D.Event.Props = &request.Props{}
	properties.D.Event.Props.VWOMeta = eventProps

	payload := convertToMap(properties)
	return RemoveNullValues(payload)
}

// RemoveNullValues removes null values from the map recursively
func RemoveNullValues(originalMap map[string]interface{}) map[string]interface{} {
	cleanedMap := make(map[string]interface{})

	for key, value := range originalMap {
		if valueMap, ok := value.(map[string]interface{}); ok {
			// Recursively remove null values from nested maps
			value = RemoveNullValues(valueMap)
		}
		if value != nil {
			cleanedMap[key] = value
		}
	}

	return cleanedMap
}

// GenerateMsgID generates a message ID
func GenerateMsgID(uuid string) string {
	return fmt.Sprintf("%s-%d", uuid, time.Now().UnixNano()/1e6)
}

// GenerateSessionID generates a session ID
func GenerateSessionID() int64 {
	return time.Now().Unix()
}

// GenerateRandom generates a random number as string
func GenerateRandom() string {
	return fmt.Sprintf("%.16f", rand.Float64())
}

// CreateHeaders creates headers for the request
func CreateHeaders(userAgent string, ipAddress string) map[string]string {
	headers := make(map[string]string)
	if userAgent != "" {
		headers[enums.UserAgent.GetHeader()] = userAgent
	}
	if ipAddress != "" {
		headers[enums.IP.GetHeader()] = ipAddress
	}
	return headers
}

// Helper functions for internal use
func generateRandom() string {
	return GenerateRandom()
}

func removeNullValues(originalMap map[string]interface{}) map[string]interface{} {
	return RemoveNullValues(originalMap)
}

func createHeaders(userAgent string, ipAddress string) map[string]string {
	return CreateHeaders(userAgent, ipAddress)
}

func setOptionalVisitorData(eventArchData *request.EventArchData, visitorUserAgent string, ipAddress string) {
	if visitorUserAgent != "" {
		eventArchData.VisitorUA = visitorUserAgent
	}
	if ipAddress != "" {
		eventArchData.VisitorIP = ipAddress
	}
}

func createEvent(sdkKey string, eventName string) *request.Event {
	props := createProps(sdkKey)
	return &request.Event{
		Props: props,
		Name:  eventName,
		Time:  time.Now().UnixNano() / 1e6,
	}
}

func createProps(sdkKey string) *request.Props {
	return &request.Props{
		SDKName:    constants.SDKName,
		SDKVersion: constants.SDKVersion,
		EnvKey:     sdkKey,
		Product:    constants.FME,
	}
}

func createVisitor(sdkKey string) *request.Visitor {
	visitorProps := map[string]interface{}{
		constants.VWOFsEnvironment: sdkKey,
	}
	return &request.Visitor{
		Props: visitorProps,
	}
}

func addCustomEventProperties(properties *request.EventArchPayload, eventProperties map[string]interface{}) {
	if eventProperties != nil {
		properties.D.Event.Props.AdditionalProperties = eventProperties
	}
}

func convertToMap(properties *request.EventArchPayload) map[string]interface{} {
	// Convert struct to map using JSON marshaling/unmarshaling
	jsonData, _ := json.Marshal(properties)
	var result map[string]interface{}
	json.Unmarshal(jsonData, &result)
	return result
}
