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
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// EvaluateRuleResult represents the result of rule evaluation
type EvaluateRuleResult struct {
	PreSegmentationResult bool
	WhitelistedObject     *campaign.Variation
	UpdatedDecision       map[string]interface{}
}

// EvaluateRule evaluates the rule for a given feature and campaign and returns the rule evaluation result
func EvaluateRule(
	serviceContainer interfaces.ServiceContainerInterface, // ServiceContainer object
	feature *campaign.Feature,
	campaignModel *campaign.Campaign,
	context *user.VWOContext,
	evaluatedFeatureMap map[string]interface{},
	megGroupWinnerCampaigns map[int]string,
	storageService interfaces.StorageServiceInterface,
	decision map[string]interface{},
) (ruleEvaluationResult map[string]interface{}) {
	// Perform whitelisting and pre-segmentation checks
	defer func() {
		if r := recover(); r != nil {
			// Log error using serviceContainer logger
			if serviceContainer != nil {
				serviceContainer.GetLoggerService().Error("ERROR_EVALUATING_RULE", map[string]interface{}{
					"err": r,
				}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
			}
			ruleEvaluationResult = map[string]interface{}{
				enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): false,
				enums.EvaluatedRuleResultWhitelistedObject.GetValue():     nil,
				enums.EvaluatedRuleResultUpdatedDecision.GetValue():       decision,
			}
		}
	}()

	// Check if the campaign satisfies the whitelisting and pre-segmentation
	checkResult := CheckWhitelistingAndPreSeg(
		serviceContainer,
		feature,
		campaignModel,
		context,
		evaluatedFeatureMap,
		megGroupWinnerCampaigns,
		storageService,
		decision,
	)

	// Extract the results of the evaluation
	preSegmentationResult := checkResult[enums.EvaluatedRuleResultPreSegmentationResult.GetValue()].(bool)
	var whitelistedObject *campaign.Variation = nil
	if checkResult[enums.EvaluatedRuleResultWhitelistedObject.GetValue()] != nil {
		if variation, ok := checkResult[enums.EvaluatedRuleResultWhitelistedObject.GetValue()].(*campaign.Variation); ok {
			whitelistedObject = variation
		}
	}

	// If pre-segmentation is successful and a whitelisted object exists, proceed to send an impression
	if preSegmentationResult && whitelistedObject != nil && whitelistedObject.ID != 0 {
		// Update the decision object with campaign and variation details
		decision[enums.DecisionExperimentID.GetValue()] = campaignModel.ID
		decision[enums.DecisionExperimentKey.GetValue()] = campaignModel.Key
		decision[enums.DecisionExperimentVariationID.GetValue()] = whitelistedObject.ID
	}

	// Return the results of the evaluation
	ruleEvaluationResult = map[string]interface{}{
		enums.EvaluatedRuleResultPreSegmentationResult.GetValue(): preSegmentationResult,
		enums.EvaluatedRuleResultWhitelistedObject.GetValue():     whitelistedObject,
		enums.EvaluatedRuleResultUpdatedDecision.GetValue():       decision,
	}

	return ruleEvaluationResult
}
