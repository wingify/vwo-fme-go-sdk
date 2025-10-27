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

import "github.com/wingify/vwo-fme-go-sdk/pkg/models/user"

// TestData represents a single test case
type TestData struct {
	Description string           `json:"description"`
	Settings    string           `json:"settings"`
	Context     *user.VWOContext `json:"context"`
	UserIds     []string         `json:"userIds"`
	Expectation *Expectation     `json:"expectation"`
	FeatureKey  string           `json:"featureKey"`
	FeatureKey2 string           `json:"featureKey2"`
}
