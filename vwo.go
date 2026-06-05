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

package vwo

import (
	wingify "github.com/wingify/wingify-fme-go-sdk"
	"github.com/wingify/wingify-fme-go-sdk/pkg/enums"
)

// VWOClient is the legacy compatibility alias for the Wingify FME client.
type VWOClient = wingify.WingifyClient

// Init initializes the VWO FME client using the Wingify SDK with the vwo host profile.
func Init(options map[string]interface{}) (*VWOClient, error) {
	if options == nil {
		options = map[string]interface{}{}
	}
	options[enums.OptionHostProfile.GetValue()] = "vwo"
	return wingify.Init(options)
}

// GetUUID generates a UUID for a user based on their userId and accountId.
func GetUUID(userID string, accountID string) (string, error) {
	return wingify.GetUUID(userID, accountID)
}
