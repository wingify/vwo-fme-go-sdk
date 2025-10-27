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

// Campaign represents a VWO campaign
type Campaign struct {
	ID                       int                    `json:"id"`
	Segments                 map[string]interface{} `json:"segments,omitempty"`
	Salt                     string                 `json:"salt,omitempty"`
	PercentTraffic           int                    `json:"percentTraffic,omitempty"`
	IsUserListEnabled        bool                   `json:"isUserListEnabled,omitempty"`
	Key                      string                 `json:"key"`
	Type                     string                 `json:"type"`
	Name                     string                 `json:"name"`
	IsForcedVariationEnabled bool                   `json:"isForcedVariationEnabled,omitempty"`
	Variations               []Variation            `json:"variations,omitempty"`
	Metrics                  []Metric               `json:"metrics,omitempty"`
	Variables                []Variable             `json:"variables,omitempty"`
	VariationID              int                    `json:"variationId,omitempty"`
	CampaignID               int                    `json:"campaignId,omitempty"`
	RuleKey                  string                 `json:"ruleKey,omitempty"`
	StartRangeVariation      int                    `json:"startRangeVariation,omitempty"`
	EndRangeVariation        int                    `json:"endRangeVariation,omitempty"`
	Weight                   float64                `json:"weight,omitempty"`
}

// GetID returns the campaign ID
func (c *Campaign) GetID() int {
	return c.ID
}

// GetName returns the campaign name
func (c *Campaign) GetName() string {
	return c.Name
}

// GetSegments returns the campaign segments
func (c *Campaign) GetSegments() map[string]interface{} {
	return c.Segments
}

// GetTraffic returns the campaign traffic percentage
func (c *Campaign) GetTraffic() int {
	return c.PercentTraffic
}

// GetType returns the campaign type
func (c *Campaign) GetType() string {
	return c.Type
}

// GetIsForcedVariationEnabled returns if forced variation is enabled
func (c *Campaign) GetIsForcedVariationEnabled() bool {
	return c.IsForcedVariationEnabled
}

// GetKey returns the campaign key
func (c *Campaign) GetKey() string {
	return c.Key
}

// GetVariations returns the campaign variations
func (c *Campaign) GetVariations() []Variation {
	return c.Variations
}

// GetCampaignID returns the campaign ID
func (c *Campaign) GetCampaignID() int {
	return c.CampaignID
}

// GetIsUserListEnabled returns if user list is enabled
func (c *Campaign) GetIsUserListEnabled() bool {
	return c.IsUserListEnabled
}

// GetMetrics returns the campaign metrics
func (c *Campaign) GetMetrics() []Metric {
	return c.Metrics
}

// GetVariables returns the campaign variables
func (c *Campaign) GetVariables() []Variable {
	return c.Variables
}

// GetRuleKey returns the rule key
func (c *Campaign) GetRuleKey() string {
	return c.RuleKey
}

// GetSalt returns the campaign salt
func (c *Campaign) GetSalt() string {
	return c.Salt
}

// GetStartRangeVariation returns the start range variation
func (c *Campaign) GetStartRangeVariation() int {
	return c.StartRangeVariation
}

// GetEndRangeVariation returns the end range variation
func (c *Campaign) GetEndRangeVariation() int {
	return c.EndRangeVariation
}

// GetWeight returns the campaign weight
func (c *Campaign) GetWeight() float64 {
	return c.Weight
}

// SetWeight sets the campaign weight
func (c *Campaign) SetWeight(weight float64) {
	c.Weight = weight
}

// SetStartRangeVariation sets the start range variation
func (c *Campaign) SetStartRangeVariation(startRangeVariation int) {
	c.StartRangeVariation = startRangeVariation
}

// SetEndRangeVariation sets the end range variation
func (c *Campaign) SetEndRangeVariation(endRangeVariation int) {
	c.EndRangeVariation = endRangeVariation
}
