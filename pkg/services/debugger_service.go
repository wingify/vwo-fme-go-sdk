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

package services

// DebuggerService handles debug event properties for different categories
type DebuggerService struct {
	// Map of debug event props for each category
	// key: category name
	// value: map of debug props
	debugEventProps map[string]map[string]interface{}

	// Map of standard debug props
	// key: prop name
	// value: prop value
	standardDebugProps map[string]interface{}
}

// NewDebuggerService creates a new DebuggerService instance
func NewDebuggerService() *DebuggerService {
	return &DebuggerService{
		debugEventProps:    make(map[string]map[string]interface{}),
		standardDebugProps: make(map[string]interface{}),
	}
}

// AddStandardDebugProps adds a map of standard debug props. This is used to add props that are common to all categories.
// @param standardDebugProps Map of standard debug props
func (ds *DebuggerService) AddStandardDebugProps(standardDebugProps map[string]interface{}) {
	for key, value := range standardDebugProps {
		ds.standardDebugProps[key] = value
	}
}

// AddStandardDebugProp adds a single standard debug prop.
// @param key Prop name
// @param value Prop value
func (ds *DebuggerService) AddStandardDebugProp(key string, value interface{}) {
	ds.standardDebugProps[key] = value
}

// GetStandardDebugProps returns the standard debug props.
// @return Map of standard debug props
func (ds *DebuggerService) GetStandardDebugProps() map[string]interface{} {
	return ds.standardDebugProps
}

// AddCategoryDebugProps adds a map of debug props to a specific category.
// @param category Category name
// @param eventProps Map of debug props
func (ds *DebuggerService) AddCategoryDebugProps(category string, eventProps map[string]interface{}) {
	ds.debugEventProps[category] = eventProps
}

// AddCategoryDebugProp adds a single debug prop to a specific category.
// @param category Category name
// @param key Prop name
// @param value Prop value
func (ds *DebuggerService) AddCategoryDebugProp(category string, key string, value interface{}) {
	if ds.debugEventProps[category] == nil {
		ds.debugEventProps[category] = make(map[string]interface{})
	}
	ds.debugEventProps[category][key] = value
}

// GetDebugEventProps returns a map of all debug event props for a specific category, with standard props merged into each category.
// Category-specific props override standard props.
// @param category Category name
// @return Map of debug event props for a specific category
func (ds *DebuggerService) GetDebugEventProps(category string) map[string]interface{} {
	// Copy all categories and merge standard props into each
	categoryProps := make(map[string]interface{})

	// Add standard props first
	for key, value := range ds.standardDebugProps {
		categoryProps[key] = value
	}

	// Category-specific props override standard ones
	if categoryEventProps, exists := ds.debugEventProps[category]; exists {
		for key, value := range categoryEventProps {
			categoryProps[key] = value
		}
	}

	return categoryProps
}

// Clear clears all debug event props and standard props.
func (ds *DebuggerService) Clear() {
	ds.debugEventProps = make(map[string]map[string]interface{})
	ds.standardDebugProps = make(map[string]interface{})
}
