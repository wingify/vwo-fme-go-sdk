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
	"regexp"
	"strings"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// ProcessSettings processes the settings file and modifies it as required
// This method is called before the settings are used by the SDK
// It sets the variation allocation for each campaign
// It adds linked campaigns to each feature in the settings based on rules
// It adds isGatewayServiceRequired flag to each feature in the settings based on pre segmentation
func ProcessSettings(settingsModel *settings.Settings, logManager interfaces.LoggerServiceInterface) {
	campaigns := settingsModel.Campaigns

	for i := range campaigns {
		SetVariationAllocation(&campaigns[i], logManager)
	}

	addLinkedCampaignsToSettings(settingsModel)
	addIsGatewayServiceRequiredFlag(settingsModel)
}

// addLinkedCampaignsToSettings adds linked campaigns to each feature in the settings based on rules
func addLinkedCampaignsToSettings(settingsModel *settings.Settings) {
	// Create a map for quick access to campaigns by ID
	campaignMap := make(map[int]*campaign.Campaign)
	for i := range settingsModel.Campaigns {
		campaignMap[settingsModel.Campaigns[i].ID] = &settingsModel.Campaigns[i]
	}

	// Loop over all features
	for featureIndex := range settingsModel.Features {
		feature := &settingsModel.Features[featureIndex]
		var rulesLinkedCampaignModel []campaign.Campaign

		for _, rule := range feature.Rules {
			originalCampaign, exists := campaignMap[rule.CampaignID]
			if !exists {
				continue
			}

			// Create a copy of the campaign
			linkedCampaign := campaign.Campaign{
				ID:                       originalCampaign.ID,
				Segments:                 originalCampaign.Segments,
				Salt:                     originalCampaign.Salt,
				PercentTraffic:           originalCampaign.PercentTraffic,
				IsUserListEnabled:        originalCampaign.IsUserListEnabled,
				Key:                      originalCampaign.Key,
				Type:                     originalCampaign.Type,
				Name:                     originalCampaign.Name,
				IsForcedVariationEnabled: originalCampaign.IsForcedVariationEnabled,
				Variations:               make([]campaign.Variation, len(originalCampaign.Variations)),
				Metrics:                  originalCampaign.Metrics,
				Variables:                originalCampaign.Variables,
				VariationID:              originalCampaign.VariationID,
				CampaignID:               originalCampaign.CampaignID,
				RuleKey:                  rule.RuleKey,
			}

			// Deep copy variations
			copy(linkedCampaign.Variations, originalCampaign.Variations)

			// If a variationId is specified, find and add the variation
			if rule.VariationID != 0 {
				for _, variation := range linkedCampaign.Variations {
					if variation.ID == rule.VariationID {
						linkedCampaign.Variations = []campaign.Variation{variation}
						break
					}
				}
			}

			rulesLinkedCampaignModel = append(rulesLinkedCampaignModel, linkedCampaign)
		}

		// Assign the linked campaigns to the feature
		feature.SetRulesLinkedCampaign(rulesLinkedCampaignModel)
	}
}

// addIsGatewayServiceRequiredFlag adds isGatewayServiceRequired flag to each feature in the settings based on pre segmentation
func addIsGatewayServiceRequiredFlag(settingsModel *settings.Settings) {
	// Updated pattern without using lookbehind
	patternString := `\b(country|region|city|os|device_type|browser_string|ua|browser_version|os_version)\b|"custom_variable"\s*:\s*\{\s*"name"\s*:\s*"inlist\([^)]*\)"`
	pattern := regexp.MustCompile(patternString)

	for featureIndex := range settingsModel.Features {
		feature := &settingsModel.Features[featureIndex]
		rules := feature.RulesLinkedCampaign

		for _, rule := range rules {
			var segments interface{}

			if rule.Type == string(enums.CampaignTypeRollout) || rule.Type == string(enums.CampaignTypePersonalize) {
				if len(rule.Variations) > 0 {
					segments = rule.Variations[0].Segments
				}
			} else {
				segments = rule.Segments
			}

			if segments != nil {
				jsonSegments, err := json.Marshal(segments)
				if err != nil {
					continue
				}

				jsonSegmentsStr := string(jsonSegments)
				matches := pattern.FindAllStringSubmatchIndex(jsonSegmentsStr, -1)
				foundMatch := false

				for _, match := range matches {
					matchStr := jsonSegmentsStr[match[0]:match[1]]

					// Check if it's one of the specific keywords
					if regexp.MustCompile(`\b(country|region|city|os|device_type|browser_string|ua|browser_version|os_version)\b`).MatchString(matchStr) {
						// Check if within "custom_variable" block
						if !isWithinCustomVariable(match[0], jsonSegmentsStr) {
							foundMatch = true
							break
						}
					} else {
						foundMatch = true
						break
					}
				}

				if foundMatch {
					feature.SetIsGatewayServiceRequired(true)
					break
				}
			}
		}
	}
}

// isWithinCustomVariable checks if a match is within "custom_variable"
func isWithinCustomVariable(startIndex int, input string) bool {
	index := strings.LastIndex(input[:startIndex], `"custom_variable"`)
	if index == -1 {
		return false
	}

	closingBracketIndex := strings.Index(input[index:], "}")
	return closingBracketIndex != -1 && startIndex < (index+closingBracketIndex)
}
