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

// Variable represents a campaign variable
type Variable struct {
	ID    int         `json:"id"`
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// GetID returns the variable ID
func (vr *Variable) GetID() int {
	return vr.ID
}

// GetValue returns the variable value
func (vr *Variable) GetValue() interface{} {
	return vr.Value
}

// GetType returns the variable type
func (vr *Variable) GetType() string {
	return vr.Type
}

// GetKey returns the variable key
func (vr *Variable) GetKey() string {
	return vr.Key
}

// SetValue sets the variable value
func (vr *Variable) SetValue(value interface{}) {
	vr.Value = value
}

// SetKey sets the variable key
func (vr *Variable) SetKey(key string) {
	vr.Key = key
}

// SetType sets the variable type
func (vr *Variable) SetType(varType string) {
	vr.Type = varType
}
