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
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
)

// DoesEventBelongToAnyFeature checks if an event belongs to any feature in the settings
func DoesEventBelongToAnyFeature(eventName string, settings *settings.Settings) bool {
	if settings == nil || settings.GetFeatures() == nil {
		return false
	}

	for _, feature := range settings.GetFeatures() {
		if feature.GetMetrics() != nil {
			for _, metric := range feature.GetMetrics() {
				if metric.GetIdentifier() == eventName {
					return true
				}
			}
		}
	}
	return false
}

// GetFeatureFromKey finds a feature by its key
func GetFeatureFromKey(settings *settings.Settings, featureKey string) *campaign.Feature {
	if settings == nil || settings.GetFeatures() == nil {
		return nil
	}

	for _, feature := range settings.GetFeatures() {
		if feature.GetKey() == featureKey {
			return &feature
		}
	}
	return nil
}

// GetSpecificRulesBasedOnType gets specific rules based on campaign type
func GetSpecificRulesBasedOnType(feature *campaign.Feature, campaignType enums.CampaignTypeEnum) []*campaign.Campaign {
	if feature == nil {
		return nil
	}

	rules := feature.GetRulesLinkedCampaign()
	if rules == nil {
		return nil
	}

	var result []*campaign.Campaign
	for _, rule := range rules {
		if rule.GetType() == campaignType.GetValue() {
			// Create a copy of the rule to avoid memory issues
			ruleCopy := rule
			result = append(result, &ruleCopy)
		}
	}
	return result
}

// GetAllExperimentRules gets all AB and Personalize rules
func GetAllExperimentRules(feature *campaign.Feature) []*campaign.Campaign {
	if feature == nil || feature.GetRulesLinkedCampaign() == nil {
		return nil
	}

	var result []*campaign.Campaign
	for _, rule := range feature.GetRulesLinkedCampaign() {
		if rule.GetType() == enums.CampaignTypeAB.GetValue() || rule.GetType() == enums.CampaignTypePersonalize.GetValue() {
			// Create a copy of the rule to avoid memory issues
			ruleCopy := rule
			result = append(result, &ruleCopy)
		}
	}
	return result
}
