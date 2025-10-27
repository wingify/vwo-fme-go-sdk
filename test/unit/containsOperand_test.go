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

func TestContainsOperand(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("ExactMatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "something",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SuffixMatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "asdn3kn42knsdsomething",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("PrefixMatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "somethingdfgdwerewew",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ContainsMatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "asdn3kn42knsdsomethingjsbdj",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SpecialCharactersTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*f25u!v@b#k$6%9^f&o*v(m)w_-=+s,./` + "`" + `(*&^%$#@!*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "A-N-Y-T-H-I-N-G---f25u!v@b#k$6%9^f&o*v(m)w_-=+s,./`(*&^%$#@!----A-N-Y-T-H-I-N-G",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ContainsOperandFalsyTestWithSpecialCharacter", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"a":"wildcard(*some*thing*)"}}]}`
		customVariables := map[string]interface{}{
			"a":           "hellosomethingworld",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SpacesTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*nice to see you. will    you be   my        friend?*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "Hello there!! nice to see you. will    you be   my        friend? Yes, Great!!",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("UpperCaseTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*HgUvshFRjsbTnvsdiUFFTGHFHGvDRT.YGHGH*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "A-N-Y-T-H-I-N-G---HgUvshFRjsbTnvsdiUFFTGHFHGvDRT.YGHGH----A-N-Y-T-H-I-N-G",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*123*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          365412363,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*123.456*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          765123.4567364,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*123*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          765123.7364,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("StringifiedFloatTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*123.456*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "87654123.4567902",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*true*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          true,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*false*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          false,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "qwertyu",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("PartOfTextTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*zzsomethingzz*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SingleCharTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*zzsomethingzz*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "i",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("CaseMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "Something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("CaseMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "SOMETHING",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NoValueProvidedTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MissingKeyValueTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NullValueProvidedTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          nil,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("IncorrectKeyTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"neq":         "something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("IncorrectKeyCaseTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*something*)"}}]}`
		customVariables := map[string]interface{}{
			"EQ":          "something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*123*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          12,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*false*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          true,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"wildcard(*true*)"}}]}`
		customVariables := map[string]interface{}{
			"eq":          false,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})
}
