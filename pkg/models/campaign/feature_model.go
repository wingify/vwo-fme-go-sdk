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

package campaign

// Feature represents a VWO feature
type Feature struct {
	ID                       int        `json:"id"`
	Key                      string     `json:"key"`
	Name                     string     `json:"name"`
	Type                     string     `json:"type"`
	IsEnabled                bool       `json:"isEnabled,omitempty"`
	Variables                []Variable `json:"variables,omitempty"`
	Metrics                  []Metric   `json:"metrics,omitempty"`
	Rules                    []Rule     `json:"rules,omitempty"`
	ImpactCampaign           *Campaign  `json:"impactCampaign,omitempty"`
	RulesLinkedCampaign      []Campaign `json:"rulesLinkedCampaign,omitempty"`
	IsGatewayServiceRequired bool       `json:"isGatewayServiceRequired,omitempty"`
	IsDebuggerEnabled        bool       `json:"isDebuggerEnabled,omitempty"`
}

// GetName returns the feature name
func (f *Feature) GetName() string {
	return f.Name
}

// GetType returns the feature type
func (f *Feature) GetType() string {
	return f.Type
}

// GetID returns the feature ID
func (f *Feature) GetID() int {
	return f.ID
}

// GetKey returns the feature key
func (f *Feature) GetKey() string {
	return f.Key
}

// GetRules returns the feature rules
func (f *Feature) GetRules() []Rule {
	return f.Rules
}

// GetImpactCampaign returns the impact campaign
func (f *Feature) GetImpactCampaign() *Campaign {
	return f.ImpactCampaign
}

// GetRulesLinkedCampaign returns the rules linked campaigns
func (f *Feature) GetRulesLinkedCampaign() []Campaign {
	return f.RulesLinkedCampaign
}

// SetRulesLinkedCampaign sets the rules linked campaigns
func (f *Feature) SetRulesLinkedCampaign(rulesLinkedCampaign []Campaign) {
	f.RulesLinkedCampaign = rulesLinkedCampaign
}

// GetMetrics returns the feature metrics
func (f *Feature) GetMetrics() []Metric {
	return f.Metrics
}

// GetIsGatewayServiceRequired returns if gateway service is required
func (f *Feature) GetIsGatewayServiceRequired() bool {
	return f.IsGatewayServiceRequired
}

// SetIsGatewayServiceRequired sets if gateway service is required
func (f *Feature) SetIsGatewayServiceRequired(isRequired bool) {
	f.IsGatewayServiceRequired = isRequired
}

// GetIsEnabled returns if the feature is enabled
func (f *Feature) GetIsEnabled() bool {
	return f.IsEnabled
}

// SetIsEnabled sets if the feature is enabled
func (f *Feature) SetIsEnabled(isEnabled bool) {
	f.IsEnabled = isEnabled
}

// GetVariables returns the feature variables
func (f *Feature) GetVariables() []Variable {
	return f.Variables
}

// SetVariables sets the feature variables
func (f *Feature) SetVariables(variables []Variable) {
	f.Variables = variables
}

// GetIsDebuggerEnabled returns if debugger is enabled
func (f *Feature) GetIsDebuggerEnabled() bool {
	return f.IsDebuggerEnabled
}

// SetIsDebuggerEnabled sets if debugger is enabled
func (f *Feature) SetIsDebuggerEnabled(isDebuggerEnabled bool) {
	f.IsDebuggerEnabled = isDebuggerEnabled
}
