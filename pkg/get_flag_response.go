/**
 * Copyright 2024 Wingify Software Pvt. Ltd.
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
package sdk

import (
	"encoding/json"
	"errors"
)

// GetFlagResponse represents the response from the GetFlag API.
type GetFlagResponse struct {
	isEnabled bool
	variables []interface{}
}

// Variable represents a variable object in the getVariables response.
type Variable struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Id    int         `json:"id"`
	Value interface{} `json:"value"`
}

// NewGetFlagResponse initializes a new instance of GetFlagResponse.
// It takes a JSON-encoded response and returns a GetFlagResponse object.
func NewGetFlagResponse(response []byte) (*GetFlagResponse, error) {
	var data map[string]interface{}
	err := json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}

	isEnabled, ok := data["isEnabled"].(bool)
	if !ok {
		return nil, errors.New("invalid or missing isEnabled key in response")
	}

	// variables is a list of map containing the key value pair of the variables
	variables, ok := data["variables"].([]interface{})
	if !ok {
		return nil, errors.New("invalid or missing variables key in response")
	}

	return &GetFlagResponse{
		isEnabled: isEnabled,
		variables: variables,
	}, nil
}

// IsEnabled returns the value of isEnabled.
// It indicates whether the flag is enabled or not.
func (r *GetFlagResponse) IsEnabled() bool {
	return r.isEnabled
}

// GetVariables returns the value of variables.
// It provides a list of variables associated with the flag.
func (r *GetFlagResponse) GetVariables() []interface{} {
	return r.variables
}

// GetVariable returns the value of the variable with the given key or the default value if not found.
// It searches for a variable by its key in the list of variables.
func (r *GetFlagResponse) GetVariable(key string, defaultValue interface{}) interface{} {
	for _, v := range r.variables {
		varMap, ok := v.(map[string]interface{})
		if ok {
			if varMap["key"] == key {
				return varMap["value"]
			}
		}
	}
	return defaultValue
}
