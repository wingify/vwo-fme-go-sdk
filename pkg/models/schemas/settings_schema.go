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

package schemas

import (
	"fmt"
	"strings"

	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
)

// SettingsSchema represents the validation result for settings
type SettingsSchema struct {
	Valid  bool
	Errors []string
}

// NewSettingsSchema creates a new SettingsSchema with default valid state
func NewSettingsSchema() *SettingsSchema {
	return &SettingsSchema{
		Valid:  true,
		Errors: []string{},
	}
}

// NewSettingsSchemaWithValid creates a new SettingsSchema with specified valid state
func NewSettingsSchemaWithValid(valid bool) *SettingsSchema {
	return &SettingsSchema{
		Valid:  valid,
		Errors: []string{},
	}
}

// IsValid returns if the settings are valid
func (settingsSchema *SettingsSchema) IsValid() bool {
	return settingsSchema.Valid
}

// SetValid sets the validity of the settings
func (settingsSchema *SettingsSchema) SetValid(valid bool) {
	settingsSchema.Valid = valid
}

// GetErrors returns the errors in the settings
func (settingsSchema *SettingsSchema) GetErrors() []string {
	return settingsSchema.Errors
}

// AddError adds an error to the settings
func (settingsSchema *SettingsSchema) AddError(error string) {
	settingsSchema.Errors = append(settingsSchema.Errors, error)
	settingsSchema.Valid = false
}

// GetErrorsAsString returns the errors as a string
func (settingsSchema *SettingsSchema) GetErrorsAsString() string {
	return strings.Join(settingsSchema.Errors, "; ")
}

// IsSettingsValid returns if the settings are valid
func (settingsSchema *SettingsSchema) IsSettingsValid(settingsObj *settings.Settings) bool {
	return settingsSchema.ValidateSettings(settingsObj).IsValid()
}

// ValidateSettings validates the settings
func (settingsSchema *SettingsSchema) ValidateSettings(settingsObj *settings.Settings) *SettingsSchema {
	result := NewSettingsSchema()

	defer func() {
		if r := recover(); r != nil {
			result.AddError(fmt.Sprintf("Error validating settings: %v", r))
			result.SetValid(false)
		}
	}()

	if settingsObj == nil {
		result.AddError("Settings object is null")
		return result
	}

	// Validate Settings fields
	if settingsObj.Version == 0 {
		result.AddError("Settings version is null")
	}

	if settingsObj.AccountID == 0 {
		result.AddError("Settings accountId is null")
	}

	if settingsObj.Campaigns == nil {
		result.AddError("Settings campaigns list is null")
	} else {
		for i, camp := range settingsObj.Campaigns {
			campaignResult := settingsSchema.validateCampaign(&camp, i)
			if !campaignResult.IsValid() {
				result.Errors = append(result.Errors, campaignResult.GetErrors()...)
				result.SetValid(false)
			}
		}
	}

	if settingsObj.Features != nil {
		for i, feat := range settingsObj.Features {
			featureResult := settingsSchema.validateFeature(&feat, i)
			if !featureResult.IsValid() {
				result.Errors = append(result.Errors, featureResult.GetErrors()...)
				result.SetValid(false)
			}
		}
	}

	return result
}

// validateCampaign validates the campaign
func (settingsSchema *SettingsSchema) validateCampaign(camp *campaign.Campaign, index int) *SettingsSchema {
	result := NewSettingsSchema()
	prefix := fmt.Sprintf("Campaign[%d]: ", index)

	if camp == nil {
		result.AddError(prefix + "Campaign object is null")
		return result
	}

	if camp.ID == 0 {
		result.AddError(prefix + "Campaign id is null")
	}

	if camp.Type == "" {
		result.AddError(prefix + "Campaign type is null")
	}

	if camp.Key == "" {
		result.AddError(prefix + "Campaign key is null")
	}

	// Note: Status field validation skipped as it's not present in current Campaign model

	if camp.Name == "" {
		result.AddError(prefix + "Campaign name is null")
	}

	if camp.Variations == nil {
		result.AddError(prefix + "Campaign variations list is null")
	} else if len(camp.Variations) == 0 {
		result.AddError(prefix + "Campaign variations list is empty")
	} else {
		for i, variation := range camp.Variations {
			variationResult := settingsSchema.validateCampaignVariation(&variation, index, i)
			if !variationResult.IsValid() {
				result.Errors = append(result.Errors, variationResult.GetErrors()...)
				result.SetValid(false)
			}
		}
	}

	return result
}

// validateCampaignVariation validates the campaign variation
func (settingsSchema *SettingsSchema) validateCampaignVariation(variation *campaign.Variation, campaignIndex int, variationIndex int) *SettingsSchema {
	result := NewSettingsSchema()
	prefix := fmt.Sprintf("Campaign[%d].Variation[%d]: ", campaignIndex, variationIndex)

	if variation == nil {
		result.AddError(prefix + "Variation object is null")
		return result
	}

	if variation.ID == 0 {
		result.AddError(prefix + "Variation id is null")
	}

	if variation.Name == "" {
		result.AddError(prefix + "Variation name is null")
	}

	if variation.Weight == 0 {
		result.AddError(prefix + "Variation weight is empty")
	}

	if variation.Variables != nil {
		for i, variable := range variation.Variables {
			variableResult := settingsSchema.validateVariableObject(&variable, fmt.Sprintf("Campaign[%d].Variation[%d].Variable[%d]", campaignIndex, variationIndex, i))
			if !variableResult.IsValid() {
				result.Errors = append(result.Errors, variableResult.GetErrors()...)
				result.SetValid(false)
			}
		}
	}

	return result
}

// validateVariableObject validates the variable object
func (settingsSchema *SettingsSchema) validateVariableObject(variable *campaign.Variable, context string) *SettingsSchema {
	result := NewSettingsSchema()
	prefix := context + ": "

	if variable == nil {
		result.AddError(prefix + "Variable object is null")
		return result
	}

	if variable.ID == 0 {
		result.AddError(prefix + "Variable id is null")
	}

	if variable.Type == "" {
		result.AddError(prefix + "Variable type is null")
	}

	if variable.Key == "" {
		result.AddError(prefix + "Variable key is null")
	}

	if variable.Value == nil {
		result.AddError(prefix + "Variable value is null")
	}

	return result
}

// validateFeature validates the feature
func (settingsSchema *SettingsSchema) validateFeature(feature *campaign.Feature, index int) *SettingsSchema {
	result := NewSettingsSchema()
	prefix := fmt.Sprintf("Feature[%d]: ", index)

	if feature == nil {
		result.AddError(prefix + "Feature object is null")
		return result
	}

	if feature.ID == 0 {
		result.AddError(prefix + "Feature id is null")
	}

	if feature.Key == "" {
		result.AddError(prefix + "Feature key is null")
	}

	// Note: Status field validation skipped as it's not present in current Feature model

	if feature.Name == "" {
		result.AddError(prefix + "Feature name is null")
	}

	if feature.Type == "" {
		result.AddError(prefix + "Feature type is null")
	}

	if feature.Metrics == nil {
		result.AddError(prefix + "Feature metrics list is null")
	} else if len(feature.Metrics) == 0 {
		result.AddError(prefix + "Feature metrics list is empty")
	} else {
		for i, metric := range feature.Metrics {
			metricResult := settingsSchema.validateCampaignMetric(&metric, index, i)
			if !metricResult.IsValid() {
				result.Errors = append(result.Errors, metricResult.GetErrors()...)
				result.SetValid(false)
			}
		}
	}

	if feature.Rules != nil {
		for i, rule := range feature.Rules {
			ruleResult := settingsSchema.validateRule(&rule, index, i)
			if !ruleResult.IsValid() {
				result.Errors = append(result.Errors, ruleResult.GetErrors()...)
				result.SetValid(false)
			}
		}
	}

	// Note: Variables validation skipped as Feature model doesn't have Variables field
	// Variables are in Campaign model in the current implementation

	return result
}

// validateCampaignMetric validates the campaign metric
func (settingsSchema *SettingsSchema) validateCampaignMetric(metric *campaign.Metric, featureIndex int, metricIndex int) *SettingsSchema {
	result := NewSettingsSchema()
	prefix := fmt.Sprintf("Feature[%d].Metric[%d]: ", featureIndex, metricIndex)

	if metric == nil {
		result.AddError(prefix + "Metric object is null")
		return result
	}

	if metric.ID == 0 {
		result.AddError(prefix + "Metric id is null")
	}

	if metric.Type == "" {
		result.AddError(prefix + "Metric type is null")
	}

	if metric.Identifier == "" {
		result.AddError(prefix + "Metric identifier is null")
	}

	return result
}

// validateRule validates the rule
func (settingsSchema *SettingsSchema) validateRule(rule *campaign.Rule, featureIndex int, ruleIndex int) *SettingsSchema {
	result := NewSettingsSchema()
	prefix := fmt.Sprintf("Feature[%d].Rule[%d]: ", featureIndex, ruleIndex)

	if rule == nil {
		result.AddError(prefix + "Rule object is null")
		return result
	}

	if rule.Type == "" {
		result.AddError(prefix + "Rule type is null")
	}

	if rule.RuleKey == "" {
		result.AddError(prefix + "Rule ruleKey is null")
	}

	if rule.CampaignID == 0 {
		result.AddError(prefix + "Rule campaignId is null")
	}

	return result
}
