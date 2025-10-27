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
	"regexp"
	"strings"
)

// CheckValuePresent checks if the actual values match the expected values specified in the map
func CheckValuePresent(expectedMap map[string][]string, actualMap map[string]string) bool {
	for key, actualValue := range actualMap {
		if expectedValues, exists := expectedMap[key]; exists {
			// Convert expectedValues to lowercase
			lowercaseExpectedValues := make([]string, len(expectedValues))
			for i, val := range expectedValues {
				lowercaseExpectedValues[i] = strings.ToLower(val)
			}

			// Handle wildcard patterns for all keys
			for _, val := range expectedValues {
				if strings.HasPrefix(val, "wildcard(") && strings.HasSuffix(val, ")") {
					// Extract pattern from wildcard string
					wildcardPattern := val[9 : len(val)-1]
					// Convert wildcard pattern to regex
					regexPattern := strings.ReplaceAll(wildcardPattern, "*", ".*")
					regex, err := regexp.Compile("(?i)" + regexPattern) // Case insensitive
					if err == nil && regex.MatchString(actualValue) {
						return true // Match found
					}
				}
			}

			// Direct value check for all keys
			if contains(lowercaseExpectedValues, strings.ToLower(strings.TrimSpace(actualValue))) {
				return true // Direct value match found
			}
		}
	}
	return false // No matches found
}

// ValuesMatch compares expected location values with user's location to determine a match
func ValuesMatch(expectedLocationMap map[string]interface{}, userLocation map[string]string) bool {
	for key, value := range expectedLocationMap {
		if userValue, exists := userLocation[key]; exists {
			normalizedValue1 := NormalizeValue(value)
			normalizedValue2 := NormalizeValue(userValue)
			if normalizedValue1 != normalizedValue2 {
				return false
			}
		} else {
			return false
		}
	}
	return true // If all values match, return true
}

// NormalizeValue normalizes a value to a consistent format for comparison
func NormalizeValue(value interface{}) string {
	if value == nil {
		return ""
	}
	// Convert to string
	str := ""
	switch v := value.(type) {
	case string:
		str = v
	default:
		str = strings.Trim(strings.TrimSpace(v.(string)), "\"")
	}
	// Remove quotes and trim whitespace
	str = strings.Trim(strings.TrimSpace(str), "\"")
	return str
}

// MatchWithRegex matches a string against a regular expression and returns the match result
func MatchWithRegex(str string, regexPattern string) bool {
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return false
	}
	return regex.MatchString(str)
}

// GetKeyValue extracts the first key-value pair from a map
func GetKeyValue(node map[string]interface{}) (string, interface{}) {
	for key, value := range node {
		return key, value
	}
	return "", nil
}

// contains checks if a string slice contains a specific value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
