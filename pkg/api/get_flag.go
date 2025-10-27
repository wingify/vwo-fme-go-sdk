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

package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/decorators"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/storage"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	loggerEnums "github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/services"
	"github.com/wingify/vwo-fme-go-sdk/pkg/utils"
)

// GetFlag retrieves a feature flag value and returns the decision object
func GetFlag(featureKey string, context *user.VWOContext, serviceContainer interfaces.ServiceContainerInterface) models.GetFlagResponse {
	getFlag := models.NewGetFlag()
	shouldCheckForExperimentsRules := false

	passedRulesInformation := make(map[string]interface{})
	evaluatedFeatureMap := make(map[string]interface{})

	// Get feature object from feature key
	feature := utils.GetFeatureFromKey(serviceContainer.GetSettings(), featureKey)

	// Decision object to be sent for the integrations
	decision := map[string]interface{}{
		enums.DecisionFeatureName.GetValue(): nil,
		enums.DecisionFeatureID.GetValue():   nil,
		enums.DecisionFeatureKey.GetValue():  nil,
		enums.DecisionUserID.GetValue():      nil,
		enums.DecisionAPI.GetValue():         enums.ApiGetFlag,
	}

	if feature != nil {
		decision[enums.DecisionFeatureName.GetValue()] = feature.GetName()
		decision[enums.DecisionFeatureID.GetValue()] = feature.GetID()
		decision[enums.DecisionFeatureKey.GetValue()] = feature.GetKey()
	}
	if context != nil {
		decision[enums.DecisionUserID.GetValue()] = context.GetID()
	}

	// create standard debug props
	standardDebugProps := map[string]interface{}{
		enums.DebugPropAPI.GetValue():        enums.ApiGetFlag,
		enums.DebugPropFeatureKey.GetValue(): featureKey,
		enums.DebugPropUUID.GetValue():       context.GetUUID(),
		enums.DebugPropSessionID.GetValue():  context.GetSessionId(),
	}

	// add standard debug props to the debugger service
	serviceContainer.GetDebuggerService().AddStandardDebugProps(standardDebugProps)

	// Check storage for existing data
	storageService := services.NewStorageService()
	storageDecorator := decorators.NewStorageDecorator()
	storedDataMap := storageDecorator.GetFeatureFromStorage(featureKey, context, storageService, serviceContainer)

	// If feature is found in storage, return the stored variation
	if storedDataMap != nil {
		storedData, err := parseStoredData(storedDataMap)
		if err != nil {
			// Log error parsing stored data
			serviceContainer.GetLoggerService().Error("ERROR_READING_DATA_FROM_STORAGE", map[string]interface{}{"err": err.Error()}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		} else if storedData != nil {
			if storedData.FeatureID != 0 && utils.IsFeaturePresentInSettings(serviceContainer.GetSettings(), storedData.FeatureID) {
				// Check for experiment variation
				if storedData.ExperimentVariationID != 0 {
					if storedData.ExperimentKey != "" {
						variation := utils.GetVariationFromCampaignKey(serviceContainer.GetSettings(), storedData.ExperimentKey, storedData.ExperimentVariationID)
						if variation.GetID() != 0 {
							// Log using proper structured logging
							serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["STORED_VARIATION_FOUND"], map[string]interface{}{
								"variationKey":   variation.GetName(),
								"userId":         context.GetID(),
								"experimentType": "experiment",
								"experimentKey":  storedData.ExperimentKey,
							}))
							getFlag.SetIsEnabled(true)
							variables := variation.GetVariables()
							if variables != nil {
								getFlag.SetVariables(convertVariationsToVariables(variables))
							} else {
								getFlag.SetVariables([]*models.Variable{})
							}
							return getFlag
						}
					}
				} else if storedData.RolloutKey != "" && storedData.RolloutID != 0 {
					variation := utils.GetVariationFromCampaignKey(serviceContainer.GetSettings(), storedData.RolloutKey, storedData.RolloutVariationID)
					if variation.GetID() != 0 {
						// Log using proper structured logging
						serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["STORED_VARIATION_FOUND"], map[string]interface{}{
							"variationKey":   variation.GetName(),
							"userId":         context.GetID(),
							"experimentType": "rollout",
							"experimentKey":  storedData.RolloutKey,
						}))

						serviceContainer.GetLoggerService().Debug(log.BuildMessage(log.DebugLogMessagesEnum["EXPERIMENTS_EVALUATION_WHEN_ROLLOUT_PASSED"], map[string]interface{}{
							"userId": context.GetID(),
						}))

						getFlag.SetIsEnabled(true)
						variables := variation.GetVariables()
						if variables != nil {
							getFlag.SetVariables(convertVariationsToVariables(variables))
						} else {
							getFlag.SetVariables([]*models.Variable{})
						}
						shouldCheckForExperimentsRules = true
						featureInfo := map[string]interface{}{
							enums.DecisionRolloutID.GetValue():          storedData.RolloutID,
							enums.DecisionRolloutKey.GetValue():         storedData.RolloutKey,
							enums.DecisionRolloutVariationID.GetValue(): storedData.RolloutVariationID,
						}
						evaluatedFeatureMap[featureKey] = featureInfo
						// Copy featureInfo to passedRulesInformation
						for k, v := range featureInfo {
							passedRulesInformation[k] = v
						}
					}
				}
			}
		}
	}

	// If feature is not found, return false
	if feature == nil {
		serviceContainer.GetLoggerService().Error("FEATURE_NOT_FOUND", map[string]interface{}{
			"featureKey": featureKey,
		}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		getFlag.SetIsEnabled(false)
		return getFlag
	}

	// Set contextual data for segmentation
	serviceContainer.GetSegmentationManager().SetContextualData(serviceContainer, feature, context)

	// Get all rollout rules and evaluate them
	rolloutRules := utils.GetSpecificRulesBasedOnType(feature, enums.CampaignTypeRollout)
	if len(rolloutRules) > 0 && !getFlag.IsEnabled() {
		rolloutRulesToEvaluate := []*campaign.Campaign{}
		for _, rule := range rolloutRules {
			evaluateRuleResult := utils.EvaluateRule(
				serviceContainer,
				feature,
				rule,
				context,
				evaluatedFeatureMap,
				make(map[int]string),
				storageService,
				decision,
			)
			preSegmentationResult := evaluateRuleResult[enums.EvaluatedRuleResultPreSegmentationResult.GetValue()].(bool)
			if preSegmentationResult {
				rolloutRulesToEvaluate = append(rolloutRulesToEvaluate, rule)
				featureMap := map[string]interface{}{
					enums.DecisionRolloutID.GetValue():          rule.GetID(),
					enums.DecisionRolloutKey.GetValue():         strings.Split(rule.GetKey(), featureKey+"_")[1],
					enums.DecisionRolloutVariationID.GetValue(): rule.GetVariations()[0].GetID(),
				}
				evaluatedFeatureMap[featureKey] = featureMap
				break
			}
		}

		// Evaluate the passed rollout rule traffic and get the variation
		if len(rolloutRulesToEvaluate) > 0 {
			passedRolloutCampaign := rolloutRulesToEvaluate[0]
			variation := utils.EvaluateTrafficAndGetVariation(
				serviceContainer,
				passedRolloutCampaign,
				context.GetID(),
			)
			if variation != nil {
				getFlag.SetIsEnabled(true)
				variables := variation.GetVariables()
				if variables != nil {
					getFlag.SetVariables(convertVariationsToVariables(variables))
				} else {
					getFlag.SetVariables([]*models.Variable{})
				}
				shouldCheckForExperimentsRules = true
				updateIntegrationsDecisionObject(passedRolloutCampaign, variation, passedRulesInformation, decision)
				utils.CreateAndSendImpressionForVariationShown(
					serviceContainer,
					passedRolloutCampaign.GetID(),
					variation.GetID(),
					context,
					featureKey,
				)
			}
		}
	} else if !shouldCheckForExperimentsRules {
		serviceContainer.GetLoggerService().Debug(log.BuildMessage(log.DebugLogMessagesEnum["EXPERIMENTS_EVALUATION_WHEN_NO_ROLLOUT_PRESENT"], map[string]interface{}{}))
		shouldCheckForExperimentsRules = true
	}

	// If any rollout rule passed, check for experiment rules
	if shouldCheckForExperimentsRules {
		experimentRulesToEvaluate := []*campaign.Campaign{}
		experimentRules := utils.GetAllExperimentRules(feature)
		megGroupWinnerCampaigns := make(map[int]string)
		for _, rule := range experimentRules {
			evaluateRuleResult := utils.EvaluateRule(
				serviceContainer,
				feature,
				rule,
				context,
				evaluatedFeatureMap,
				megGroupWinnerCampaigns,
				storageService,
				decision,
			)
			preSegmentationResult := evaluateRuleResult[enums.EvaluatedRuleResultPreSegmentationResult.GetValue()].(bool)
			if preSegmentationResult {
				whitelistedObject := evaluateRuleResult[enums.EvaluatedRuleResultWhitelistedObject.GetValue()]

				// Check for nil properly - handle typed nil case (Go gotcha: typed nil != untyped nil)
				var isNil bool
				if whitelistedObject == nil {
					isNil = true
				} else if variation, ok := whitelistedObject.(*campaign.Variation); ok && variation == nil {
					isNil = true
				}

				if isNil {
					experimentRulesToEvaluate = append(experimentRulesToEvaluate, rule)
				} else {
					// If whitelisted object is not null, update the decision object
					if variation, ok := whitelistedObject.(*campaign.Variation); ok {
						getFlag.SetIsEnabled(true)
						variables := variation.GetVariables()
						if variables != nil {
							getFlag.SetVariables(convertVariationsToVariables(variables))
						} else {
							getFlag.SetVariables([]*models.Variable{})
						}
						passedRulesInformation[enums.DecisionExperimentID.GetValue()] = rule.GetID()
						passedRulesInformation[enums.DecisionExperimentKey.GetValue()] = rule.GetKey()
						passedRulesInformation[enums.DecisionExperimentVariationID.GetValue()] = variation.GetID()

						// create and send impression for whitelisted variation
						utils.CreateAndSendImpressionForVariationShown(
							serviceContainer,
							rule.GetID(),
							variation.GetID(),
							context,
							featureKey,
						)
					}
				}
				break
			}
		}

		// Evaluate the passed experiment rule traffic and get the variation
		if len(experimentRulesToEvaluate) > 0 {
			campaign := experimentRulesToEvaluate[0]
			variation := utils.EvaluateTrafficAndGetVariation(
				serviceContainer,
				campaign,
				context.GetID(),
			)
			if variation != nil {
				getFlag.SetIsEnabled(true)
				variables := variation.GetVariables()
				if variables != nil {
					getFlag.SetVariables(convertVariationsToVariables(variables))
				} else {
					getFlag.SetVariables([]*models.Variable{})
				}
				updateIntegrationsDecisionObject(campaign, variation, passedRulesInformation, decision)
				utils.CreateAndSendImpressionForVariationShown(
					serviceContainer,
					campaign.GetID(),
					variation.GetID(),
					context,
					featureKey,
				)
			}
		}
	}

	// Store data if flag is enabled
	if getFlag.IsEnabled() {
		storageMap := map[string]interface{}{
			enums.StorageFeatureKey.GetValue(): feature.GetKey(),
			enums.StorageUserID.GetValue():     context.GetID(),
			enums.StorageFeatureID.GetValue():  feature.GetID(),
		}
		// Copy passedRulesInformation to storageMap
		for k, v := range passedRulesInformation {
			storageMap[k] = v
		}
		storageDecorator.SetDataInStorage(storageMap, storageService, serviceContainer)
	}

	// Execute the integrations
	serviceContainer.GetHooksManager().Set(decision)
	serviceContainer.GetHooksManager().Execute(serviceContainer.GetHooksManager().Get())

	// if debugger is enabled, update the debug event props
	if feature.GetIsDebuggerEnabled() {
		updateDebugEventProps(serviceContainer, decision)
		utils.SendDebugEventToVWO(serviceContainer.GetSettingsManager(), serviceContainer.GetDebuggerService().GetDebugEventProps(enums.DebuggerCategoryDecision.GetValue()))
	}

	// Handle impact campaign
	if feature.GetImpactCampaign() != nil && feature.GetImpactCampaign().GetCampaignID() != 0 {
		// Log impact analysis
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["IMPACT_ANALYSIS"], map[string]interface{}{
			"userId":     context.GetID(),
			"featureKey": featureKey,
			"status": func() string {
				if getFlag.IsEnabled() {
					return "enabled"
				}
				return "disabled"
			}(),
		}))

		// Send impression for impact campaign
		variationID := 1 // disabled
		if getFlag.IsEnabled() {
			variationID = 2 // enabled
		}
		utils.CreateAndSendImpressionForVariationShown(
			serviceContainer,
			feature.GetImpactCampaign().GetCampaignID(),
			variationID,
			context,
			featureKey,
		)
	}

	return getFlag
}

// parseStoredData parses stored data from map
func parseStoredData(data map[string]interface{}) (*storage.StorageData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var storedData storage.StorageData
	err = json.Unmarshal(jsonData, &storedData)
	if err != nil {
		return nil, err
	}

	return &storedData, nil
}

// convertVariationsToVariables converts campaign variations to models variables
func convertVariationsToVariables(variations []campaign.Variable) []*models.Variable {
	result := make([]*models.Variable, 0, len(variations))
	for _, variation := range variations {
		result = append(result, &models.Variable{
			Key:   variation.GetKey(),
			Value: normalizeVariableValue(variation.GetType(), variation.GetValue()),
			Type:  variation.GetType(),
			Id:    variation.GetID(),
		})
	}
	return result
}

// normalizeVariableValue: for "integer" variables, if the value is float64, return int64; otherwise return the original value
func normalizeVariableValue(varType string, val interface{}) interface{} {
	switch strings.ToLower(varType) {
	case "integer":
		if f, ok := val.(float64); ok {
			return int64(f)
		}
		return val
	default:
		return val
	}
}

// updateIntegrationsDecisionObject updates the decision object with campaign and variation details
func updateIntegrationsDecisionObject(campaign *campaign.Campaign, variation *campaign.Variation, passedRulesInformation map[string]interface{}, decision map[string]interface{}) {
	if campaign.GetType() == enums.CampaignTypeRollout.GetValue() {
		passedRulesInformation[enums.DecisionRolloutID.GetValue()] = campaign.GetID()
		passedRulesInformation[enums.DecisionRolloutKey.GetValue()] = campaign.GetKey()
		passedRulesInformation[enums.DecisionRolloutVariationID.GetValue()] = variation.GetID()
	} else {
		passedRulesInformation[enums.DecisionExperimentID.GetValue()] = campaign.GetID()
		passedRulesInformation[enums.DecisionExperimentKey.GetValue()] = campaign.GetKey()
		passedRulesInformation[enums.DecisionExperimentVariationID.GetValue()] = variation.GetID()
	}

	// Copy passedRulesInformation to decision
	for k, v := range passedRulesInformation {
		decision[k] = v
	}
}

// updateDebugEventProps updates the debug event props with the decision keys
// @param serviceContainer ServiceContainer object containing the debugger service
// @param decision Map containing the decision object
func updateDebugEventProps(serviceContainer interfaces.ServiceContainerInterface, decision map[string]interface{}) {
	decisionKeys := make(map[string]interface{})
	featureKey := decision[enums.DecisionFeatureKey.GetValue()].(string)
	message := fmt.Sprintf("Flag decision given for feature:%s.", featureKey)

	// Check for rollout information
	if rolloutKey, exists := decision[enums.DecisionRolloutKey.GetValue()]; exists && rolloutKey != nil && rolloutKey != "" {
		if rolloutVariationId, exists := decision[enums.DecisionRolloutVariationID.GetValue()]; exists && rolloutVariationId != nil && rolloutVariationId != "" {
			rolloutKeyStr := rolloutKey.(string)
			// Split rollout key to extract just the rollout part (remove featureKey_ prefix)
			if strings.HasPrefix(rolloutKeyStr, featureKey+"_") {
				rolloutKeyStr = rolloutKeyStr[len(featureKey)+1:]
			}
			message += fmt.Sprintf(" Got Rollout:%s. Rollout variation id:%v.", rolloutKeyStr, rolloutVariationId)
		}
	}

	// Check for experiment information
	if experimentKey, exists := decision[enums.DecisionExperimentKey.GetValue()]; exists && experimentKey != nil && experimentKey != "" {
		if experimentVariationId, exists := decision[enums.DecisionExperimentVariationID.GetValue()]; exists && experimentVariationId != nil && experimentVariationId != "" {
			experimentKeyStr := experimentKey.(string)
			// Split experiment key to extract just the experiment part (remove featureKey_ prefix)
			if strings.HasPrefix(experimentKeyStr, featureKey+"_") {
				experimentKeyStr = experimentKeyStr[len(featureKey)+1:]
			}
			message += fmt.Sprintf(" Got Experiment:%s. Experiment variation id:%v.", experimentKeyStr, experimentVariationId)
		}
	}

	decisionKeys[enums.DebugPropMessage.GetValue()] = message
	decisionKeys[enums.DebugPropMessageType.GetValue()] = constants.FLAG_DECISION
	decisionKeys[enums.DebugPropLogLevel.GetValue()] = loggerEnums.LogLevelInfo.String()
	serviceContainer.GetDebuggerService().AddCategoryDebugProps(enums.DebuggerCategoryDecision.GetValue(), decisionKeys)
}
