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

func TestAndOperator(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("SingleAndOperatorMatchingTest", func(t *testing.T) {
		dsl := `{"and":[{"custom_variable":{"eq":"eq_value"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SingleAndOperatorMismatchTest", func(t *testing.T) {
		dsl := `{"and":[{"custom_variable":{"eq":"eq_value"}}]}`
		customVariables := map[string]interface{}{
			"a":           "n_eq_value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SingleAndOperatorCaseMismatchTest", func(t *testing.T) {
		dsl := `{"and":[{"custom_variable":{"eq":"eq_value"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "Eq_Value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleAndOperatorTest2", func(t *testing.T) {
		dsl := `{"and":[{"and":[{"and":[{"and":[{"and":[{"custom_variable":{"eq":"eq_value"}}]}]}]}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleAndOperatorWithSingleCorrectValueTest", func(t *testing.T) {
		dsl := `{"and":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"reg":         "wrong",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleAndOperatorWithSingleCorrectValueTest2", func(t *testing.T) {
		dsl := `{"and":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "wrong",
			"reg":         "myregexxxxxx",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleAndOperatorWithAllCorrectCustomVariablesTest", func(t *testing.T) {
		dsl := `{"and":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"reg":         "myregexxxxxx",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleAndOperatorWithAllIncorrectCorrectCustomVariablesTest", func(t *testing.T) {
		dsl := `{"and":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "wrong",
			"reg":         "wrong",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})
}
