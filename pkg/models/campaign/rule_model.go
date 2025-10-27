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

// Rule represents a feature rule
type Rule struct {
	Type        string `json:"type"`
	Status      bool   `json:"status"`
	VariationID int    `json:"variationId"`
	CampaignID  int    `json:"campaignId"`
	RuleKey     string `json:"ruleKey"`
}

// GetCampaignID returns the rule campaign ID
func (r *Rule) GetCampaignID() int {
	return r.CampaignID
}

// GetVariationID returns the rule variation ID
func (r *Rule) GetVariationID() int {
	return r.VariationID
}

// GetStatus returns the rule status
func (r *Rule) GetStatus() bool {
	return r.Status
}

// GetType returns the rule type
func (r *Rule) GetType() string {
	return r.Type
}

// GetRuleKey returns the rule key
func (r *Rule) GetRuleKey() string {
	return r.RuleKey
}
