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

// Variable represents a variable object in the getVariables response.
type Variable struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Id    int         `json:"id"`
	Value interface{} `json:"value"`
}

// NewVariable creates a new Variable instance
func NewVariable(key string, value interface{}, varType string, id int) *Variable {
	return &Variable{
		Key:   key,
		Value: value,
		Type:  varType,
		Id:    id,
	}
}

// GetValue returns the variable value
func (v *Variable) GetValue() interface{} {
	return v.Value
}

// SetValue sets the variable value
func (v *Variable) SetValue(value interface{}) {
	v.Value = value
}

// GetType returns the variable type
func (v *Variable) GetType() string {
	return v.Type
}

// SetType sets the variable type
func (v *Variable) SetType(varType string) {
	v.Type = varType
}

// GetKey returns the variable key
func (v *Variable) GetKey() string {
	return v.Key
}

// SetKey sets the variable key
func (v *Variable) SetKey(key string) {
	v.Key = key
}

// GetID returns the variable ID
func (v *Variable) GetID() int {
	return v.Id
}

// SetID sets the variable ID
func (v *Variable) SetID(id int) {
	v.Id = id
}
