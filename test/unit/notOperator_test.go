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

func TestNotOperator(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("ExactMatchTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ExactMatchWithSpecialCharactersTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"f25u!v@b#k$6%9^f&o*v(m)w_-=+s,./` + "`" + `(*&^%$#@!"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "f25u!v@b#k$6%9^f&o*v(m)w_-=+s,./`(*&^%$#@!",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ExactMatchWithSpacesTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"nice to see you. will    you be   my        friend?"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "nice to see you. will    you be   my        friend?",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ExactMatchWithUpperCaseTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"HgUvshFRjsbTnvsdiUFFTGHFHGvDRT.YGHGH"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "HgUvshFRjsbTnvsdiUFFTGHFHGvDRT.YGHGH",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123.456"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123.456,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeExtraDecimalZerosTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123.456"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123.456000000,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123.0,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("StringifiedFloatTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123.456"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "123.456000000",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("StringifiedFloatTest2", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123.0"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("StringifiedFloatTest3", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123.4560000"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123.456,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"true"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          true,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeTest2", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"false"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          false,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MismatchTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "notsomething",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("PartOfTextTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"zzsomethingzz"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "something",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SingleCharTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"zzsomethingzz"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "i",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("CaseMismatchTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "Something",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("CaseMismatchTest2", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "SOMETHING",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NoValueProvidedTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          "",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MissingKeyValueTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NullValueProvidedTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          nil,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("IncorrectKeyTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"neq":         "something",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("IncorrectKeyCaseTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"something"}}]}}`
		customVariables := map[string]interface{}{
			"EQ":          "something",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          12,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123.456"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"123.456"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          123.4567,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"false"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          true,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"not":{"or":[{"custom_variable":{"eq":"true"}}]}}`
		customVariables := map[string]interface{}{
			"eq":          false,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NestedNotOperatorTest", func(t *testing.T) {
		dsl := `{"or":[{"or":[{"not":{"or":[{"or":[{"custom_variable":{"eq":"eq_value"}}]}]}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleNotOperatorTest", func(t *testing.T) {
		dsl := `{"or":[{"not":{"or":[{"not":{"or":[{"custom_variable":{"eq":"eq_value"}}]}}]}}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleNotOperatorTest2", func(t *testing.T) {
		dsl := `{"and":[{"and":[{"not":{"and":[{"and":[{"custom_variable":{"eq":"eq_value"}}]}]}}]}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleNotOperatorTest3", func(t *testing.T) {
		dsl := `{"and":[{"not":{"and":[{"not":{"and":[{"custom_variable":{"eq":"eq_value"}}]}}]}}]}`
		customVariables := map[string]interface{}{
			"eq":          "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleNotOperatorTest4", func(t *testing.T) {
		dsl := `{"not":{"or":[{"not":{"or":[{"not":{"or":[{"not":{"or":[{"custom_variable":{"neq":"eq_value"}}]}}]}}]}}]}}`
		customVariables := map[string]interface{}{
			"neq":         "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleNotOperatorTest5", func(t *testing.T) {
		dsl := `{"not":{"or":[{"not":{"or":[{"not":{"or":[{"not":{"or":[{"custom_variable":{"neq":"not_eq_value"}}]}}]}}]}}]}}`
		customVariables := map[string]interface{}{
			"neq":         "eq_value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleNotOperatorTest6", func(t *testing.T) {
		dsl := `{"not":{"or":[{"not":{"or":[{"not":{"or":[{"not":{"or":[{"not":{"or":[{"custom_variable":{"neq":"eq_value"}}]}}]}}]}}]}}]}}`
		customVariables := map[string]interface{}{
			"neq":         "eq_value",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MultipleNotOperatorTest7", func(t *testing.T) {
		dsl := `{"not":{"or":[{"not":{"or":[{"not":{"or":[{"not":{"or":[{"not":{"or":[{"custom_variable":{"neq":"neq_value"}}]}}]}}]}}]}}]}}`
		customVariables := map[string]interface{}{
			"neq":         "eq_value",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})
}
