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

// GetFlagResponse exposes a minimal read-only contract for feature flag consumers
// so that only the whitelisted methods are available to SDK users.
type GetFlagResponse interface {
	IsEnabled() bool
	GetVariable(key string, defaultValue interface{}) interface{}
	GetVariables() []map[string]interface{}
}

// GetFlag represents the result of a feature flag evaluation
type GetFlag struct {
	Enabled   bool        `json:"isEnabled"`
	Variables []*Variable `json:"variables,omitempty"`
	Reason    string      `json:"reason,omitempty"`
}

// NewGetFlag creates a new GetFlag instance
func NewGetFlag() *GetFlag {
	return &GetFlag{
		Enabled:   false,
		Variables: make([]*Variable, 0),
	}
}

// IsEnabled returns whether the flag is enabled
func (gf *GetFlag) IsEnabled() bool {
	return gf.Enabled
}

// SetIsEnabled sets whether the flag is enabled
func (gf *GetFlag) SetIsEnabled(isEnabled bool) {
	gf.Enabled = isEnabled
}

// SetVariables sets the variables list
func (gf *GetFlag) SetVariables(variables []*Variable) {
	gf.Variables = variables
}

// GetVariablesValue returns the variables list
func (gf *GetFlag) GetVariablesValue() []*Variable {
	return gf.Variables
}

// GetVariable returns a specific variable value given a key, with a default fallback
func (gf *GetFlag) GetVariable(key string, defaultValue interface{}) interface{} {
	for _, variable := range gf.GetVariablesValue() {
		if variable.GetKey() == key {
			return variable.GetValue()
		}
	}
	return defaultValue
}

// GetVariables returns variables as a list of maps
func (gf *GetFlag) GetVariables() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, variable := range gf.GetVariablesValue() {
		result = append(result, gf.convertVariableModelToMap(variable))
	}
	return result
}

// convertVariableModelToMap converts a Variable model to a map
func (gf *GetFlag) convertVariableModelToMap(variable *Variable) map[string]interface{} {
	return map[string]interface{}{
		"key":   variable.GetKey(),
		"value": variable.GetValue(),
		"type":  variable.GetType(),
		"id":    variable.GetID(),
	}
}
