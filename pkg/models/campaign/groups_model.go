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

// Groups represents a group of campaigns
type Groups struct {
	Name      string             `json:"name"`
	Campaigns []string           `json:"campaigns,omitempty"`
	Et        int                `json:"et,omitempty"` // Algorithm type
	P         []string           `json:"p,omitempty"`  // Priority
	Wt        map[string]float64 `json:"wt,omitempty"` // Weight
}

// GetName returns the group name
func (g *Groups) GetName() string {
	return g.Name
}

// SetName sets the group name
func (g *Groups) SetName(name string) {
	g.Name = name
}

// GetCampaigns returns the campaigns in the group
func (g *Groups) GetCampaigns() []string {
	return g.Campaigns
}

// SetCampaigns sets the campaigns in the group
func (g *Groups) SetCampaigns(campaigns []string) {
	g.Campaigns = campaigns
}

// GetEt returns the algorithm type (default to 1 for random if not set)
func (g *Groups) GetEt() int {
	if g.Et == 0 {
		return 1 // Default to random
	}
	return g.Et
}

// SetEt sets the algorithm type
func (g *Groups) SetEt(et int) {
	g.Et = et
}

// GetP returns the priority
func (g *Groups) GetP() []string {
	return g.P
}

// SetP sets the priority
func (g *Groups) SetP(p []string) {
	g.P = p
}

// GetWt returns the weight map
func (g *Groups) GetWt() map[string]float64 {
	return g.Wt
}

// SetWt sets the weight map
func (g *Groups) SetWt(wt map[string]float64) {
	g.Wt = wt
}
