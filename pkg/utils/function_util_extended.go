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
	"encoding/json"

	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
)

// CloneObject creates a deep copy of an object using JSON marshaling/unmarshaling
func CloneObject(obj interface{}) interface{} {
	// Convert to JSON
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	// Determine the type and create appropriate instance
	switch obj.(type) {
	case *campaign.Campaign:
		var cloned campaign.Campaign
		err = json.Unmarshal(jsonBytes, &cloned)
		if err != nil {
			return nil
		}
		return &cloned
	case campaign.Campaign:
		var cloned campaign.Campaign
		err = json.Unmarshal(jsonBytes, &cloned)
		if err != nil {
			return nil
		}
		return cloned
	case *campaign.Variation:
		var cloned campaign.Variation
		err = json.Unmarshal(jsonBytes, &cloned)
		if err != nil {
			return nil
		}
		return &cloned
	case campaign.Variation:
		var cloned campaign.Variation
		err = json.Unmarshal(jsonBytes, &cloned)
		if err != nil {
			return nil
		}
		return cloned
	default:
		// For other types, return a generic map
		var cloned map[string]interface{}
		err = json.Unmarshal(jsonBytes, &cloned)
		if err != nil {
			return nil
		}
		return cloned
	}
}
