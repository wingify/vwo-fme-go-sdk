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

package settings

import (
	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
)

// Settings represents VWO settings
type Settings struct {
	SDKKey              string                     `json:"sdkKey"`
	AccountID           int                        `json:"accountId"`
	Version             int                        `json:"version"`
	CollectionPrefix    string                     `json:"collectionPrefix,omitempty"`
	UsageStatsAccountID int                        `json:"usageStatsAccountId,omitempty"`
	Features            []campaign.Feature         `json:"features,omitempty"`
	Campaigns           []campaign.Campaign        `json:"campaigns,omitempty"`
	CampaignGroups      map[string]int             `json:"campaignGroups,omitempty"`
	Groups              map[string]campaign.Groups `json:"groups,omitempty"`
	PollInterval        int                        `json:"pollInterval,omitempty"`
	SDKMetaInfo         map[string]interface{}     `json:"sdkMetaInfo,omitempty"`
}

// GetFeatures returns the features
func (s *Settings) GetFeatures() []campaign.Feature {
	return s.Features
}

// GetCampaigns returns the campaigns
func (s *Settings) GetCampaigns() []campaign.Campaign {
	return s.Campaigns
}

// GetSDKKey returns the SDK key
func (s *Settings) GetSDKKey() string {
	return s.SDKKey
}

// GetAccountID returns the account ID
func (s *Settings) GetAccountID() int {
	return s.AccountID
}

// GetVersion returns the version
func (s *Settings) GetVersion() int {
	return s.Version
}

// GetCollectionPrefix returns the collection prefix
func (s *Settings) GetCollectionPrefix() string {
	return s.CollectionPrefix
}

// GetCampaignGroups returns the campaign groups
func (s *Settings) GetCampaignGroups() map[string]int {
	return s.CampaignGroups
}

// GetGroups returns the groups
func (s *Settings) GetGroups() map[string]campaign.Groups {
	return s.Groups
}

// SetPollInterval sets the poll interval
func (s *Settings) SetPollInterval(value int) {
	s.PollInterval = value
}

// GetPollInterval returns the poll interval
// Returns default poll interval if not set
func (s *Settings) GetPollInterval() int {
	if s.PollInterval == 0 {
		// If poll interval is not set (0), return default value
		return constants.DefaultPollInterval
	}
	return s.PollInterval
}

// GetUsageStatsAccountID returns the usage stats account ID
func (s *Settings) GetUsageStatsAccountID() int {
	return s.UsageStatsAccountID
}
