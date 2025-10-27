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

func TestCombination(t *testing.T) {
	// Initialize segmentation manager
	logManager := loggerCore.NewLogManager(nil)
	segmentationManager := segmentationCore.NewSegmentationManagerWithEvaluator(logManager, true)

	t.Run("DslWithAllOperandsTest", func(t *testing.T) {
		dsl := `{"or":[{"or":[{"and":[{"or":[{"custom_variable":{"start_with":"wildcard(my_start_with_val*)"}}]},{"not":{"or":[{"custom_variable":{"neq":"not_eq_value"}}]}}]},{"or":[{"custom_variable":{"contain":"wildcard(*my_contain_val*)"}}]}]},{"and":[{"or":[{"custom_variable":{"eq":"eq_value"}}]},{"or":[{"custom_variable":{"reg":"regex(myregex+)"}}]}]}]}`

		testCases := map[string]map[string]interface{}{
			"matchingStartWithValue": {
				"start_with":  "my_start_with_valzzzzzzzzzzzzzzzz",
				"neq":         1,
				"contain":     1,
				"eq":          1,
				"reg":         1,
				"expectation": true,
			},
			"matchingNotEqualToValue": {
				"start_with":  1,
				"neq":         "not_eq_value",
				"contain":     1,
				"eq":          1,
				"reg":         1,
				"expectation": false,
			},
			"matchingBothStartWithAndNotEqualToValue": {
				"start_with":  "my_start_with_valzzzzzzzzzzzzzzzz",
				"neq":         "not_eq_value",
				"contain":     1,
				"eq":          1,
				"reg":         1,
				"expectation": false,
			},
			"matchingContainsWithValue": {
				"start_with":  "m1y_1sta1rt_with_val",
				"neq":         false,
				"contain":     "zzzzzzmy_contain_valzzzzz",
				"eq":          1,
				"reg":         1,
				"expectation": true,
			},
			"matchingEqualToValue": {
				"start_with":  "m1y_1sta1rt_with_val",
				"neq":         nil,
				"contain":     "my_ contain _val",
				"eq":          "eq_value",
				"reg":         1,
				"expectation": false,
			},
			"matchingRegexValue": {
				"start_with":  "m1y_1sta1rt_with_val",
				"neq":         123,
				"contain":     "my_ contain _val",
				"eq":          "eq__value",
				"reg":         "myregexxxxxx",
				"expectation": false,
			},
			"matchingBothEqualToAndRegexValue": {
				"start_with":  "m1y_1sta1rt_with_val",
				"neq":         "not_matching",
				"contain":     "my$contain$val",
				"eq":          "eq_value",
				"reg":         "myregexxxxxx",
				"expectation": true,
			},
		}

		validateAllCases(t, segmentationManager, dsl, testCases)
	})

	t.Run("EmptyDslTest", func(t *testing.T) {
		dsl := `{}`
		testCases := map[string]map[string]interface{}{
			"matchingStartWithValue": {
				"start_with":  "m1y_1sta1rt_with_val",
				"neq":         nil,
				"contain":     "my_ contain _val",
				"eq":          "eq_value",
				"reg":         1,
				"expectation": false,
			},
		}

		validateAllCases(t, segmentationManager, dsl, testCases)
	})

	t.Run("DslWithAllOperandsTest2", func(t *testing.T) {
		dsl := `{"or":[{"and":[{"and":[{"not":{"or":[{"custom_variable":{"notvwo":"notvwo"}}]}},{"or":[{"custom_variable":{"vwovwovwo":"regex(vwovwovwo)"}}]}]},{"or":[{"custom_variable":{"regex_vwo":"regex(this\\s+is\\s+vwo)"}}]}]},{"and":[{"and":[{"not":{"or":[{"custom_variable":{"vwo_not_equal_to":"owv"}}]}},{"or":[{"custom_variable":{"vwo_equal_to":"vwo"}}]}]},{"or":[{"or":[{"custom_variable":{"vwo_starts_with":"wildcard(owv vwo*)"}}]},{"or":[{"custom_variable":{"vwo_contains":"wildcard(*vwo vwo vwo vwo vwo*)"}}]}]}]}]}`

		testCases := map[string]map[string]interface{}{
			"false_1": {
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"vwo_contains":     "vwo",
				"vwo_equal_to":     "vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_starts_with":  "v owv vwo",
				"vwovwovwo":        "vwovovwo",
				"expectation":      false,
			},
			"false_2": {
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"vwo_contains":     "vwo",
				"vwo_equal_to":     "vwovwo",
				"vwo_not_equal_to": "vwo",
				"vwo_starts_with":  "owv vwo",
				"vwovwovwo":        "vwovw",
				"expectation":      false,
			},
			"false_3": {
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"vwo_contains":     "vwo",
				"vwo_equal_to":     "vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_starts_with":  "vwo owv vwo",
				"vwovwovwo":        "vwovwovw",
				"expectation":      false,
			},
			"false_4": {
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"vwo_contains":     "vwo",
				"vwo_equal_to":     "vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_starts_with":  "vwo owv vwo",
				"vwovwovwo":        "vwo",
				"expectation":      false,
			},
		}

		validateAllCases(t, segmentationManager, dsl, testCases)
	})

	t.Run("DslWithAllOperandsTest3", func(t *testing.T) {
		dsl := `{"or":[{"and":[{"and":[{"not":{"or":[{"custom_variable":{"notvwo":"notvwo"}}]}},{"or":[{"custom_variable":{"vwovwovwo":"regex(vwovwovwo)"}}]}]},{"or":[{"custom_variable":{"regex_vwo":"regex(this\\s+is\\s+vwo)"}}]}]},{"and":[{"and":[{"not":{"or":[{"custom_variable":{"vwo_not_equal_to":"owv"}}]}},{"or":[{"custom_variable":{"vwo_equal_to":"vwo"}}]}]},{"or":[{"or":[{"custom_variable":{"vwo_starts_with":"wildcard(owv vwo*)"}}]},{"or":[{"custom_variable":{"vwo_contains":"wildcard(*vwo vwo vwo vwo vwo*)"}}]}]}]}]}`

		testCases := map[string]map[string]interface{}{
			"true_1": {
				"vwo_starts_with":  "vwo owv vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_equal_to":     "vwo",
				"notvwo":           "vo",
				"regex_vwo":        "this   is vwo",
				"vwovwovwo":        "vwovwovwo",
				"vwo_contains":     "vw",
				"expectation":      true,
			},
			"true_2": {
				"vwo_starts_with":  "owv vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_equal_to":     "vwo",
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"vwovwovwo":        "vwovwovwo",
				"vwo_contains":     "vwo",
				"expectation":      true,
			},
			"true_3": {
				"vwo_starts_with":  "owv vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_equal_to":     "vwo",
				"notvwo":           "vwovwo",
				"regex_vwo":        "this   isvwo",
				"vwovwovwo":        "vwovwovwo",
				"vwo_contains":     "vwo",
				"expectation":      true,
			},
			"true_4": {
				"vwo_starts_with":  "owv vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_equal_to":     "vwo",
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"vwovwovwo":        "vwo",
				"vwo_contains":     "vwo",
				"expectation":      true,
			},
		}

		validateAllCases(t, segmentationManager, dsl, testCases)
	})

	t.Run("DslWithAllOperandsTest4", func(t *testing.T) {
		dsl := `{"and":[{"or":[{"custom_variable":{"contains_vwo":"wildcard(*vwo*)"}}]},{"and":[{"and":[{"or":[{"and":[{"or":[{"and":[{"or":[{"custom_variable":{"regex_for_all_letters":"regex(^[A-z]+$)"}}]},{"or":[{"custom_variable":{"regex_for_capital_letters":"regex(^[A-Z]+$)"}}]}]},{"or":[{"custom_variable":{"regex_for_small_letters":"regex(^[a-z]+$)"}}]}]},{"or":[{"custom_variable":{"regex_for_no_zeros":"regex(^[1-9]+$)"}}]}]},{"or":[{"custom_variable":{"regex_for_zeros":"regex(^[0]+$)"}}]}]},{"or":[{"custom_variable":{"regex_real_number":"regex(^\\\\d+(\\\\.\\\\d+)?)"}}]}]},{"or":[{"or":[{"custom_variable":{"this_is_regex":"regex(this\\\\s+is\\\\s+text)"}}]},{"and":[{"and":[{"or":[{"custom_variable":{"starts_with":"wildcard(starts_with_variable*)"}}]},{"or":[{"custom_variable":{"contains":"wildcard(*contains_variable*)"}}]}]},{"or":[{"not":{"or":[{"custom_variable":{"is_not_equal_to":"is_not_equal_to_variable"}}]}},{"or":[{"custom_variable":{"is_equal_to":"equal_to_variable"}}]}]}]}]}]}]}`

		testCases := map[string]map[string]interface{}{
			"false_5": {
				"contains":                  "contains_variable",
				"contains_vwo":              "legends say that vwo is the best",
				"is_equal_to":               "equal_to_variable",
				"is_not_equal_to":           "is_not_equal_to_variable",
				"regex_for_all_letters":     "dsfASF6",
				"regex_for_capital_letters": "SADFLSDLF",
				"regex_for_no_zeros":        12231023,
				"regex_for_small_letters":   "sadfksjdf",
				"regex_for_zeros":           "0001000",
				"regex_real_number":         12321.2242,
				"starts_with":               "starts_with_variable",
				"this_is_regex":             "this    is    regex",
				"expectation":               false,
			},
			"true_5": {
				"contains":                  "contains_variable",
				"contains_vwo":              "legends say that vwo is the best",
				"is_equal_to":               "equal_to_variable",
				"is_not_equal_to":           "is_not_equal_to_variable",
				"regex_for_all_letters":     "dsfASF",
				"regex_for_capital_letters": "SADFLSDLF",
				"regex_for_no_zeros":        1223123,
				"regex_for_small_letters":   "sadfksjdf",
				"regex_for_zeros":           0,
				"regex_real_number":         12321.2242,
				"starts_with":               "starts_with_variable",
				"this_is_regex":             "this    is    regex",
				"expectation":               false,
			},
		}

		validateAllCases(t, segmentationManager, dsl, testCases)
	})

	t.Run("DslWithAllOperandsTest5", func(t *testing.T) {
		dsl := `{"and":[{"or":[{"and":[{"not":{"or":[{"custom_variable":{"thanos":"snap"}}]}},{"or":[{"custom_variable":{"batman":"wildcard(*i am batman*)"}}]}]},{"or":[{"custom_variable":{"joker":"regex((joker)+)"}}]}]},{"and":[{"or":[{"or":[{"custom_variable":{"lol":"lolololololol"}}]},{"or":[{"custom_variable":{"blablabla":"wildcard(*bla*)"}}]}]},{"and":[{"and":[{"not":{"or":[{"custom_variable":{"notvwo":"notvwo"}}]}},{"or":[{"and":[{"or":[{"custom_variable":{"vwovwovwo":"regex(vwovwovwo)"}}]},{"or":[{"custom_variable":{"regex_vwo":"regex(this\\s+is\\s+vwo)"}}]}]},{"or":[{"and":[{"not":{"or":[{"custom_variable":{"vwo_not_equal_to":"owv"}}]}},{"or":[{"custom_variable":{"vwo_equal_to":"vwo"}}]}]},{"or":[{"custom_variable":{"vwo_starts_with":"wildcard(owv vwo*)"}}]}]}]}]},{"or":[{"custom_variable":{"vwo_contains":"wildcard(*vwo vwo vwo vwo vwo*)"}}]}]}]}]}`

		testCases := map[string]map[string]interface{}{
			"false_11": {
				"batman":           "hello i am batman world",
				"blablabla":        "lba",
				"joker":            "joker joker joker",
				"lol":              "lollolololol",
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"thanos":           "snap",
				"vwo_contains":     "vwo vwo vwo vwo vwo",
				"vwo_equal_to":     "vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_starts_with":  "vwo",
				"vwovwovwo":        "vwovwovwo",
				"expectation":      false,
			},
			"true_9": {
				"batman":           "hello i am batman world",
				"blablabla":        "bla bla bla",
				"joker":            "joker joker joker",
				"lol":              "lollolololol",
				"notvwo":           "vwo",
				"regex_vwo":        "this   is vwo",
				"thanos":           "half universe",
				"vwo_contains":     "vwo vwo vwo vwo vwo",
				"vwo_equal_to":     "vwo",
				"vwo_not_equal_to": "vwo",
				"vwo_starts_with":  "owv vwo",
				"vwovwovwo":        "vwovwovwo",
				"expectation":      true,
			},
		}

		validateAllCases(t, segmentationManager, dsl, testCases)
	})
}

func validateAllCases(t *testing.T, segmentationManager *segmentationCore.SegmentationManager, dsl string, testCases map[string]map[string]interface{}) {
	for testName, customVariables := range testCases {
		t.Run(testName, func(t *testing.T) {
			expected := customVariables["expectation"].(bool)
			result := segmentationManager.ValidateSegmentation(dsl, customVariables)
			assert.Equal(t, expected, result)
		})
	}
}
