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

	"github.com/stretchr/testify/assert"
	loggerCore "github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/core"
	segmentationCore "github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/core"
)

func TestEqualityOperand(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("ExactMatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "something",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ExactMatchWithSpecialCharactersTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"f25u!v@b#k$6%9^f&o*v(m)w_-=+s,./` + "`" + `(*&^%$#@!"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "f25u!v@b#k$6%9^f&o*v(m)w_-=+s,./`(*&^%$#@!",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ExactMatchWithSpacesTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"nice to see you. will    you be   my        friend?"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "nice to see you. will    you be   my        friend?",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ExactMatchWithSpacesTest1", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"   nice to see you. will    you be   my        friend?   "}}]}`
		customVariables := map[string]interface{}{
			"eq":          "nice to see you. will you be my friend?",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("ExactMatchWithUpperCaseTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"HgUvshFRjsbTnvsdiUFFTGHFHGvDRT.YGHGH"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "HgUvshFRjsbTnvsdiUFFTGHFHGvDRT.YGHGH",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("TrimValueTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"  hi  "}}]}`
		customVariables := map[string]interface{}{
			"eq":          "          hi          ",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123.456"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123.456,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeExtraDecimalZerosTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123.456"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123.456000000,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123.0,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("StringifiedFloatTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123.456"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "123.456000000",
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("StringifiedFloatTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123.0"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("StringifiedFloatTest3", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123.4560000"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123.456,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"true"}}]}`
		customVariables := map[string]interface{}{
			"eq":          true,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"false"}}]}`
		customVariables := map[string]interface{}{
			"eq":          false,
			"expectation": true,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "notsomething",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("PartOfTextTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"zzsomethingzz"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("SingleCharTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"zzsomethingzz"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "i",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("CaseMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "Something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("CaseMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "SOMETHING",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NoValueProvidedTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"eq":          "",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("MissingKeyValueTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NullValueProvidedTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"eq":          nil,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("IncorrectKeyTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"neq":         "something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("IncorrectKeyCaseTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"something"}}]}`
		customVariables := map[string]interface{}{
			"EQ":          "something",
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("NumericDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123"}}]}`
		customVariables := map[string]interface{}{
			"eq":          12,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123.456"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("FloatDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"123.456"}}]}`
		customVariables := map[string]interface{}{
			"eq":          123.4567,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeMismatchTest", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"false"}}]}`
		customVariables := map[string]interface{}{
			"eq":          true,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})

	t.Run("BooleanDataTypeMismatchTest2", func(t *testing.T) {
		dsl := `{"or":[{"custom_variable":{"eq":"true"}}]}`
		customVariables := map[string]interface{}{
			"eq":          false,
			"expectation": false,
		}
		verifyExpectation(t, segmentationManager, dsl, customVariables)
	})
}

func verifyExpectation(t *testing.T, segmentationManager *segmentationCore.SegmentationManager, dsl string, customVariables map[string]interface{}) {
	expected := customVariables["expectation"].(bool)
	result := segmentationManager.ValidateSegmentation(dsl, customVariables)
	assert.Equal(t, expected, result)
}
