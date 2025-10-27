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

package models

// GlobalRequestModel represents global request configuration
type GlobalRequestModel struct {
	BaseURL           string
	Timeout           int
	Query             map[string]interface{}
	Body              map[string]interface{}
	Headers           map[string]string
	IsDevelopmentMode bool
}

// NewGlobalRequestModel creates a new GlobalRequestModel with default values
func NewGlobalRequestModel(
	baseURL string,
	query map[string]interface{},
	body map[string]interface{},
	headers map[string]string,
) *GlobalRequestModel {
	return &GlobalRequestModel{
		BaseURL: baseURL,
		Timeout: 3000, // Default timeout 3 seconds
		Query:   query,
		Body:    body,
		Headers: headers,
	}
}
