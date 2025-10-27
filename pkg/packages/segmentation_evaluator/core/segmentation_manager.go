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

package core

import (
	"encoding/json"
	"fmt"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/gateway"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/evaluators"
)

// SegmentationManager manages segmentation evaluation
// Ensure SegmentationManager implements SegmentationManagerInterface
var _ interfaces.SegmentationManagerInterface = (*SegmentationManager)(nil)

type SegmentationManager struct {
	evaluator  *evaluators.SegmentEvaluator
	logManager interfaces.LoggerServiceInterface
}

// NewSegmentationManager creates a new SegmentationManager instance
func NewSegmentationManager(logManager interfaces.LoggerServiceInterface) *SegmentationManager {
	return &SegmentationManager{
		logManager: logManager,
	}
}

// NewSegmentationManagerWithEvaluator creates a new SegmentationManager instance with evaluator initialization
func NewSegmentationManagerWithEvaluator(logManager interfaces.LoggerServiceInterface, shouldInitializeEvaluator bool) *SegmentationManager {
	sm := &SegmentationManager{
		logManager: logManager,
	}
	if shouldInitializeEvaluator {
		sm.evaluator = evaluators.NewSegmentEvaluator(nil, nil, nil)
	}
	return sm
}

// SetContextualData sets the contextual data required for segmentation
func (sm *SegmentationManager) SetContextualData(
	serviceContainer interfaces.ServiceContainerInterface, // ServiceContainer object containing the settings manager
	feature *campaign.Feature, // FeatureModel object containing the feature settings
	context *user.VWOContext, // VWOContext object containing the user context
) {

	defer func() {
		if r := recover(); r != nil {
			sm.logManager.Error("ERROR_SETTING_SEGMENTATION_CONTEXT", map[string]interface{}{"err": fmt.Sprintf("%v", r)}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	// Create a new evaluator with the contextual data
	sm.evaluator = evaluators.NewSegmentEvaluator(
		context,
		serviceContainer,
		feature,
	)

	// Set the serviceContainer in the evaluator
	sm.evaluator.ServiceContainer = serviceContainer

	// If user agent and ipAddress both are null or empty, return
	if context.UserAgent == "" && context.IPAddress == "" {
		return
	}

	// If gateway service is required and the base URL is not the default one, fetch the data from the gateway service
	baseURL := serviceContainer.GetBaseUrl()
	if feature.IsGatewayServiceRequired && baseURL != constants.HostName && context.VWO == nil {
		queryParams := make(map[string]string)

		if context.UserAgent == "" && context.IPAddress == "" {
			return
		}

		if context.UserAgent != "" {
			queryParams[constants.QueryParamUserAgent] = context.UserAgent
		}

		if context.IPAddress != "" {
			queryParams[constants.QueryParamIPAddress] = context.IPAddress
		}

		vwoData, err := gateway.GetFromGatewayService(serviceContainer, queryParams, constants.EndpointGetUserData)
		if err != nil {
			return
		}

		if vwoData != "" {
			// Parse the gateway service response
			var gatewayServiceModel user.ContextVWO
			err := json.Unmarshal([]byte(vwoData), &gatewayServiceModel)
			if err != nil {
				sm.logManager.Error("ERROR_SETTING_SEGMENTATION_CONTEXT", map[string]interface{}{"err": err.Error()}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
				return
			}
			context.VWO = &gatewayServiceModel
		}
	}
}

// ValidateSegmentation validates the segmentation for the given DSL and properties
func (sm *SegmentationManager) ValidateSegmentation(dsl interface{}, properties map[string]interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			sm.logManager.Error("ERROR_VALIDATING_SEGMENTATION", map[string]interface{}{"err": fmt.Sprintf("%v", r)}, sm.evaluator.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	var dslMap map[string]interface{}

	// Handle DSL as string or object
	switch v := dsl.(type) {
	case string:
		err := json.Unmarshal([]byte(v), &dslMap)
		if err != nil {
			sm.logManager.Error("ERROR_VALIDATING_SEGMENTATION", map[string]interface{}{"err": err.Error()}, sm.evaluator.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
			return false
		}
	case map[string]interface{}:
		dslMap = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			sm.logManager.Error("ERROR_VALIDATING_SEGMENTATION", map[string]interface{}{"err": err.Error()}, sm.evaluator.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
			return false
		}
		err = json.Unmarshal(jsonBytes, &dslMap)
		if err != nil {
			sm.logManager.Error("ERROR_VALIDATING_SEGMENTATION", map[string]interface{}{"err": err.Error()}, sm.evaluator.ServiceContainer.GetDebuggerService().GetStandardDebugProps())
			return false
		}
	}

	return sm.evaluator.IsSegmentationValid(dslMap, properties)
}
