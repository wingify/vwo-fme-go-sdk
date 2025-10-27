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

package evaluators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/gateway"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/enums"
	segmentUtils "github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/utils"
)

// SegmentOperandEvaluator evaluates operand conditions in segmentation
type SegmentOperandEvaluator struct {
	serviceContainer interfaces.ServiceContainerInterface // ServiceContainer object
}

// NewSegmentOperandEvaluator creates a new SegmentOperandEvaluator instance
func NewSegmentOperandEvaluator(serviceContainer interfaces.ServiceContainerInterface) *SegmentOperandEvaluator {
	return &SegmentOperandEvaluator{
		serviceContainer: serviceContainer,
	}
}

// PreProcessOperandValueResult holds the result of pre-processing an operand value
type PreProcessOperandValueResult struct {
	OperandType  enums.SegmentOperandValue
	OperandValue string
}

// EvaluateCustomVariableDSL evaluates custom variable DSL
func (s *SegmentOperandEvaluator) EvaluateCustomVariableDSL(dslOperandValue map[string]interface{}, properties map[string]interface{}) bool {
	// Get the first key-value pair from dslOperandValue
	operandKey, operandValueNode := segmentUtils.GetKeyValue(dslOperandValue)
	operandValue := fmt.Sprint(operandValueNode)

	// Check if the property exists
	if _, exists := properties[operandKey]; !exists {
		return false
	}

	// Handle 'inlist' operand
	if strings.Contains(operandValue, "inlist") {
		listIDPattern := regexp.MustCompile(`inlist\(([^)]+)\)`)
		matches := listIDPattern.FindStringSubmatch(operandValue)
		if len(matches) < 2 {
			s.serviceContainer.GetLoggerService().Error("INVALID_ATTRIBUTE_LIST_FORMAT", nil, s.serviceContainer.GetDebuggerService().GetStandardDebugProps())
			return false
		}
		listID := matches[1]

		// Process the tag value and prepare query parameters
		tagValue := properties[operandKey]
		attributeValue := s.PreProcessTagValue(fmt.Sprint(tagValue))
		queryParams := map[string]string{
			"attribute": attributeValue,
			"listId":    listID,
		}

		response, err := gateway.GetFromGatewayService(s.serviceContainer, queryParams, constants.EndpointAttributeCheck)
		if err != nil {
			return false
		}

		if response == "true" {
			return true
		} else {
			return false
		}
	} else {
		// Process other types of operands
		tagValue := properties[operandKey]
		if tagValue == nil {
			tagValue = ""
		}
		// Convert tagValue to string properly to avoid scientific notation
		var tagValueStr string
		if f, ok := tagValue.(float64); ok {
			if f == float64(int(f)) {
				tagValueStr = strconv.Itoa(int(f))
			} else {
				tagValueStr = strconv.FormatFloat(f, 'f', -1, 64)
			}
		} else {
			tagValueStr = fmt.Sprint(tagValue)
		}
		tagValue = s.PreProcessTagValue(tagValueStr)
		preProcessOperandValue := s.PreProcessOperandValue(operandValue)
		processedValues := s.ProcessValues(preProcessOperandValue.OperandValue, tagValue)

		// Convert numeric values to strings if processing wildcard pattern
		operandType := preProcessOperandValue.OperandType
		if operandType == enums.SegmentOperandStartingEndingStarValue ||
			operandType == enums.SegmentOperandStartingStarValue ||
			operandType == enums.SegmentOperandEndingStarValue ||
			operandType == enums.SegmentOperandRegexValue {
		}

		tagValue = processedValues["tagValue"]
		return s.ExtractResult(operandType, strings.TrimSpace(strings.ReplaceAll(fmt.Sprint(processedValues["operandValue"]), "\"", "")), fmt.Sprint(tagValue))
	}
}

// PreProcessOperandValue pre-processes the operand value to extract operand type and value
func (s *SegmentOperandEvaluator) PreProcessOperandValue(operand string) *PreProcessOperandValueResult {
	var operandType enums.SegmentOperandValue
	var operandValue string

	if segmentUtils.MatchWithRegex(operand, enums.SegmentOperandRegexLowerMatch.GetRegex()) {
		operandType = enums.SegmentOperandLowerValue
		operandValue = s.ExtractOperandValue(operand, enums.SegmentOperandRegexLowerMatch.GetRegex())
	} else if segmentUtils.MatchWithRegex(operand, enums.SegmentOperandRegexWildcardMatch.GetRegex()) {
		operandValue = s.ExtractOperandValue(operand, enums.SegmentOperandRegexWildcardMatch.GetRegex())
		startingStar := segmentUtils.MatchWithRegex(operandValue, enums.SegmentOperandRegexStartingStar.GetRegex())
		endingStar := segmentUtils.MatchWithRegex(operandValue, enums.SegmentOperandRegexEndingStar.GetRegex())
		if startingStar && endingStar {
			operandType = enums.SegmentOperandStartingEndingStarValue
		} else if startingStar {
			operandType = enums.SegmentOperandStartingStarValue
		} else if endingStar {
			operandType = enums.SegmentOperandEndingStarValue
		} else {
			operandType = enums.SegmentOperandRegexValue
		}
		operandValue = regexp.MustCompile(enums.SegmentOperandRegexStartingStar.GetRegex()).ReplaceAllString(operandValue, "")
		operandValue = regexp.MustCompile(enums.SegmentOperandRegexEndingStar.GetRegex()).ReplaceAllString(operandValue, "")
	} else if segmentUtils.MatchWithRegex(operand, enums.SegmentOperandRegexRegexMatch.GetRegex()) {
		operandType = enums.SegmentOperandRegexValue
		operandValue = s.ExtractOperandValue(operand, enums.SegmentOperandRegexRegexMatch.GetRegex())
	} else if segmentUtils.MatchWithRegex(operand, enums.SegmentOperandRegexGreaterThanMatch.GetRegex()) {
		operandType = enums.SegmentOperandGreaterThanValue
		operandValue = s.ExtractOperandValue(operand, enums.SegmentOperandRegexGreaterThanMatch.GetRegex())
	} else if segmentUtils.MatchWithRegex(operand, enums.SegmentOperandRegexGreaterThanEqualToMatch.GetRegex()) {
		operandType = enums.SegmentOperandGreaterThanEqualToValue
		operandValue = s.ExtractOperandValue(operand, enums.SegmentOperandRegexGreaterThanEqualToMatch.GetRegex())
	} else if segmentUtils.MatchWithRegex(operand, enums.SegmentOperandRegexLessThanMatch.GetRegex()) {
		operandType = enums.SegmentOperandLessThanValue
		operandValue = s.ExtractOperandValue(operand, enums.SegmentOperandRegexLessThanMatch.GetRegex())
	} else if segmentUtils.MatchWithRegex(operand, enums.SegmentOperandRegexLessThanEqualToMatch.GetRegex()) {
		operandType = enums.SegmentOperandLessThanEqualToValue
		operandValue = s.ExtractOperandValue(operand, enums.SegmentOperandRegexLessThanEqualToMatch.GetRegex())
	} else {
		operandType = enums.SegmentOperandEqualValue
		operandValue = operand
	}

	return &PreProcessOperandValueResult{
		OperandType:  operandType,
		OperandValue: operandValue,
	}
}

// EvaluateUserDSL evaluates user DSL
func (s *SegmentOperandEvaluator) EvaluateUserDSL(dslOperandValue string, properties map[string]interface{}) bool {
	users := strings.Split(dslOperandValue, ",")
	for _, user := range users {
		userTrimmed := strings.TrimSpace(strings.ReplaceAll(user, "\"", ""))
		if vwoUserID, exists := properties["_vwoUserId"]; exists {
			if userTrimmed == fmt.Sprint(vwoUserID) {
				return true
			}
		}
	}
	return false
}

// EvaluateUserAgentDSL evaluates user agent DSL
func (s *SegmentOperandEvaluator) EvaluateUserAgentDSL(dslOperandValue string, context *user.VWOContext) bool {
	if context == nil || context.UserAgent == "" {
		s.serviceContainer.GetLoggerService().Error("INVALID_USER_AGENT_IN_CONTEXT_FOR_PRE_SEGMENTATION", nil, s.serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return false
	}
	tagValue := context.UserAgent // Already URL decoded
	preProcessOperandValue := s.PreProcessOperandValue(dslOperandValue)
	processedValues := s.ProcessValues(preProcessOperandValue.OperandValue, tagValue)
	tagValue = fmt.Sprint(processedValues["tagValue"])
	operandType := preProcessOperandValue.OperandType
	return s.ExtractResult(operandType, strings.TrimSpace(strings.ReplaceAll(fmt.Sprint(processedValues["operandValue"]), "\"", "")), tagValue)
}

// EvaluateStringOperandDSL evaluates a given string tag value against a DSL operand value
// Supported operand types: ip_address, browser_version, os_version
func (s *SegmentOperandEvaluator) EvaluateStringOperandDSL(dslOperandValue string, context *user.VWOContext, operandType enums.SegmentOperatorValue) bool {
	operand := dslOperandValue

	// Determine the tag value based on operand type
	tagValue := s.getTagValueForOperandType(context, operandType)
	if tagValue == "" {
		s.logMissingContextError(operandType)
		return false
	}

	preProcessOperandValue := s.PreProcessOperandValue(operand)
	processedValues := s.ProcessValuesForOperandType(preProcessOperandValue.OperandValue, tagValue, operandType)
	processedTagValue := fmt.Sprint(processedValues["tagValue"])

	return s.ExtractResult(
		preProcessOperandValue.OperandType,
		strings.TrimSpace(strings.ReplaceAll(fmt.Sprint(processedValues["operandValue"]), "\"", "")),
		processedTagValue,
	)
}

// getTagValueForOperandType returns the appropriate tag value based on the operand type
func (s *SegmentOperandEvaluator) getTagValueForOperandType(context *user.VWOContext, operandType enums.SegmentOperatorValue) string {
	if operandType == enums.SegmentOperatorIP {
		if context != nil {
			return context.GetIPAddress()
		}
		return ""
	} else if operandType == enums.SegmentOperatorBrowserVersion {
		return s.getBrowserVersionFromContext(context)
	}
	// Default for OS version
	return s.getOsVersionFromContext(context)
}

// getBrowserVersionFromContext extracts browser version from VWO context
func (s *SegmentOperandEvaluator) getBrowserVersionFromContext(context *user.VWOContext) string {
	if context == nil || context.GetVWO() == nil {
		return ""
	}
	ua := context.GetVWO().GetUaInfo()
	if len(ua) == 0 {
		return ""
	}
	if v, ok := ua[enums.SegmentOperatorBrowserVersion.String()]; ok && v != "" {
		return v
	}
	return ""
}

// getOsVersionFromContext extracts OS version from VWO context
func (s *SegmentOperandEvaluator) getOsVersionFromContext(context *user.VWOContext) string {
	if context == nil || context.GetVWO() == nil {
		return ""
	}
	ua := context.GetVWO().GetUaInfo()
	if len(ua) == 0 {
		return ""
	}
	if v, ok := ua[enums.SegmentOperatorOSVersion.String()]; ok && v != "" {
		return v
	}
	return ""
}

// logMissingContextError logs appropriate error message for missing context
func (s *SegmentOperandEvaluator) logMissingContextError(operandType enums.SegmentOperatorValue) {
	if operandType == enums.SegmentOperatorIP {
		s.serviceContainer.GetLoggerService().Error("INVALID_IP_ADDRESS_IN_CONTEXT_FOR_PRE_SEGMENTATION", nil, s.serviceContainer.GetDebuggerService().GetStandardDebugProps())
	} else if operandType == enums.SegmentOperatorBrowserVersion {
		s.serviceContainer.GetLoggerService().Error("INVALID_USER_AGENT_IN_CONTEXT_FOR_PRE_SEGMENTATION", nil, s.serviceContainer.GetDebuggerService().GetStandardDebugProps())
	} else {
		s.serviceContainer.GetLoggerService().Error("INVALID_USER_AGENT_IN_CONTEXT_FOR_PRE_SEGMENTATION", nil, s.serviceContainer.GetDebuggerService().GetStandardDebugProps())
	}
}

// PreProcessTagValue pre-processes the tag value
func (s *SegmentOperandEvaluator) PreProcessTagValue(tagValue string) string {
	if tagValue == "" {
		return ""
	}
	// Simple boolean check
	if tagValue == "true" || tagValue == "false" {
		boolVal, _ := strconv.ParseBool(tagValue)
		return strconv.FormatBool(boolVal)
	}
	return strings.TrimSpace(tagValue)
}

// ProcessValues processes operand and tag values
func (s *SegmentOperandEvaluator) ProcessValues(operandValue interface{}, tagValue interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Convert values to strings first, but handle floats properly to avoid scientific notation
	var operandValueStr, tagValueStr string

	// Handle operand value
	if f, ok := operandValue.(float64); ok {
		if f == float64(int(f)) {
			operandValueStr = strconv.Itoa(int(f))
		} else {
			operandValueStr = strconv.FormatFloat(f, 'f', -1, 64)
		}
	} else {
		operandValueStr = fmt.Sprint(operandValue)
	}

	// Handle tag value
	if f, ok := tagValue.(float64); ok {
		if f == float64(int(f)) {
			tagValueStr = strconv.Itoa(int(f))
		} else {
			tagValueStr = strconv.FormatFloat(f, 'f', -1, 64)
		}
	} else {
		tagValueStr = fmt.Sprint(tagValue)
	}

	// Check if tag value contains non-numeric characters (except decimal point)
	// This matches the Node.js NON_NUMERIC_PATTERN = /[^0-9.]/
	hasNonNumeric := false
	for _, char := range tagValueStr {
		if char < '0' || char > '9' {
			if char != '.' {
				hasNonNumeric = true
				break
			}
		}
	}

	// If tag value has non-numeric characters, return original values as strings
	if hasNonNumeric {
		result["operandValue"] = operandValueStr
		result["tagValue"] = tagValueStr
		return result
	}

	// Try to convert both values to numbers
	operandFloat, operandErr := strconv.ParseFloat(operandValueStr, 64)
	tagFloat, tagErr := strconv.ParseFloat(tagValueStr, 64)

	// If either conversion fails, return original values as strings
	if operandErr != nil || tagErr != nil {
		result["operandValue"] = operandValueStr
		result["tagValue"] = tagValueStr
		return result
	}

	// Both are valid numbers, convert them back to strings
	// Check if the numeric value is actually an integer
	if operandFloat == float64(int(operandFloat)) {
		result["operandValue"] = strconv.Itoa(int(operandFloat))
	} else {
		result["operandValue"] = strconv.FormatFloat(operandFloat, 'f', -1, 64)
	}

	if tagFloat == float64(int(tagFloat)) {
		result["tagValue"] = strconv.Itoa(int(tagFloat))
	} else {
		result["tagValue"] = strconv.FormatFloat(tagFloat, 'f', -1, 64)
	}

	return result
}

// ProcessValuesForOperandType processes values but skips numeric handling for specified operand types
func (s *SegmentOperandEvaluator) ProcessValuesForOperandType(operandValue interface{}, tagValue interface{}, operandType enums.SegmentOperatorValue) map[string]interface{} {
	if operandType == enums.SegmentOperatorIP || operandType == enums.SegmentOperatorBrowserVersion || operandType == enums.SegmentOperatorOSVersion {
		return map[string]interface{}{
			"operandValue": fmt.Sprint(operandValue),
			"tagValue":     fmt.Sprint(tagValue),
		}
	}
	return s.ProcessValues(operandValue, tagValue)
}

// ConvertValue converts a value to appropriate format
func (s *SegmentOperandEvaluator) ConvertValue(value interface{}) interface{} {
	// Check if the value is a boolean
	if b, ok := value.(bool); ok {
		return strconv.FormatBool(b)
	}

	// Try to convert to numeric value
	valueStr := fmt.Sprint(value)
	if numericValue, err := strconv.ParseFloat(valueStr, 64); err == nil {
		// Check if the numeric value is actually an integer
		if numericValue == float64(int(numericValue)) {
			return strconv.Itoa(int(numericValue))
		}
		// Format float to avoid scientific notation
		return strconv.FormatFloat(numericValue, 'f', -1, 64)
	}

	// Return the value as-is if it's not a number
	return valueStr
}

// ExtractResult extracts the result of the evaluation based on the operand type and values
func (s *SegmentOperandEvaluator) ExtractResult(operandType enums.SegmentOperandValue, operandValue interface{}, tagValue string) bool {
	result := false
	operandValueStr := fmt.Sprint(operandValue)

	switch operandType {
	case enums.SegmentOperandLowerValue:
		result = strings.EqualFold(operandValueStr, tagValue)
	case enums.SegmentOperandStartingEndingStarValue:
		result = strings.Contains(tagValue, operandValueStr)
	case enums.SegmentOperandStartingStarValue:
		result = strings.HasSuffix(tagValue, operandValueStr)
	case enums.SegmentOperandEndingStarValue:
		result = strings.HasPrefix(tagValue, operandValueStr)
	case enums.SegmentOperandRegexValue:
		regex, err := regexp.Compile(operandValueStr)
		if err != nil {
			result = false
		} else {
			result = regex.MatchString(tagValue)
		}
	case enums.SegmentOperandGreaterThanValue:
		if s.isVersionString(tagValue) && s.isVersionString(operandValueStr) {
			result = s.compareVersions(tagValue, operandValueStr) > 0
		} else {
			operandFloat, operandErr := strconv.ParseFloat(operandValueStr, 64)
			tagFloat, tagErr := strconv.ParseFloat(tagValue, 64)
			if operandErr != nil || tagErr != nil {
				result = false
			} else {
				result = tagFloat > operandFloat
			}
		}
	case enums.SegmentOperandGreaterThanEqualToValue:
		if s.isVersionString(tagValue) && s.isVersionString(operandValueStr) {
			result = s.compareVersions(tagValue, operandValueStr) >= 0
		} else {
			operandFloat, operandErr := strconv.ParseFloat(operandValueStr, 64)
			tagFloat, tagErr := strconv.ParseFloat(tagValue, 64)
			if operandErr != nil || tagErr != nil {
				result = false
			} else {
				result = tagFloat >= operandFloat
			}
		}
	case enums.SegmentOperandLessThanValue:
		if s.isVersionString(tagValue) && s.isVersionString(operandValueStr) {
			result = s.compareVersions(tagValue, operandValueStr) < 0
		} else {
			operandFloat, operandErr := strconv.ParseFloat(operandValueStr, 64)
			tagFloat, tagErr := strconv.ParseFloat(tagValue, 64)
			if operandErr != nil || tagErr != nil {
				result = false
			} else {
				result = tagFloat < operandFloat
			}
		}
	case enums.SegmentOperandLessThanEqualToValue:
		if s.isVersionString(tagValue) && s.isVersionString(operandValueStr) {
			result = s.compareVersions(tagValue, operandValueStr) <= 0
		} else {
			operandFloat, operandErr := strconv.ParseFloat(operandValueStr, 64)
			tagFloat, tagErr := strconv.ParseFloat(tagValue, 64)
			if operandErr != nil || tagErr != nil {
				result = false
			} else {
				result = tagFloat <= operandFloat
			}
		}
	default:
		// If both look like version strings, use version comparison
		if s.isVersionString(tagValue) && s.isVersionString(operandValueStr) {
			result = s.compareVersions(tagValue, operandValueStr) == 0
		} else {
			result = tagValue == operandValueStr
		}
	}

	return result
}

// ExtractOperandValue extracts the operand value based on the provided regex pattern
func (s *SegmentOperandEvaluator) ExtractOperandValue(operand string, regexPattern string) string {
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(operand)
	if len(matches) >= 2 {
		return matches[1]
	}
	return operand
}

// isVersionString checks if a string appears to be a version string (digits and dots)
func (s *SegmentOperandEvaluator) isVersionString(str string) bool {
	if str == "" {
		return false
	}
	// ^(\d+\.)*\d+$
	pattern := regexp.MustCompile(`^(\d+\.)*\d+$`)
	return pattern.MatchString(str)
}

// compareVersions compares two semantic-like version strings
// Returns -1 if v1 < v2, 0 if equal, 1 if v1 > v2
func (s *SegmentOperandEvaluator) compareVersions(v1 string, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		p1 := 0
		p2 := 0
		if i < len(parts1) {
			if m, err := strconv.Atoi(matchDigits(parts1[i])); err == nil {
				p1 = m
			}
		}
		if i < len(parts2) {
			if m, err := strconv.Atoi(matchDigits(parts2[i])); err == nil {
				p2 = m
			}
		}
		if p1 < p2 {
			return -1
		} else if p1 > p2 {
			return 1
		}
	}
	return 0
}

// matchDigits extracts leading digits from a string part, defaults to 0 when no digits
func matchDigits(s string) string {
	if s == "" {
		return "0"
	}
	re := regexp.MustCompile(`^\d+`)
	if re.MatchString(s) {
		return re.FindString(s)
	}
	return "0"
}
