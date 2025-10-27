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
	"math"
	"strconv"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/decorators"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	campaign "github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/storage"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/decision_maker"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// MegUtil handles Mutually Exclusive Group logic
type MegUtil struct{}

// NewMegUtil creates a new MegUtil instance
func NewMegUtil() *MegUtil {
	return &MegUtil{}
}

// EvaluateGroupsResult represents the result of MEG evaluation
type EvaluateGroupsResult struct {
	FeatureKeys      []string
	GroupCampaignIds []string
}

// EligibleCampaignsResult represents eligible campaigns result
type EligibleCampaignsResult struct {
	EligibleCampaigns            []*campaign.Campaign
	EligibleCampaignsWithStorage []*campaign.Campaign
	InEligibleCampaigns          []*campaign.Campaign
}

// EvaluateGroups evaluates groups for a given feature and group ID and returns the winner variation
func (mu *MegUtil) EvaluateGroups(
	serviceContainer interfaces.ServiceContainerInterface,
	feature *campaign.Feature,
	groupId int,
	evaluatedFeatureMap map[string]interface{},
	context *user.VWOContext,
	storageService interfaces.StorageServiceInterface,
) *campaign.Variation {
	featureToSkip := []string{}
	campaignMap := make(map[string][]*campaign.Campaign)

	// Get all feature keys and all campaignIds from the groupId
	settings := serviceContainer.GetSettings()
	featureKeysAndGroupCampaignIds := mu.GetFeatureKeysFromGroup(settings, groupId)
	featureKeys := featureKeysAndGroupCampaignIds.FeatureKeys
	groupCampaignIds := featureKeysAndGroupCampaignIds.GroupCampaignIds

	for _, featureKey := range featureKeys {
		currentFeature := GetFeatureFromKey(settings, featureKey)

		// Check if the feature is already evaluated
		if contains(featureToSkip, featureKey) {
			continue
		}

		// Evaluate the feature rollout rules
		isRolloutRulePassed := mu.isRolloutRuleForFeaturePassed(
			serviceContainer,
			currentFeature,
			evaluatedFeatureMap,
			&featureToSkip,
			context,
			storageService,
		)
		if isRolloutRulePassed {
			// Build campaign map for this feature
			for _, feature1 := range settings.GetFeatures() {
				if feature1.GetKey() == featureKey {
					for _, ruleCampaign := range feature1.GetRulesLinkedCampaign() {
						campaignIDStr := strconv.Itoa(ruleCampaign.GetID())
						campaignWithVariationID := campaignIDStr + "_" + strconv.Itoa(ruleCampaign.GetVariations()[0].GetID())

						if contains(groupCampaignIds, campaignIDStr) || contains(groupCampaignIds, campaignWithVariationID) {
							if campaignMap[featureKey] == nil {
								campaignMap[featureKey] = []*campaign.Campaign{}
							}

							// Check if campaign with same rule key already exists
							exists := false
							for _, existingCampaign := range campaignMap[featureKey] {
								if existingCampaign.GetRuleKey() == ruleCampaign.GetRuleKey() {
									exists = true
									break
								}
							}

							if !exists {
								campaignCopy := ruleCampaign
								campaignMap[featureKey] = append(campaignMap[featureKey], &campaignCopy)
							}
						}
					}
				}
			}
		}
	}

	eligibleCampaignsMap := mu.getEligibleCampaigns(serviceContainer, campaignMap, context, storageService)
	eligibleCampaigns := eligibleCampaignsMap.EligibleCampaigns
	eligibleCampaignsWithStorage := eligibleCampaignsMap.EligibleCampaignsWithStorage

	return mu.findWinnerCampaignAmongEligibleCampaigns(
		serviceContainer,
		feature.GetKey(),
		eligibleCampaigns,
		eligibleCampaignsWithStorage,
		groupId,
		context,
		storageService,
	)
}

// GetFeatureKeysFromGroup retrieves feature keys associated with a group based on the group ID and returns the feature keys and group campaign IDs
func (mu *MegUtil) GetFeatureKeysFromGroup(settings *settings.Settings, groupId int) *EvaluateGroupsResult {
	groupCampaignIds := GetCampaignsByGroupID(settings, groupId)
	featureKeys := GetFeatureKeysFromCampaignIDs(settings, groupCampaignIds)

	return &EvaluateGroupsResult{
		FeatureKeys:      featureKeys,
		GroupCampaignIds: groupCampaignIds,
	}
}

// isRolloutRuleForFeaturePassed evaluates the feature rollout rules for a given feature and returns true if the rollout rule is passed
func (mu *MegUtil) isRolloutRuleForFeaturePassed(
	serviceContainer interfaces.ServiceContainerInterface,
	feature *campaign.Feature,
	evaluatedFeatureMap map[string]interface{},
	featureToSkip *[]string,
	context *user.VWOContext,
	storageService interfaces.StorageServiceInterface,
) bool {
	if evaluatedFeatureMap[feature.GetKey()] != nil {
		if featureMap, ok := evaluatedFeatureMap[feature.GetKey()].(map[string]interface{}); ok {
			if _, exists := featureMap["rolloutId"]; exists {
				return true
			}
		}
	}

	rollOutRules := GetSpecificRulesBasedOnType(feature, enums.CampaignTypeRollout)
	if len(rollOutRules) > 0 {
		var ruleToTestForTraffic *campaign.Campaign

		for _, rule := range rollOutRules {
			preSegmentationResult := EvaluateRule(
				serviceContainer,
				feature,
				rule,
				context,
				evaluatedFeatureMap,
				nil,
				storageService,
				make(map[string]interface{}),
			)
			if preSegmentationResult[enums.EvaluatedRuleResultPreSegmentationResult.GetValue()].(bool) {
				ruleToTestForTraffic = rule
				break
			}
		}

		if ruleToTestForTraffic != nil {
			variation := EvaluateTrafficAndGetVariation(serviceContainer, ruleToTestForTraffic, context.GetID())
			if variation != nil {
				rollOutInformation := map[string]interface{}{
					enums.DecisionRolloutID.GetValue():          variation.GetID(),
					enums.DecisionRolloutKey.GetValue():         variation.GetName(),
					enums.DecisionRolloutVariationID.GetValue(): variation.GetID(),
				}
				evaluatedFeatureMap[feature.GetKey()] = rollOutInformation
				return true
			}
		}

		// No rollout rule passed
		*featureToSkip = append(*featureToSkip, feature.GetKey())
		return false
	}

	// No rollout rule, evaluate experiments
	serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_SKIP_ROLLOUT_EVALUATE_EXPERIMENTS"], map[string]interface{}{
		"featureKey": feature.GetKey(),
	}))
	return true
}

// getEligibleCampaigns retrieves eligible campaigns based on the provided campaign map and context and returns the eligible campaigns and ineligible campaigns
func (mu *MegUtil) getEligibleCampaigns(
	serviceContainer interfaces.ServiceContainerInterface,
	campaignMap map[string][]*campaign.Campaign,
	context *user.VWOContext,
	storageService interfaces.StorageServiceInterface,
) *EligibleCampaignsResult {
	eligibleCampaigns := []*campaign.Campaign{}
	eligibleCampaignsWithStorage := []*campaign.Campaign{}
	inEligibleCampaigns := []*campaign.Campaign{}

	// Note: Storage decorator would be used here, but avoiding import cycle

	for _, campaigns := range campaignMap {
		for _, campaign := range campaigns {
			// Check storage for existing data
			// Note: Storage decorator call would go here, but avoiding import cycle
			// var storedDataMap map[string]interface{}
			storedDataMap := map[string]interface{}{} // Placeholder to avoid nil check

			if len(storedDataMap) > 0 {
				// Convert map to JSON and then to Storage struct
				jsonData, err := json.Marshal(storedDataMap)
				if err == nil {
					var storedData storage.StorageData
					if json.Unmarshal(jsonData, &storedData) == nil {
						if storedData.GetExperimentVariationID() != 0 {
							if storedData.GetExperimentKey() != "" && storedData.GetExperimentKey() == campaign.GetKey() {
								variation := GetVariationFromCampaignKey(serviceContainer.GetSettings(), storedData.GetExperimentKey(), storedData.GetExperimentVariationID())
								if variation.GetID() != 0 {
									serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_CAMPAIGN_FOUND_IN_STORAGE"], map[string]interface{}{
										"campaignKey": storedData.GetExperimentKey(),
										"userId":      context.GetID(),
									}))

									// Check if campaign already exists in eligibleCampaignsWithStorage
									exists := false
									for _, existingCampaign := range eligibleCampaignsWithStorage {
										if existingCampaign.GetKey() == campaign.GetKey() {
											exists = true
											break
										}
									}
									if !exists {
										eligibleCampaignsWithStorage = append(eligibleCampaignsWithStorage, campaign)
									}
									continue
								}
							}
						}
					}
				}
			}

			// Check if user is eligible for the campaign
			if GetPreSegmentationDecision(campaign, context, serviceContainer) &&
				IsUserPartOfCampaign(context.GetID(), campaign, serviceContainer) {
				serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_CAMPAIGN_ELIGIBLE"], map[string]interface{}{
					"campaignKey": getCampaignKey(campaign),
					"userId":      context.GetID(),
				}))
				eligibleCampaigns = append(eligibleCampaigns, campaign)
				continue
			}

			inEligibleCampaigns = append(inEligibleCampaigns, campaign)
		}
	}

	return &EligibleCampaignsResult{
		EligibleCampaigns:            eligibleCampaigns,
		EligibleCampaignsWithStorage: eligibleCampaignsWithStorage,
		InEligibleCampaigns:          inEligibleCampaigns,
	}
}

// findWinnerCampaignAmongEligibleCampaigns evaluates the eligible campaigns and determines the winner campaign and returns the winner variation
func (mu *MegUtil) findWinnerCampaignAmongEligibleCampaigns(
	serviceContainer interfaces.ServiceContainerInterface,
	featureKey string,
	eligibleCampaigns []*campaign.Campaign,
	eligibleCampaignsWithStorage []*campaign.Campaign,
	groupId int,
	context *user.VWOContext,
	storageService interfaces.StorageServiceInterface,
) *campaign.Variation {
	campaignIds := GetCampaignIDsFromFeatureKey(serviceContainer.GetSettings(), featureKey)
	var winnerCampaign *campaign.Variation

	defer func() {
		if r := recover(); r != nil {
			serviceContainer.GetLoggerService().Error("ERROR_FINDING_WINNER_CAMPAIGN", map[string]interface{}{
				"method": "findWinnerCampaignAmongEligibleCampaigns",
				"err":    fmt.Sprint(r),
			}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	groups := serviceContainer.GetSettings().GetGroups()
	group, exists := groups[strconv.Itoa(groupId)]
	megAlgoNumber := constants.RandomAlgo
	if exists && group.GetEt() != 0 {
		megAlgoNumber = group.GetEt()
	}

	if len(eligibleCampaignsWithStorage) == 1 {
		// Convert campaign to variation
		jsonData, err := json.Marshal(eligibleCampaignsWithStorage[0])
		if err == nil {
			json.Unmarshal(jsonData, &winnerCampaign)
		}

		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_WINNER_CAMPAIGN"], map[string]interface{}{
			"campaignKey": getCampaignKey(eligibleCampaignsWithStorage[0]),
			"groupId":     strconv.Itoa(groupId),
			"userId":      context.GetID(),
		}))
	} else if len(eligibleCampaignsWithStorage) > 1 && megAlgoNumber == constants.RandomAlgo {
		winnerCampaign = mu.normalizeWeightsAndFindWinningCampaign(serviceContainer, eligibleCampaignsWithStorage, context, campaignIds, groupId, storageService)
	} else if len(eligibleCampaignsWithStorage) > 1 {
		winnerCampaign = mu.getCampaignUsingAdvancedAlgo(serviceContainer, eligibleCampaignsWithStorage, context, campaignIds, groupId, storageService)
	}

	if len(eligibleCampaignsWithStorage) == 0 {
		if len(eligibleCampaigns) == 1 {
			// Convert campaign to variation
			jsonData, err := json.Marshal(eligibleCampaigns[0])
			if err == nil {
				json.Unmarshal(jsonData, &winnerCampaign)
			}

			serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_WINNER_CAMPAIGN"], map[string]interface{}{
				"campaignKey": getCampaignKey(eligibleCampaigns[0]),
				"groupId":     strconv.Itoa(groupId),
				"userId":      context.GetID(),
				"algo":        "",
			}))
		} else if len(eligibleCampaigns) > 1 && megAlgoNumber == constants.RandomAlgo {
			winnerCampaign = mu.normalizeWeightsAndFindWinningCampaign(serviceContainer, eligibleCampaigns, context, campaignIds, groupId, storageService)
		} else if len(eligibleCampaigns) > 1 {
			winnerCampaign = mu.getCampaignUsingAdvancedAlgo(serviceContainer, eligibleCampaigns, context, campaignIds, groupId, storageService)
		}
	}

	return winnerCampaign
}

// normalizeWeightsAndFindWinningCampaign normalizes the weights of shortlisted campaigns and determines the winning campaign using random allocation and returns the winner variation
func (mu *MegUtil) normalizeWeightsAndFindWinningCampaign(
	serviceContainer interfaces.ServiceContainerInterface,
	shortlistedCampaigns []*campaign.Campaign,
	context *user.VWOContext,
	calledCampaignIds []int,
	groupId int,
	storageService interfaces.StorageServiceInterface,
) *campaign.Variation {
	defer func() {
		if r := recover(); r != nil {
			serviceContainer.GetLoggerService().Error("ERROR_FINDING_WINNER_CAMPAIGN", map[string]interface{}{
				"method": "normalizeWeightsAndFindWinningCampaign",
				"err":    fmt.Sprint(r),
			}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	// Normalize weights
	normalizedWeight := roundToFourDecimals(100.0 / float64(len(shortlistedCampaigns)))
	for _, campaign := range shortlistedCampaigns {
		campaign.SetWeight(normalizedWeight)
	}

	// Convert campaigns to variations
	variations := []campaign.Variation{}
	for _, ruleCampaign := range shortlistedCampaigns {
		jsonData, err := json.Marshal(ruleCampaign)
		if err == nil {
			var variation campaign.Variation
			if json.Unmarshal(jsonData, &variation) == nil {
				variations = append(variations, variation)
			}
		}
	}

	SetCampaignAllocation(variations)

	bucketValue := decision_maker.CalculateBucketValue(GetBucketingSeed(context.GetID(), nil, &groupId))

	winnerVariation := GetVariation(variations, bucketValue)

	if winnerVariation != nil {
		winnerCampaignKey := ""
		if winnerVariation.GetType() == string(enums.CampaignTypeAB) {
			winnerCampaignKey = winnerVariation.GetKey()
		} else {
			winnerCampaignKey = winnerVariation.GetName() + "_" + winnerVariation.GetRuleKey()
		}
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_WINNER_CAMPAIGN"], map[string]interface{}{
			"campaignKey": winnerCampaignKey,
			"groupId":     strconv.Itoa(groupId),
			"userId":      context.GetID(),
			"algo":        "using random algorithm",
		}))

		// Store the result
		dataToStore := map[string]interface{}{
			enums.StorageFeatureKey.GetValue():            constants.VWOMetaMegKey + strconv.Itoa(groupId),
			enums.StorageUserID.GetValue():                context.GetID(),
			enums.DecisionExperimentID.GetValue():         winnerVariation.GetID(),
			enums.StorageExperimentKey.GetValue():         winnerVariation.GetKey(),
			enums.StorageExperimentVariationID.GetValue(): getVariationIDForPersonalizeFromVariation(winnerVariation),
		}

		storageDecorator := decorators.NewStorageDecorator()
		storageDecorator.SetDataInStorage(dataToStore, storageService, serviceContainer)

		if containsInt(calledCampaignIds, winnerVariation.GetID()) {
			return winnerVariation
		}
	} else {
		serviceContainer.GetLoggerService().Info("No winner campaign found for MEG group: " + strconv.Itoa(groupId))
	}

	return nil
}

// getCampaignUsingAdvancedAlgo advanced algorithm to find the winning campaign based on priority order and weighted random distribution and returns the winner variation
func (mu *MegUtil) getCampaignUsingAdvancedAlgo(
	serviceContainer interfaces.ServiceContainerInterface,
	shortlistedCampaigns []*campaign.Campaign,
	context *user.VWOContext,
	calledCampaignIds []int,
	groupId int,
	storageService interfaces.StorageServiceInterface,
) *campaign.Variation {
	var winnerCampaign *campaign.Variation
	found := false

	defer func() {
		if r := recover(); r != nil {
			serviceContainer.GetLoggerService().Error("ERROR_FINDING_WINNER_CAMPAIGN", map[string]interface{}{
				"method": "getCampaignUsingAdvancedAlgo",
				"err":    fmt.Sprint(r),
			}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	groups := serviceContainer.GetSettings().GetGroups()
	group, exists := groups[strconv.Itoa(groupId)]
	priorityOrder := []string{}
	wt := make(map[string]float64)

	if exists {
		if len(group.GetP()) > 0 {
			priorityOrder = group.GetP()
		}
		if len(group.GetWt()) > 0 {
			wt = group.GetWt()
		}
	}

	// Check priority order first
	for _, priorityID := range priorityOrder {
		for _, ruleCampaign := range shortlistedCampaigns {
			if strconv.Itoa(ruleCampaign.GetID()) == priorityID {
				jsonData, err := json.Marshal(ruleCampaign)
				if err == nil {
					json.Unmarshal(jsonData, &winnerCampaign)
				}
				found = true
				break
			} else if (strconv.Itoa(ruleCampaign.GetID()) + "_" + strconv.Itoa(ruleCampaign.GetVariations()[0].GetID())) == priorityID {
				jsonData, err := json.Marshal(ruleCampaign)
				if err == nil {
					json.Unmarshal(jsonData, &winnerCampaign)
				}
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// If no priority match, use weighted random
	if winnerCampaign == nil {
		participatingCampaignList := []*campaign.Campaign{}
		for _, ruleCampaign := range shortlistedCampaigns {
			campaignID := ruleCampaign.GetID()
			campaignIDStr := strconv.Itoa(campaignID)
			campaignWithVariationID := campaignIDStr + "_" + strconv.Itoa(ruleCampaign.GetVariations()[0].GetID())

			if weight, exists := wt[campaignIDStr]; exists {
				campaignCopy := *ruleCampaign
				campaignCopy.SetWeight(weight)
				participatingCampaignList = append(participatingCampaignList, &campaignCopy)
			} else if weight, exists := wt[campaignWithVariationID]; exists {
				campaignCopy := *ruleCampaign
				campaignCopy.SetWeight(weight)
				participatingCampaignList = append(participatingCampaignList, &campaignCopy)
			}
		}

		// Convert to variations
		variations := []campaign.Variation{}
		for _, participatingCampaign := range participatingCampaignList {
			jsonData, err := json.Marshal(participatingCampaign)
			if err == nil {
				var variation campaign.Variation
				if json.Unmarshal(jsonData, &variation) == nil {
					variations = append(variations, variation)
				}
			}
		}

		SetCampaignAllocation(variations)

		bucketValue := decision_maker.CalculateBucketValue(GetBucketingSeed(context.GetID(), nil, &groupId))
		winnerCampaign = GetVariation(variations, bucketValue)
	}

	if winnerCampaign != nil {
		winnerCampaignKey := ""
		if winnerCampaign.GetType() == string(enums.CampaignTypeAB) {
			winnerCampaignKey = winnerCampaign.GetKey()
		} else {
			winnerCampaignKey = winnerCampaign.GetName() + "_" + winnerCampaign.GetRuleKey()
		}
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["MEG_WINNER_CAMPAIGN"], map[string]interface{}{
			"campaignKey": winnerCampaignKey,
			"groupId":     strconv.Itoa(groupId),
			"userId":      context.GetID(),
			"algo":        "using advanced algorithm",
		}))

		// Store the result
		dataToStore := map[string]interface{}{
			enums.StorageFeatureKey.GetValue():            constants.VWOMetaMegKey + strconv.Itoa(groupId),
			enums.StorageUserID.GetValue():                context.GetID(),
			enums.DecisionExperimentID.GetValue():         winnerCampaign.GetID(),
			enums.StorageExperimentKey.GetValue():         winnerCampaign.GetKey(),
			enums.StorageExperimentVariationID.GetValue(): getVariationIDForPersonalizeFromVariation(winnerCampaign),
		}

		storageDecorator := decorators.NewStorageDecorator()
		storageDecorator.SetDataInStorage(dataToStore, storageService, serviceContainer)

		if containsInt(calledCampaignIds, winnerCampaign.GetID()) {
			return winnerCampaign
		}
	} else {
		serviceContainer.GetLoggerService().Info(log.BuildMessage(log.InfoLogMessagesEnum["NO_WINNER_CAMPAIGN_FOUND_FOR_MEG_GROUP"], map[string]interface{}{
			"groupId": strconv.Itoa(groupId),
		}))
	}

	return nil
}

// contains checks if a string slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// containsInt checks if an int slice contains an int
func containsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// roundToFourDecimals rounds to 4 decimal places
func roundToFourDecimals(value float64) float64 {
	return math.Round(value*10000) / 10000
}

// getVariationIDForPersonalizeFromVariation returns the variation ID for personalize campaigns
func getVariationIDForPersonalizeFromVariation(variation *campaign.Variation) int {
	if variation.GetType() == string(enums.CampaignTypePersonalize) {
		return variation.GetVariations()[0].GetID()
	}
	return -1
}
