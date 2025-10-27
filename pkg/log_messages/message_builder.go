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

package log_messages

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`\{([0-9a-zA-Z_]+)\}`)

// BuildMessage constructs a message by replacing placeholders in a template with corresponding values from a data map.
// Placeholders are in the format {key}.
//
// Parameters:
//   - template: The message template containing placeholders in the format `{key}`
//   - data: A map containing keys and values used to replace the placeholders in the template
//
// Returns:
//   - The constructed message with all placeholders replaced by their corresponding values from the data map
func BuildMessage(template string, data map[string]interface{}) string {
	if data == nil {
		data = make(map[string]interface{})
	}

	// Replace all placeholders with their corresponding values
	result := placeholderRegex.ReplaceAllStringFunc(template, func(match string) string {
		// Extract the key from the placeholder (remove { and })
		key := strings.Trim(match, "{}")

		// Retrieve the value from the data map
		value, exists := data[key]

		// If the key does not exist or the value is nil, return an empty string
		if !exists || value == nil {
			return ""
		}

		// Convert value to string based on its type
		switch v := value.(type) {
		case string:
			return v
		case int, int8, int16, int32, int64:
			return strings.TrimSpace(strings.Replace(match, key, fmt.Sprint(v), 1))
		case uint, uint8, uint16, uint32, uint64:
			return strings.TrimSpace(strings.Replace(match, key, fmt.Sprint(v), 1))
		case float32, float64:
			return strings.TrimSpace(strings.Replace(match, key, fmt.Sprint(v), 1))
		case bool:
			if v {
				return "true"
			}
			return "false"
		case func() string:
			// If the value is a function, evaluate it
			return v()
		default:
			// For other types (including custom string types like ApiEnum), use fmt.Sprint
			return fmt.Sprint(v)
		}
	})

	return result
}
