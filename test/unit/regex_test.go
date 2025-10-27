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

func TestRegex(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("RegexOperandTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}`
		customVariables := map[string]interface{}{
			"reg":         "myregexxxxxx",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("RegexOperandTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"reg":"regex(<(W[^>]*)(.*?)>)"}}]}`
		customVariables := map[string]interface{}{
			"reg":         "<WingifySDK id=1></WingifySDK>",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("RegexOperandMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"reg":"regex(<(W[^>]*)(.*?)>)"}}]}`
		customVariables := map[string]interface{}{
			"reg":         "<wingifySDK id=1></wingifySDK>",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("RegexOperandCaseMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}`
		customVariables := map[string]interface{}{
			"reg":         "myregeXxxxxx",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("InvalidRegexTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"reg":"regex(*)"}}]}`
		customVariables := map[string]interface{}{
			"reg":         "*",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("InvalidRegexTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"reg":"regex(*)"}}]}`
		customVariables := map[string]interface{}{
			"reg":         "asdf",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})
}
