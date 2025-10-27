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
	"strconv"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	_ "github.com/wingify/vwo-fme-go-sdk/pkg/models/settings" // Used indirectly by GetGroupDetailsIfCampaignPartOfIt
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/storage"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/decision_maker"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// CheckWhitelistingAndPreSeg checks whitelisting and pre-segmentation for a campaign
func CheckWhitelistingAndPreSeg(
	serviceContainer interfaces.ServiceContainerInterface, // ServiceContainer object
	feature *campaign.Feature,
	campaignModel *campaign.Campaign,
	context *user.VWOContext,
	evaluatedFeatureMap map[string]interface{},
	megGroupWinnerCampaigns map[int]string,
	storageService interfaces.StorageServiceInterface,
	decision map[string]interface{},
) map[string]interface{} {
	accountID := serviceContainer.GetSettingsManager().GetAccountID()

	// Check for nil pointers first
	if context == nil {
		return map[string]interface{}{
			enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
			enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
		}
	}

	vwoUserID := GetUUID(context.ID, accountID)
	if campaignModel == nil {
		return map[string]interface{}{
			enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
			enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
		}
	}

	if serviceContainer == nil {
		return map[string]interface{}{
			enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
			enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
		}
	}

	if campaignModel.Type == string(enums.CampaignTypeAB) {
		variationTargetingVariables := make(map[string]interface{})
		if context.VariationTargetingVariables != nil {
			for k, v := range context.VariationTargetingVariables {
				variationTargetingVariables[k] = v
			}
		}

		if campaignModel.IsUserListEnabled {
			variationTargetingVariables[constants.VariationTargetingUserIDKey] = vwoUserID
		} else {
			variationTargetingVariables[constants.VariationTargetingUserIDKey] = context.ID
		}
		context.VariationTargetingVariables = variationTargetingVariables
		decision[enums.ContextVariationTargetingVariables.GetValue()] = context.VariationTargetingVariables // for integration

		if campaignModel.IsForcedVariationEnabled {
			whitelistingResult := checkCampaignWhitelisting(campaignModel, context, serviceContainer)
			if whitelistingResult != nil {
				return map[string]interface{}{
					enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): true,
					enums.EvaluatedRuleResultWhitelistedObject.GetValue():     whitelistingResult["variation"],
				}
			}
		} else {
			serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["WHITELISTING_SKIP"], map[string]interface{}{
				"userId":      context.ID,
				"campaignKey": campaignModel.RuleKey,
			}))
		}
	}

	// Set _vwoUserId for custom variables
	customVariables := make(map[string]interface{})
	if context.CustomVariables != nil {
		for k, v := range context.CustomVariables {
			customVariables[k] = v
		}
	}

	if campaignModel.IsUserListEnabled {
		customVariables[constants.VariationTargetingUserIDKey] = vwoUserID
	} else {
		customVariables[constants.VariationTargetingUserIDKey] = context.ID
	}
	context.CustomVariables = customVariables
	decision[enums.ContextCustomVariables.GetValue()] = context.CustomVariables // for integration

	// Check if Rule being evaluated is part of Mutually Exclusive Group
	settingsModel := serviceContainer.GetSettings()
	groupDetails := GetGroupDetailsIfCampaignPartOfIt(settingsModel, campaignModel.ID, getVariationIDForPersonalize(campaignModel))
	groupId := ""
	if groupDetails != nil {
		if id, exists := groupDetails["groupId"]; exists && id != "" {
			groupId = id
		}
	}

	// MEG handling
	if groupId != "" && groupId != "null" {
		// Check if the group is already evaluated for the user
		groupIdInt, err := strconv.Atoi(groupId)
		if err == nil {
			groupWinnerCampaignId, exists := megGroupWinnerCampaigns[groupIdInt]
			if exists && groupWinnerCampaignId != "" {
				if campaignModel.Type == string(enums.CampaignTypeAB) {
					if groupWinnerCampaignId == strconv.Itoa(campaignModel.ID) {
						// If the campaign is the winner of the MEG, return true
						return map[string]interface{}{
							enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): true,
							enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
						}
					}
				} else if campaignModel.Type == string(enums.CampaignTypePersonalize) {
					// if personalise then check if the requested variation is the winner
					if groupWinnerCampaignId == strconv.Itoa(campaignModel.ID)+"_"+strconv.Itoa(campaignModel.Variations[0].ID) {
						// If the campaign is the winner of the MEG, return true
						return map[string]interface{}{
							enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): true,
							enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
						}
					}
				}
				// If the campaign is not the winner of the MEG, return false
				return map[string]interface{}{
					enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
					enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
				}
			} else {
				// Check in storage if the group is already evaluated for the user
				storedDataMap, err := storageService.GetDataInStorage(constants.VWOMetaMegKey+strconv.Itoa(groupIdInt), context)
				if err != nil {
					serviceContainer.GetLoggerService().Error("ERROR_READING_DATA_FROM_STORAGE", map[string]interface{}{"err": err.Error()}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
				} else {
					if len(storedDataMap) > 0 {
						// Convert map to JSON and then to Storage struct
						jsonData, err := json.Marshal(storedDataMap)
						if err == nil {
							var storedData storage.StorageData
							if json.Unmarshal(jsonData, &storedData) == nil {
								if storedData.GetExperimentID() != 0 && storedData.GetExperimentKey() != "" {
									serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_CAMPAIGN_FOUND_IN_STORAGE"], map[string]interface{}{
										"campaignKey": storedData.GetExperimentKey(),
										"userId":      context.ID,
									}))

									if storedData.GetExperimentID() == campaignModel.ID {
										if campaignModel.Type == string(enums.CampaignTypePersonalize) {
											// if personalise then check if the requested variation is the winner
											if storedData.GetExperimentVariationID() == campaignModel.Variations[0].ID {
												return map[string]interface{}{
													enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): true,
													enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
												}
											} else {
												// store the campaign in local cache, so that it can be used later without looking into user storage again
												megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(storedData.GetExperimentID()) + "_" + strconv.Itoa(storedData.GetExperimentVariationID())
												return map[string]interface{}{
													enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
													enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
												}
											}
										} else {
											// return the campaign if the called campaignId matches
											return map[string]interface{}{
												enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): true,
												enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
											}
										}
									}

									// if experimentId is not -1 then campaign is personalise campaign, store the details and return
									if storedData.GetExperimentVariationID() != -1 {
										megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(storedData.GetExperimentID()) + "_" + strconv.Itoa(storedData.GetExperimentVariationID())
									} else {
										// else store the campaignId only and return
										megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(storedData.GetExperimentID())
									}
									return map[string]interface{}{
										enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
										enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// If Whitelisting is skipped/failed, Check campaign's pre-segmentation
	isPreSegmentationPassed := GetPreSegmentationDecision(campaignModel, context, serviceContainer)

	// MEG evaluation - full implementation using MegUtil
	if isPreSegmentationPassed && groupId != "" && groupId != "null" {
		groupIdInt, err := strconv.Atoi(groupId)
		if err == nil {
			// Find the feature for this campaign
			var feature *campaign.Feature
			for _, f := range serviceContainer.GetSettings().GetFeatures() {
				for _, rule := range f.GetRulesLinkedCampaign() {
					if rule.GetID() == campaignModel.ID {
						feature = &f
						break
					}
				}
				if feature != nil {
					break
				}
			}

			if feature != nil {
				megUtil := NewMegUtil()
				variationModel := megUtil.EvaluateGroups(
					serviceContainer,
					feature,
					groupIdInt,
					evaluatedFeatureMap,
					context,
					storageService,
				)

				// Check if the current campaign is the winner
				// this condition would be true only when the current campaignId match with group winner campaignId
				// for personalise campaign, all personalise variations have same campaignId, so we check for campaignId_variationId
				if variationModel != nil && variationModel.GetID() != 0 && variationModel.GetID() == campaignModel.ID {
					// if campaign is AB then return true
					if variationModel.GetType() == string(enums.CampaignTypeAB) {
						return map[string]interface{}{
							enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): true,
							enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
						}
					} else {
						// if personalise then check if the requested variation is the winner
						if len(variationModel.GetVariations()) > 0 && variationModel.GetVariations()[0].GetID() == campaignModel.Variations[0].ID {
							return map[string]interface{}{
								enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): true,
								enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
							}
						} else {
							// store the campaign in local cache, so that it can be used later
							if len(variationModel.GetVariations()) > 0 {
								megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(variationModel.GetID()) + "_" + strconv.Itoa(variationModel.GetVariations()[0].GetID())
							} else {
								megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(variationModel.GetID())
							}
							return map[string]interface{}{
								enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
								enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
							}
						}
					}
				} else if variationModel != nil && variationModel.GetID() != 0 {
					// when there is a winner but not the current campaign
					if variationModel.GetType() == string(enums.CampaignTypeAB) {
						// if campaign is AB then store only the campaignId
						megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(variationModel.GetID())
					} else {
						// if campaign is personalise then store the campaignId_variationId
						if len(variationModel.GetVariations()) > 0 {
							megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(variationModel.GetID()) + "_" + strconv.Itoa(variationModel.GetVariations()[0].GetID())
						} else {
							megGroupWinnerCampaigns[groupIdInt] = strconv.Itoa(variationModel.GetID())
						}
					}
					return map[string]interface{}{
						enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
						enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
					}
				}
				// store -1 if no winner found, so that we don't evaluate the group again as the result would be the same for the current getFlag call
				megGroupWinnerCampaigns[groupIdInt] = "-1"
				return map[string]interface{}{
					enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
					enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
				}
			}
		}
	}

	return map[string]interface{}{
		enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): isPreSegmentationPassed,
		enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
	}
}

// EvaluateTrafficAndGetVariation evaluates the traffic for a given campaign and get the variation
func EvaluateTrafficAndGetVariation(
	serviceContainer interfaces.ServiceContainerInterface, // ServiceContainer object
	campaignModel *campaign.Campaign,
	userID string,
) *campaign.Variation {
	// Get the variation allotted to the user
	var variation *campaign.Variation = nil
	accountID := serviceContainer.GetSettingsManager().GetAccountID()

	variation = GetVariationAllotted(userID, accountID, campaignModel, serviceContainer)

	if variation == nil {
		// Log that user did not get any variation
		if serviceContainer != nil {
			campaignKey := getCampaignKey(campaignModel)
			serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["USER_CAMPAIGN_BUCKET_INFO"], map[string]interface{}{
				"userId":      userID,
				"campaignKey": campaignKey,
				"status":      "did not get any variation",
			}))
		}
		return nil
	}

	// Log that user got a variation
	if serviceContainer != nil {
		campaignKey := getCampaignKey(campaignModel)
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["USER_CAMPAIGN_BUCKET_INFO"], map[string]interface{}{
			"userId":      userID,
			"campaignKey": campaignKey,
			"status":      "got variation: " + variation.Name,
		}))
	}

	return variation
}

// checkCampaignWhitelisting checks for whitelisting
func checkCampaignWhitelisting(
	campaignModel *campaign.Campaign,
	context *user.VWOContext,
	serviceContainer interfaces.ServiceContainerInterface,
) map[string]interface{} {
	whitelistingResult := evaluateWhitelisting(campaignModel, context, serviceContainer)

	// Log whitelisting status
	status := "failed"
	variationString := ""
	if whitelistingResult != nil {
		status = "passed"
		variationString = whitelistingResult["variationName"].(string)
	}

	if serviceContainer != nil {
		campaignKey := getCampaignKey(campaignModel)
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["WHITELISTING_STATUS"], map[string]interface{}{
			"userId":          context.ID,
			"campaignKey":     campaignKey,
			"status":          status,
			"variationString": variationString,
		}))
	}

	return whitelistingResult
}

// evaluateWhitelisting evaluates whitelisting for a campaign
func evaluateWhitelisting(
	campaignModel *campaign.Campaign,
	context *user.VWOContext,
	serviceContainer interfaces.ServiceContainerInterface,
) map[string]interface{} {
	targetedVariations := []campaign.Variation{}

	for _, variation := range campaignModel.Variations {
		if variation.Segments != nil && len(variation.Segments) == 0 {
			// Log WHITELISTING_SKIP
			if serviceContainer != nil {
				variationInfo := ""
				if variation.Name != "" {
					variationInfo = "for variation: " + variation.Name
				}
				campaignKey := getCampaignKey(campaignModel)
				serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["WHITELISTING_SKIP"], map[string]interface{}{
					"userId":      context.ID,
					"campaignKey": campaignKey,
					"variation":   variationInfo,
				}))
			}
			continue
		}

		// Check for segmentation and evaluate
		if variation.Segments != nil {
			segmentationResult := false // Default fallback
			if serviceContainer != nil {
				segmentationResult = serviceContainer.GetSegmentationManager().ValidateSegmentation(variation.Segments, context.VariationTargetingVariables)
			}

			if segmentationResult {
				clonedVariation := variation // Simple copy for now - can be enhanced with deep cloning later
				targetedVariations = append(targetedVariations, clonedVariation)
			}
		}
	}

	var whitelistedVariation *campaign.Variation

	if len(targetedVariations) > 1 {
		ScaleVariationWeights(targetedVariations)
		currentAllocation := 0
		for i := range targetedVariations {
			stepFactor := AssignRangeValues(&targetedVariations[i], currentAllocation)
			currentAllocation += stepFactor
		}
		bucketingSeed := GetBucketingSeed(context.ID, campaignModel, nil)
		bucketValue := decision_maker.CalculateBucketValue(bucketingSeed)

		whitelistedVariation = GetVariation(targetedVariations, bucketValue)
	} else if len(targetedVariations) == 1 {
		whitelistedVariation = &targetedVariations[0]
	}

	if whitelistedVariation != nil {
		// Return map with variation details
		return map[string]interface{}{
			"variation":     whitelistedVariation,
			"variationName": whitelistedVariation.Name,
			"variationId":   whitelistedVariation.ID,
		}
	}

	return nil
}

// getCampaignKey returns the campaign key based on campaign type
func getCampaignKey(campaignModel *campaign.Campaign) string {
	if campaignModel.Type == string(enums.CampaignTypeAB) {
		return campaignModel.Key
	}
	return campaignModel.Name + "_" + campaignModel.RuleKey
}

// getVariationIDForPersonalize returns variation ID for personalize campaigns
func getVariationIDForPersonalize(campaignModel *campaign.Campaign) int {
	if campaignModel.Type == string(enums.CampaignTypePersonalize) {
		return campaignModel.Variations[0].ID
	}
	return -1
}

// IsUserPartOfCampaign checks if the user is part of the campaign based on traffic allocation
func IsUserPartOfCampaign(
	userID string,
	campaignModel *campaign.Campaign,
	serviceContainer interfaces.ServiceContainerInterface,
) bool {
	if campaignModel == nil || userID == "" {
		return false
	}

	var trafficAllocation float64
	var salt string

	// Check if the campaign is of type ROLLOUT or PERSONALIZE
	campaignType := campaignModel.Type
	isRolloutOrPersonalize := campaignType == string(enums.CampaignTypeRollout) ||
		campaignType == string(enums.CampaignTypePersonalize)

		// Get salt and traffic allocation based on campaign type
	if isRolloutOrPersonalize {
		salt = campaignModel.Variations[0].Salt
		trafficAllocation = campaignModel.Variations[0].Weight
	} else {
		salt = campaignModel.Salt
		trafficAllocation = float64(campaignModel.PercentTraffic)
	}

	// Generate bucket key using salt if available, otherwise use campaign ID
	var bucketKey string
	if salt != "" {
		bucketKey = salt + "_" + userID
	} else {
		bucketKey = strconv.Itoa(campaignModel.ID) + "_" + userID
	}

	valueAssignedToUser := decision_maker.GetBucketValueForUser(bucketKey)
	isUserPart := valueAssignedToUser != 0 && float64(valueAssignedToUser) <= trafficAllocation

	campaignKey := getCampaignKey(campaignModel)
	notPart := ""
	if !isUserPart {
		notPart = "not"
	}

	// Log using serviceContainer logger
	if serviceContainer != nil {
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["USER_PART_OF_CAMPAIGN"], map[string]interface{}{
			"userId":      userID,
			"notPart":     notPart,
			"campaignKey": campaignKey,
		}))
	}

	return isUserPart
}

// GetVariation returns the variation based on the bucket value
func GetVariation(
	variations []campaign.Variation,
	bucketValue int,
) *campaign.Variation {
	for i := range variations {
		if bucketValue >= variations[i].StartRangeVariation &&
			bucketValue <= variations[i].EndRangeVariation {
			return &variations[i]
		}
	}
	return nil
}

// CheckInRange checks if the bucket value falls in the range of the variation
func CheckInRange(
	variation *campaign.Variation,
	bucketValue int,
) *campaign.Variation {
	if bucketValue >= variation.StartRangeVariation &&
		bucketValue <= variation.EndRangeVariation {
		return variation
	}
	return nil
}

// BucketUserToVariation buckets the user to a variation based on the bucket value
func BucketUserToVariation(
	userID string,
	accountID string,
	campaignModel *campaign.Campaign,
	serviceContainer interfaces.ServiceContainerInterface,
) *campaign.Variation {
	if campaignModel == nil || userID == "" {
		return nil
	}

	multiplier := 1
	if campaignModel.PercentTraffic == 0 {
		multiplier = 0
	}

	percentTraffic := campaignModel.PercentTraffic
	salt := campaignModel.Salt

	var bucketKey string
	// If salt is not null and not empty, use salt else use campaign id
	if salt != "" {
		bucketKey = salt + "_" + accountID + "_" + userID
	} else {
		bucketKey = strconv.Itoa(campaignModel.ID) + "_" + accountID + "_" + userID
	}

	hashValue := decision_maker.GenerateHashValue(bucketKey)
	bucketValue := decision_maker.GenerateBucketValue(hashValue, constants.MaxTrafficValue, multiplier)

	// Log using serviceContainer logger
	if serviceContainer != nil {
		serviceContainer.GetLoggerService().Debug(log.BuildMessage(log.DebugLogMessagesEnum["USER_BUCKET_TO_VARIATION"], map[string]interface{}{
			"userId":         userID,
			"campaignKey":    campaignModel.RuleKey,
			"percentTraffic": strconv.Itoa(percentTraffic),
			"bucketValue":    strconv.Itoa(bucketValue),
			"hashValue":      strconv.FormatUint(uint64(hashValue), 10),
		}))
	}

	return GetVariation(campaignModel.Variations, bucketValue)
}

// GetPreSegmentationDecision analyzes the pre-segmentation decision for the user in the campaign
func GetPreSegmentationDecision(
	campaignModel *campaign.Campaign,
	context *user.VWOContext,
	serviceContainer interfaces.ServiceContainerInterface,
) bool {
	campaignType := campaignModel.Type
	var segments map[string]interface{}

	if campaignType == string(enums.CampaignTypeRollout) ||
		campaignType == string(enums.CampaignTypePersonalize) {
		segments = campaignModel.Variations[0].Segments
	} else if campaignType == string(enums.CampaignTypeAB) {
		segments = campaignModel.Segments
	} else {
		segments = make(map[string]interface{})
	}

	campaignKey := getCampaignKey(campaignModel)

	if len(segments) == 0 {
		// Log using serviceContainer logger
		if serviceContainer != nil {
			serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["SEGMENTATION_SKIP"], map[string]interface{}{
				"userId":      context.ID,
				"campaignKey": campaignKey,
			}))
		}
		return true
	}

	// Convert CustomVariables to map[string]interface{}
	customVariables := make(map[string]interface{})
	if context.CustomVariables != nil {
		for k, v := range context.CustomVariables {
			customVariables[k] = v
		}
	}

	// Use serviceContainer segmentation manager
	var preSegmentationResult bool
	if serviceContainer != nil {
		preSegmentationResult = serviceContainer.GetSegmentationManager().ValidateSegmentation(segments, customVariables)
	}

	status := "failed"
	if preSegmentationResult {
		status = "passed"
	}

	// Log using serviceContainer logger
	if serviceContainer != nil {
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["SEGMENTATION_STATUS"], map[string]interface{}{
			"userId":      context.ID,
			"campaignKey": campaignKey,
			"status":      status,
		}))
	}

	return preSegmentationResult
}

// GetVariationAllotted gets the variation allotted to the user in the campaign
func GetVariationAllotted(
	userID string,
	accountID string,
	campaignModel *campaign.Campaign,
	serviceContainer interfaces.ServiceContainerInterface,
) *campaign.Variation {
	isUserPart := IsUserPartOfCampaign(userID, campaignModel, serviceContainer)
	if campaignModel.Type == string(enums.CampaignTypeRollout) ||
		campaignModel.Type == string(enums.CampaignTypePersonalize) {
		if isUserPart {
			return &campaignModel.Variations[0]
		}
		return nil
	}

	if isUserPart {
		return BucketUserToVariation(userID, accountID, campaignModel, serviceContainer)
	}
	return nil
}
