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

package testCases

import "github.com/wingify/vwo-fme-go-sdk/pkg/models/storage"

// Expectation represents the expected result for a test case
type Expectation struct {
	IsEnabled                 *bool                  `json:"isEnabled"`
	IntVariable               *float64               `json:"intVariable"`
	StringVariable            *string                `json:"stringVariable"`
	FloatVariable             *float64               `json:"floatVariable"`
	BooleanVariable           *bool                  `json:"booleanVariable"`
	JSONVariable              map[string]interface{} `json:"jsonVariable"`
	StorageData               *storage.StorageData   `json:"storageData"`
	ShouldReturnSameVariation *bool                  `json:"shouldReturnSameVariation"`
}
