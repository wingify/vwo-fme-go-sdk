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

package unit

import (
	"testing"

	loggerCore "github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/core"
	segmentationCore "github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/core"
)

func TestGreaterThanOperator(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("GreaterThanOperatorPass", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"gt(150)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          160,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("GreaterThanOperatorFail", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"gt(150)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          100,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("GreaterThanOperatorEqualValueFail", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"gt(150)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          150,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("GreaterThanOperatorStringValueFail", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"gt(150)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "abc",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})
}
