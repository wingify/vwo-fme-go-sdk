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

func TestOrOperator(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("SingleOrOperatorMatchingTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"eq_value"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SingleOrOperatorMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"eq_value"}}]}`
		customVariables := map[string]interface{}{
			"a":           "n_eq_value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SingleOrOperatorCaseMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"eq_value"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "Eq_Value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleOrOperatorTest", func(t *testing.T) {
		dsl := `{"or":[{"or":[{"or":[{"or":[{"or":[{"custom_variable":{"eq":"eq_value"}}]}]}]}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleOrOperatorWithSingleCorrectValueTest", func(t *testing.T) {
		dsl := `{"or":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"reg":         "wrong",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleOrOperatorWithSingleCorrectValueTest2", func(t *testing.T) {
		dsl := `{"or":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "wrong",
			"reg":         "myregexxxxxx",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleOrOperatorWithAllCorrectCustomVariablesTest", func(t *testing.T) {
		dsl := `{"or":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"reg":         "myregeXxxxxx",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleOrOperatorWithAllIncorrectCorrectCustomVariablesTest", func(t *testing.T) {
		dsl := `{"or":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "wrong",
			"reg":         "wrong",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})
}
