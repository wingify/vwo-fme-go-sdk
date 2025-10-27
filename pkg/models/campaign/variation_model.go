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

// Variation represents a campaign variation
type Variation struct {
	ID                  int                    `json:"id"`
	Key                 string                 `json:"key"`
	Name                string                 `json:"name,omitempty"`
	Weight              float64                `json:"weight"`
	RuleKey             string                 `json:"ruleKey,omitempty"`
	Salt                string                 `json:"salt,omitempty"`
	Type                string                 `json:"type,omitempty"`
	StartRangeVariation int                    `json:"startRangeVariation,omitempty"`
	EndRangeVariation   int                    `json:"endRangeVariation,omitempty"`
	Variables           []Variable             `json:"variables,omitempty"`
	Variations          []Variation            `json:"variations,omitempty"`
	Segments            map[string]interface{} `json:"segments,omitempty"`
}

// GetID returns the variation ID
func (v *Variation) GetID() int {
	if v == nil {
		return 0
	}
	return v.ID
}

// GetKey returns the variation key
func (v *Variation) GetKey() string {
	if v == nil {
		return ""
	}
	return v.Key
}

// GetWeight returns the variation weight
func (v *Variation) GetWeight() float64 {
	if v == nil {
		return 0
	}
	return v.Weight
}

// GetSegments returns the variation segments
func (v *Variation) GetSegments() map[string]interface{} {
	if v == nil {
		return nil
	}
	return v.Segments
}

// GetStartRangeVariation returns the start range
func (v *Variation) GetStartRangeVariation() int {
	if v == nil {
		return 0
	}
	return v.StartRangeVariation
}

// GetEndRangeVariation returns the end range
func (v *Variation) GetEndRangeVariation() int {
	if v == nil {
		return 0
	}
	return v.EndRangeVariation
}

// GetVariables returns the variation variables
func (v *Variation) GetVariables() []Variable {
	if v == nil {
		return []Variable{}
	}
	if v.Variables == nil {
		return []Variable{}
	}
	return v.Variables
}

// GetVariations returns nested variations
func (v *Variation) GetVariations() []Variation {
	return v.Variations
}

// GetType returns the variation type
func (v *Variation) GetType() string {
	return v.Type
}

// GetSalt returns the variation salt
func (v *Variation) GetSalt() string {
	return v.Salt
}

// GetRuleKey returns the variation rule key
func (v *Variation) GetRuleKey() string {
	return v.RuleKey
}

// GetName returns the variation name
func (v *Variation) GetName() string {
	return v.Name
}

// SetStartRange sets the start range for a variation
func (v *Variation) SetStartRange(startRange int) {
	v.StartRangeVariation = startRange
}

// SetEndRange sets the end range for a variation
func (v *Variation) SetEndRange(endRange int) {
	v.EndRangeVariation = endRange
}

// SetWeight sets the weight for a variation
func (v *Variation) SetWeight(weight float64) {
	v.Weight = weight
}
