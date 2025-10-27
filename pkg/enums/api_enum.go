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

package enums

// ApiEnum represents API method names
type ApiEnum string

const (
	ApiInit           ApiEnum = "init"
	ApiGetFlag        ApiEnum = "getFlag"
	ApiTrackEvent     ApiEnum = "trackEvent"
	ApiSetAttribute   ApiEnum = "setAttribute"
	ApiUpdateSettings ApiEnum = "updateSettings"
	ApiFlushEvents    ApiEnum = "flushEvents"
)

// CampaignTypeEnum represents campaign types
type CampaignTypeEnum string

const (
	CampaignTypeRollout     CampaignTypeEnum = "FLAG_ROLLOUT"
	CampaignTypeAB          CampaignTypeEnum = "FLAG_TESTING"
	CampaignTypePersonalize CampaignTypeEnum = "FLAG_PERSONALIZE"
)

// GetValue returns the string value of the campaign type enum
func (c CampaignTypeEnum) GetValue() string {
	return string(c)
}
