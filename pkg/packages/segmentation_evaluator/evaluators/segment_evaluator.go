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
	"reflect"
	"strconv"

	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/enums"
	segmentUtils "github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/utils"
	"github.com/wingify/vwo-fme-go-sdk/pkg/services"
)

// SegmentEvaluator evaluates segment conditions for users
type SegmentEvaluator struct {
	Context                 *user.VWOContext
	ServiceContainer        interfaces.ServiceContainerInterface // ServiceContainer object
	Feature                 *campaign.Feature
	SegmentOperandEvaluator *SegmentOperandEvaluator
}

// NewSegmentEvaluator creates a new SegmentEvaluator instance
func NewSegmentEvaluator(
	context *user.VWOContext,
	serviceContainer interfaces.ServiceContainerInterface,
	feature *campaign.Feature,
) *SegmentEvaluator {
	return &SegmentEvaluator{
		Context:                 context,
		ServiceContainer:        serviceContainer,
		Feature:                 feature,
		SegmentOperandEvaluator: NewSegmentOperandEvaluator(serviceContainer),
	}
}

// IsSegmentationValid validates if the segmentation defined in the DSL is applicable
func (s *SegmentEvaluator) IsSegmentationValid(dsl map[string]interface{}, properties map[string]interface{}) bool {
	// Get the first key-value pair from the DSL
	operator, subDsl := segmentUtils.GetKeyValue(dsl)

	// Evaluate based on the type of segmentation operator
	operatorEnum, exists := enums.SegmentOperatorValueFromString(operator)
	if !exists {
		return false
	}

	switch operatorEnum {
	case enums.SegmentOperatorNOT:
		if subDslMap, ok := subDsl.(map[string]interface{}); ok {
			result := s.IsSegmentationValid(subDslMap, properties)
			return !result
		}
		return false
	case enums.SegmentOperatorAND:
		if subDslArray, ok := subDsl.([]interface{}); ok {
			return s.Every(subDslArray, properties)
		}
		return false
	case enums.SegmentOperatorOR:
		if subDslArray, ok := subDsl.([]interface{}); ok {
			return s.Some(subDslArray, properties)
		}
		return false
	case enums.SegmentOperatorCustomVariable:
		if subDslMap, ok := subDsl.(map[string]interface{}); ok {
			return s.SegmentOperandEvaluator.EvaluateCustomVariableDSL(subDslMap, properties)
		}
		return false
	case enums.SegmentOperatorUser:
		return s.SegmentOperandEvaluator.EvaluateUserDSL(subDsl.(string), properties)
	case enums.SegmentOperatorUA:
		return s.SegmentOperandEvaluator.EvaluateUserAgentDSL(subDsl.(string), s.Context)
	case enums.SegmentOperatorIP:
		return s.SegmentOperandEvaluator.EvaluateStringOperandDSL(subDsl.(string), s.Context, enums.SegmentOperatorIP)
	case enums.SegmentOperatorBrowserVersion:
		return s.SegmentOperandEvaluator.EvaluateStringOperandDSL(subDsl.(string), s.Context, enums.SegmentOperatorBrowserVersion)
	case enums.SegmentOperatorOSVersion:
		return s.SegmentOperandEvaluator.EvaluateStringOperandDSL(subDsl.(string), s.Context, enums.SegmentOperatorOSVersion)
	default:
		return false
	}
}

// Some evaluates if any of the DSL nodes are valid using the OR logic
func (s *SegmentEvaluator) Some(dslNodes []interface{}, customVariables map[string]interface{}) bool {
	uaParserMap := make(map[string][]string)
	keyCount := 0
	isUAParser := false

	for _, node := range dslNodes {
		dsl, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		for key, value := range dsl {
			// Check for user agent related keys
			keyEnum, exists := enums.SegmentOperatorValueFromString(key)
			if !exists {
				continue
			}

			if keyEnum == enums.SegmentOperatorOperatingSystem ||
				keyEnum == enums.SegmentOperatorBrowserAgent ||
				keyEnum == enums.SegmentOperatorDeviceType ||
				keyEnum == enums.SegmentOperatorDevice {
				isUAParser = true

				if _, ok := uaParserMap[key]; !ok {
					uaParserMap[key] = []string{}
				}

				// Handle value as array or string
				if valueArray, ok := value.([]interface{}); ok {
					for _, val := range valueArray {
						if strVal, ok := val.(string); ok {
							uaParserMap[key] = append(uaParserMap[key], strVal)
						}
					}
				} else if strVal, ok := value.(string); ok {
					uaParserMap[key] = append(uaParserMap[key], strVal)
				}

				keyCount++
			}

			// Check for feature toggle based on feature ID
			if keyEnum == enums.SegmentOperatorFeatureID {
				if featureIDObject, ok := value.(map[string]interface{}); ok {
					for featureIDKey, featureIDValue := range featureIDObject {
						featureIDValueStr, ok := featureIDValue.(string)
						if !ok {
							continue
						}

						if featureIDValueStr == "on" || featureIDValueStr == "off" {
							// Find the feature by ID
							features := s.ServiceContainer.GetSettings().GetFeatures()
							var targetFeature *campaign.Feature
							for _, feature := range features {
								if strconv.Itoa(feature.GetID()) == featureIDKey {
									targetFeature = &feature
									break
								}
							}

							if targetFeature != nil {
								featureKey := targetFeature.Key
								result := s.CheckInUserStorage(featureKey, s.Context)
								if featureIDValueStr == "off" {
									return !result
								}
								return result
							} else {
								s.ServiceContainer.GetLoggerService().Debug("Feature not found with featureIdKey: " + featureIDKey)
								return false
							}
						}
					}
				}
			}
		}

		// Check if the count of keys encountered is equal to dslNodes size
		if isUAParser && keyCount == len(dslNodes) {
			uaParserResult := s.CheckUserAgentParser(uaParserMap)
			return uaParserResult
		}

		// Recursively check each DSL node
		if s.IsSegmentationValid(dsl, customVariables) {
			return true
		}
	}
	return false
}

// Every evaluates all DSL nodes using the AND logic
func (s *SegmentEvaluator) Every(dslNodes []interface{}, customVariables map[string]interface{}) bool {
	locationMap := make(map[string]interface{})

	for _, node := range dslNodes {
		dsl, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		for key := range dsl {
			// Check if the DSL node contains location-related keys
			keyEnum, exists := enums.SegmentOperatorValueFromString(key)
			if !exists {
				continue
			}

			if keyEnum == enums.SegmentOperatorCountry ||
				keyEnum == enums.SegmentOperatorRegion ||
				keyEnum == enums.SegmentOperatorCity {
				s.AddLocationValuesToMap(dsl, locationMap)
				// Check if the number of location keys matches the number of DSL nodes
				if len(locationMap) == len(dslNodes) {
					return s.CheckLocationPreSegmentation(locationMap)
				}
				continue
			}

			res := s.IsSegmentationValid(dsl, customVariables)
			if !res {
				return false
			}
		}
	}
	return true
}

// AddLocationValuesToMap adds location values from a DSL node to a map
func (s *SegmentEvaluator) AddLocationValuesToMap(dsl map[string]interface{}, locationMap map[string]interface{}) {
	// Add country, region, and city information to the location map if present
	for key, value := range dsl {
		keyEnum, exists := enums.SegmentOperatorValueFromString(key)
		if !exists {
			continue
		}

		if keyEnum == enums.SegmentOperatorCountry {
			locationMap[keyEnum.String()] = value
		}
		if keyEnum == enums.SegmentOperatorRegion {
			locationMap[keyEnum.String()] = value
		}
		if keyEnum == enums.SegmentOperatorCity {
			locationMap[keyEnum.String()] = value
		}
	}
}

// CheckLocationPreSegmentation checks if the user's location matches the expected location criteria
func (s *SegmentEvaluator) CheckLocationPreSegmentation(locationMap map[string]interface{}) bool {
	// Ensure user's IP address is available
	if s.Context == nil || s.Context.IPAddress == "" {
		s.ServiceContainer.GetLoggerService().Error("INVALID_IP_ADDRESS_IN_CONTEXT_FOR_PRE_SEGMENTATION", nil, s.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
		return false
	}

	// Check if location data is available and matches the expected values
	if s.Context.VWO == nil || s.Context.VWO.Location == nil || len(s.Context.VWO.Location) == 0 {
		if !s.ServiceContainer.GetSettingsManager().GetIsGatewayServiceProvided() {
			s.ServiceContainer.GetLoggerService().Error("GATEWAY_SERVICE_REQUIRED_FOR_PRE_SEGMENTATION", map[string]interface{}{"preSegmentationType": "location-related"}, s.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
		}
		return false
	}

	return segmentUtils.ValuesMatch(locationMap, s.Context.VWO.Location)
}

// CheckUserAgentParser checks if the user's device information matches the expected criteria
func (s *SegmentEvaluator) CheckUserAgentParser(uaParserMap map[string][]string) bool {
	// Ensure user's user agent is available
	if s.Context == nil || s.Context.UserAgent == "" {
		s.ServiceContainer.GetLoggerService().Error("INVALID_USER_AGENT_IN_CONTEXT_FOR_PRE_SEGMENTATION", nil, s.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
		return false
	}

	// Check if user agent data is available and matches the expected values
	if s.Context.VWO == nil || s.Context.VWO.GetUaInfo() == nil || len(s.Context.VWO.GetUaInfo()) == 0 {
		if !s.ServiceContainer.GetSettingsManager().GetIsGatewayServiceProvided() {
			s.ServiceContainer.GetLoggerService().Error("GATEWAY_SERVICE_REQUIRED_FOR_PRE_SEGMENTATION", map[string]interface{}{"preSegmentationType": "ua-related"}, s.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
		}
		return false
	}

	return segmentUtils.CheckValuePresent(uaParserMap, s.Context.VWO.GetUaInfo())
}

// CheckInUserStorage checks if the feature is enabled for the user by querying the storage (matches TypeScript implementation)
func (s *SegmentEvaluator) CheckInUserStorage(featureKey string, context *user.VWOContext) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			s.ServiceContainer.GetLoggerService().Error("ERROR_CHECKING_FEATURE_IN_USER_STORAGE", map[string]interface{}{"err": fmt.Sprintf("%v", r)}, s.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
			result = false
		}
	}()
	storageService := services.NewStorageService()

	// Retrieve feature data from storage directly (avoiding import cycle with decorators)
	storedData, err := storageService.GetDataInStorage(featureKey, context)
	if err != nil {
		s.ServiceContainer.GetLoggerService().Error("ERROR_READING_DATA_FROM_STORAGE", map[string]interface{}{"err": err.Error()}, s.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
		return false
	}

	// Check if the stored data is an object and not empty (matches TypeScript logic)
	if storedData != nil && reflect.TypeOf(storedData).Kind() == reflect.Map && len(storedData) > 0 {
		return true
	} else {
		return false
	}
}
