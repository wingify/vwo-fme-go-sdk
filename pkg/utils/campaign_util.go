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
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// SetVariationAllocation sets the variation allocation for a given campaign based on its type
func SetVariationAllocation(campaignModel *campaign.Campaign, logManager interfaces.LoggerServiceInterface) {
	// Check if the campaign type is rollout or personalize
	if campaignModel.Type == string(enums.CampaignTypeRollout) || campaignModel.Type == string(enums.CampaignTypePersonalize) {
		handleRolloutCampaign(campaignModel, logManager)
	} else {
		currentAllocation := 0
		// Iterate over each variation in the campaign
		for i := range campaignModel.Variations {
			variation := &campaignModel.Variations[i]
			// Assign range values to the variation and update the current allocation
			stepFactor := AssignRangeValues(variation, currentAllocation)
			currentAllocation += stepFactor

			logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["VARIATION_RANGE_ALLOCATION"], map[string]interface{}{
				"campaignKey":     campaignModel.Key,
				"variationKey":    variation.Name,
				"variationWeight": strconv.FormatFloat(variation.Weight, 'f', -1, 64),
				"startRange":      strconv.Itoa(variation.StartRangeVariation),
				"endRange":        strconv.Itoa(variation.EndRangeVariation),
			}))
		}
	}
}

// AssignRangeValues assigns start and end range values to a variation based on its weight
func AssignRangeValues(data *campaign.Variation, currentAllocation int) int {
	// Calculate the bucket range based on the variation's weight
	stepFactor := getVariationBucketRange(data.Weight)

	// Set the start and end range of the variation
	if stepFactor > 0 {
		data.StartRangeVariation = currentAllocation + 1
		data.EndRangeVariation = currentAllocation + stepFactor
	} else {
		data.StartRangeVariation = -1
		data.EndRangeVariation = -1
	}
	return stepFactor
}

// ScaleVariationWeights scales the weights of variations to sum up to 100%
func ScaleVariationWeights(variations []campaign.Variation) {
	// Calculate the total weight of all variations
	totalWeight := 0.0
	for _, variation := range variations {
		totalWeight += variation.Weight
	}

	// If total weight is zero, assign equal weight to each variation
	if totalWeight == 0 {
		equalWeight := 100.0 / float64(len(variations))
		for i := range variations {
			variations[i].Weight = equalWeight
		}
	} else {
		// Scale each variation's weight to make the total 100%
		for i := range variations {
			variations[i].Weight = (variations[i].Weight / totalWeight) * 100
		}
	}
}

// GetBucketingSeed generates a bucketing seed based on user ID, campaign, and optional group ID
func GetBucketingSeed(userID string, campaignModel *campaign.Campaign, groupID *int) string {
	// Return a seed combining group ID and user ID if group ID is provided
	if groupID != nil {
		return fmt.Sprintf("%d_%s", *groupID, userID)
	}

	// Get campaign type
	campaignType := campaignModel.Type
	// Check if campaign type is rollout or personalize
	isRolloutOrPersonalize := campaignType == string(enums.CampaignTypeRollout) ||
		campaignType == string(enums.CampaignTypePersonalize)

	// Get salt based on campaign type
	var salt string
	if isRolloutOrPersonalize && len(campaignModel.Variations) > 0 {
		salt = campaignModel.Variations[0].Salt
	} else {
		salt = campaignModel.Salt
	}

	// If salt is not null and not empty, use salt else use campaign id
	var bucketKey string
	if salt != "" {
		bucketKey = fmt.Sprintf("%s_%s", salt, userID)
	} else {
		bucketKey = fmt.Sprintf("%d_%s", campaignModel.ID, userID)
	}

	return bucketKey
}

// GetVariationFromCampaignKey retrieves a variation by its ID within a specific campaign identified by its key
func GetVariationFromCampaignKey(settingsModel *settings.Settings, campaignKey string, variationID int) campaign.Variation {
	// Find the campaign by its key
	var foundCampaign *campaign.Campaign
	for i := range settingsModel.Campaigns {
		if settingsModel.Campaigns[i].Key == campaignKey {
			foundCampaign = &settingsModel.Campaigns[i]
			break
		}
	}

	if foundCampaign != nil {
		// Find the variation by its ID within the found campaign
		for i := range foundCampaign.Variations {
			if foundCampaign.Variations[i].ID == variationID {
				return foundCampaign.Variations[i]
			}
		}
	}
	return campaign.Variation{}
}

// SetCampaignAllocation sets the allocation ranges for a list of campaigns
func SetCampaignAllocation(campaigns []campaign.Variation) {
	currentAllocation := 0
	for i := range campaigns {
		// Assign range values to each campaign and update the current allocation
		stepFactor := assignRangeValuesMEG(&campaigns[i], currentAllocation)
		currentAllocation += stepFactor
	}
}

// GetGroupDetailsIfCampaignPartOfIt determines if a campaign is part of a group
func GetGroupDetailsIfCampaignPartOfIt(settingsModel *settings.Settings, campaignID int, variationID int) map[string]string {
	groupDetails := make(map[string]string)
	campaignToCheck := strconv.Itoa(campaignID)

	// If variationId is not -1, append it to campaignId
	if variationID != -1 {
		campaignToCheck = fmt.Sprintf("%s_%d", campaignToCheck, variationID)
	}

	if settingsModel != nil && settingsModel.CampaignGroups != nil {
		if groupID, exists := settingsModel.CampaignGroups[campaignToCheck]; exists {
			groupIDStr := strconv.Itoa(groupID)
			if group, exists := settingsModel.Groups[groupIDStr]; exists {
				groupDetails["groupId"] = groupIDStr
				groupDetails["groupName"] = group.Name
			}
		}
	}
	return groupDetails
}

// FindGroupsFeaturePartOf finds all groups associated with a feature specified by its key
func FindGroupsFeaturePartOf(settingsModel *settings.Settings, featureKey string) []map[string]string {
	// Initialize an array to store all rules for the given feature
	var ruleList []campaign.Rule

	for _, feature := range settingsModel.Features {
		if feature.Key == featureKey {
			for _, rule := range feature.Rules {
				// Add rule to the array if it's not already present
				exists := false
				for _, existingRule := range ruleList {
					if existingRule.CampaignID == rule.CampaignID && existingRule.VariationID == rule.VariationID {
						exists = true
						break
					}
				}
				if !exists {
					ruleList = append(ruleList, rule)
				}
			}
		}
	}

	// Initialize an array to store all groups associated with the feature
	var groups []map[string]string

	// Iterate over each rule to find the group details
	for _, rule := range ruleList {
		variationID := -1
		if rule.Type == string(enums.CampaignTypePersonalize) {
			variationID = rule.VariationID
		}

		group := GetGroupDetailsIfCampaignPartOfIt(settingsModel, rule.CampaignID, variationID)

		// Add group to the array if it's not already present
		if len(group) > 0 {
			exists := false
			for _, existingGroup := range groups {
				if existingGroup["groupId"] == group["groupId"] {
					exists = true
					break
				}
			}
			if !exists {
				groups = append(groups, group)
			}
		}
	}

	return groups
}

// GetCampaignsByGroupID retrieves campaigns by a specific group ID
func GetCampaignsByGroupID(settingsModel *settings.Settings, groupID int) []string {
	groupIDStr := strconv.Itoa(groupID)
	if group, exists := settingsModel.Groups[groupIDStr]; exists {
		return group.Campaigns
	}
	return []string{}
}

// GetFeatureKeysFromCampaignIDs retrieves feature keys from a list of campaign IDs
func GetFeatureKeysFromCampaignIDs(settingsModel *settings.Settings, campaignIDWithVariation []string) []string {
	var featureKeys []string

	for _, campaign := range campaignIDWithVariation {
		// Split key with _ to separate campaignId and variationId
		parts := strings.Split(campaign, "_")
		campaignID, _ := strconv.Atoi(parts[0])
		var variationID *int
		if len(parts) > 1 {
			val, _ := strconv.Atoi(parts[1])
			variationID = &val
		}

		// Iterate over each feature to find the feature key
		for _, feature := range settingsModel.Features {
			// Skip if feature key is already added
			alreadyAdded := false
			for _, key := range featureKeys {
				if key == feature.Key {
					alreadyAdded = true
					break
				}
			}
			if alreadyAdded {
				continue
			}

			for _, rule := range feature.Rules {
				if rule.CampaignID == campaignID {
					// Check if variationId is provided and matches the rule's variationId
					if variationID != nil {
						if rule.VariationID == *variationID {
							featureKeys = append(featureKeys, feature.Key)
						}
					} else {
						featureKeys = append(featureKeys, feature.Key)
					}
				}
			}
		}
	}

	return featureKeys
}

// GetCampaignIDsFromFeatureKey retrieves campaign IDs from a specific feature key
func GetCampaignIDsFromFeatureKey(settingsModel *settings.Settings, featureKey string) []int {
	var campaignIDs []int

	for _, feature := range settingsModel.Features {
		if feature.Key == featureKey {
			for _, rule := range feature.Rules {
				campaignIDs = append(campaignIDs, rule.CampaignID)
			}
		}
	}

	return campaignIDs
}

// GetRuleTypeUsingCampaignIDFromFeature retrieves the rule type using a campaign ID from a specific feature
func GetRuleTypeUsingCampaignIDFromFeature(feature *campaign.Feature, campaignID int) string {
	for _, rule := range feature.Rules {
		if rule.CampaignID == campaignID {
			return rule.Type
		}
	}
	return ""
}

// assignRangeValuesMEG assigns range values to a campaign based on its weight
func assignRangeValuesMEG(data *campaign.Variation, currentAllocation int) int {
	stepFactor := getVariationBucketRange(data.Weight)

	if stepFactor > 0 {
		data.StartRangeVariation = currentAllocation + 1
		data.EndRangeVariation = currentAllocation + stepFactor
	} else {
		data.StartRangeVariation = -1
		data.EndRangeVariation = -1
	}
	return stepFactor
}

// getVariationBucketRange calculates the bucket range for a variation based on its weight
func getVariationBucketRange(variationWeight float64) int {
	if variationWeight <= 0 {
		return 0
	}
	startRange := int(math.Ceil(variationWeight * 100))
	if startRange > constants.MaxTrafficValue {
		return constants.MaxTrafficValue
	}
	return startRange
}

// handleRolloutCampaign handles the rollout campaign by setting start and end ranges for all variations
func handleRolloutCampaign(campaignModel *campaign.Campaign, logManager interfaces.LoggerServiceInterface) {
	// Set start and end ranges for all variations in the campaign
	for i := range campaignModel.Variations {
		variation := &campaignModel.Variations[i]
		endRange := int(variation.Weight * 100)
		variation.StartRangeVariation = 1
		variation.EndRangeVariation = endRange

		logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["VARIATION_RANGE_ALLOCATION"], map[string]interface{}{
			"campaignKey":     campaignModel.Key,
			"variationKey":    variation.Name,
			"variationWeight": strconv.FormatFloat(variation.Weight, 'f', -1, 64),
			"startRange":      strconv.Itoa(variation.StartRangeVariation),
			"endRange":        strconv.Itoa(variation.EndRangeVariation),
		}))
	}
}

// GetCampaignLoggingKey returns the logging key used in messages
// For AB (FLAG_TESTING): use campaign.Key
// For others: use campaign.Name + "_" + campaign.RuleKey
func GetCampaignLoggingKey(c *campaign.Campaign) string {
	if c == nil {
		return ""
	}
	if c.Type == string(enums.CampaignTypeAB) {
		return c.Key
	}
	// safeguard if rule key is empty
	if c.RuleKey != "" {
		return c.Name + "_" + c.RuleKey
	}
	return c.Name
}

// GetCampaignKeyFromCampaignID retrieves the campaign key from a campaign ID
func GetCampaignKeyFromCampaignID(settingsModel *settings.Settings, campaignID int) string {
	for _, campaign := range settingsModel.Campaigns {
		if campaign.ID == campaignID {
			return campaign.Key
		}
	}
	return ""
}

// GetVariationNameFromCampaignIdAndVariationId retrieves the variation name from a campaign ID and variation ID
func GetVariationNameFromCampaignIdAndVariationId(settingsModel *settings.Settings, campaignID int, variationID int) string {
	for _, campaign := range settingsModel.Campaigns {
		if campaign.ID == campaignID {
			for _, variation := range campaign.Variations {
				if variation.ID == variationID {
					return variation.Name
				}
			}
		}
	}
	return ""
}

// GetCampaignTypeFromCampaignId retrieves the campaign type from a campaign ID
func GetCampaignTypeFromCampaignId(settingsModel *settings.Settings, campaignID int) string {
	for _, campaign := range settingsModel.Campaigns {
		if campaign.ID == campaignID {
			return campaign.Type
		}
	}
	return ""
}

// IsFeaturePresentInSettings checks if a feature is present in the settings
func IsFeaturePresentInSettings(settingsModel *settings.Settings, featureId int) bool {
	for _, feature := range settingsModel.Features {
		if feature.ID == featureId {
			return true
		}
	}
	return false
}
