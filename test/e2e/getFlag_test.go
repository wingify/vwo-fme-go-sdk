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

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingify/vwo-fme-go-sdk"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/test/data"
	"github.com/wingify/vwo-fme-go-sdk/test/data/testCases"
)

const (
	SDK_KEY    = "abcd"
	ACCOUNT_ID = 12345
)

func TestGetFlagWithoutStorage(t *testing.T) {
	settingsReader := data.NewDummySettingsReader()
	testDataReader := data.NewTestDataReader()

	runTests(t, testDataReader.TestCases.GetFlagWithoutStorage, false, settingsReader.SettingsMap)
}

func TestGetFlagWithSalt(t *testing.T) {
	settingsReader := data.NewDummySettingsReader()
	testDataReader := data.NewTestDataReader()

	runSaltTest(t, testDataReader.TestCases.GetFlagWithSalt, settingsReader.SettingsMap)
}

func TestGetFlagWithMegRandom(t *testing.T) {
	settingsReader := data.NewDummySettingsReader()
	testDataReader := data.NewTestDataReader()

	runTests(t, testDataReader.TestCases.GetFlagMegRandom, false, settingsReader.SettingsMap)
}

func TestGetFlagWithMegAdvance(t *testing.T) {
	settingsReader := data.NewDummySettingsReader()
	testDataReader := data.NewTestDataReader()

	runTests(t, testDataReader.TestCases.GetFlagMegAdvance, false, settingsReader.SettingsMap)
}

func TestGetFlagWithStorage(t *testing.T) {
	settingsReader := data.NewDummySettingsReader()
	testDataReader := data.NewTestDataReader()

	runTests(t, testDataReader.TestCases.GetFlagWithStorage, true, settingsReader.SettingsMap)
}

func runTests(t *testing.T, tests []testCases.TestData, useStorage bool, settingsMap map[string]string) {
	for _, testData := range tests {
		t.Run(testData.Description, func(t *testing.T) {
			storage := data.NewStorageTest()

			options := map[string]interface{}{
				enums.OptionSDKKey.GetValue():    SDK_KEY,
				enums.OptionAccountID.GetValue(): ACCOUNT_ID,
				enums.OptionSettings.GetValue():  settingsMap[testData.Settings],
			}

			if useStorage {
				options[enums.OptionStorage.GetValue()] = storage
			}

			// Initialize VWO client
			vwoClient, err := vwo.Init(options)
			assert.NoError(t, err)
			assert.NotNil(t, vwoClient)

			if useStorage {
				// Check that storage is initially empty
				storageData, err := storage.Get(testData.FeatureKey, testData.Context.ID)
				assert.NoError(t, err)
				assert.Nil(t, storageData)
			}

			// Convert context to map for GetFlag call
			context := map[string]interface{}{
				enums.ContextID.GetValue(): testData.Context.ID,
			}

			// Add custom variables if present
			if testData.Context.CustomVariables != nil {
				context[enums.ContextCustomVariables.GetValue()] = testData.Context.CustomVariables
			}

			// Get feature flag
			featureFlag, err := vwoClient.GetFlag(testData.FeatureKey, context)
			assert.NoError(t, err)
			assert.NotNil(t, featureFlag)

			// Assert expectations
			if testData.Expectation.IsEnabled != nil {
				assert.Equal(t, *testData.Expectation.IsEnabled, featureFlag.IsEnabled())
			}

			if testData.Expectation.IntVariable != nil {
				intVar := featureFlag.GetVariable("int", 1)
				expected := int64(*testData.Expectation.IntVariable)
				switch v := intVar.(type) {
				case int64:
					assert.Equal(t, expected, v)
				case int:
					assert.Equal(t, expected, int64(v))
				case float64:
					assert.Equal(t, expected, int64(v))
				default:
					t.Fatalf("unexpected int variable type: %T", intVar)
				}
			}

			if testData.Expectation.StringVariable != nil {
				stringVar := featureFlag.GetVariable("string", "VWO")
				assert.Equal(t, *testData.Expectation.StringVariable, stringVar)
			}

			if testData.Expectation.FloatVariable != nil {
				floatVar := featureFlag.GetVariable("float", 1.1)
				assert.Equal(t, *testData.Expectation.FloatVariable, floatVar)
			}

			if testData.Expectation.BooleanVariable != nil {
				boolVar := featureFlag.GetVariable("boolean", false)
				assert.Equal(t, *testData.Expectation.BooleanVariable, boolVar)
			}

			if testData.Expectation.JSONVariable != nil {
				jsonVar := featureFlag.GetVariable("json", map[string]interface{}{})
				assert.Equal(t, testData.Expectation.JSONVariable, jsonVar)
			}

			// Check storage data if enabled
			if useStorage && testData.Expectation.IsEnabled != nil && *testData.Expectation.IsEnabled {
				if testData.Expectation.StorageData != nil {
					storageData, err := storage.GetStorageData(testData.FeatureKey, testData.Context.ID)
					assert.NoError(t, err)
					assert.NotNil(t, storageData)

					if testData.Expectation.StorageData.RolloutKey != "" {
						assert.Equal(t, testData.Expectation.StorageData.RolloutKey, storageData.RolloutKey)
					}
					if testData.Expectation.StorageData.RolloutVariationID != 0 {
						assert.Equal(t, testData.Expectation.StorageData.RolloutVariationID, storageData.RolloutVariationID)
					}
					if testData.Expectation.StorageData.ExperimentKey != "" {
						assert.Equal(t, testData.Expectation.StorageData.ExperimentKey, storageData.ExperimentKey)
					}
					if testData.Expectation.StorageData.ExperimentVariationID != 0 {
						assert.Equal(t, testData.Expectation.StorageData.ExperimentVariationID, storageData.ExperimentVariationID)
					}
				}
			}
		})
		break
	}
}

func runSaltTest(t *testing.T, tests []testCases.TestData, settingsMap map[string]string) {
	for _, testData := range tests {
		t.Run(testData.Description, func(t *testing.T) {
			options := map[string]interface{}{
				enums.OptionSDKKey.GetValue():    SDK_KEY,
				enums.OptionAccountID.GetValue(): ACCOUNT_ID,
				enums.OptionSettings.GetValue():  settingsMap[testData.Settings],
			}

			// Initialize VWO client
			vwoClient, err := vwo.Init(options)
			assert.NoError(t, err)
			assert.NotNil(t, vwoClient)

			for _, userID := range testData.UserIds {
				context := map[string]interface{}{
					enums.ContextID.GetValue(): userID,
				}

				// Get feature flags for both feature keys
				featureFlag, err := vwoClient.GetFlag(testData.FeatureKey, context)
				assert.NoError(t, err)
				assert.NotNil(t, featureFlag)

				featureFlag2, err := vwoClient.GetFlag(testData.FeatureKey2, context)
				assert.NoError(t, err)
				assert.NotNil(t, featureFlag2)

				// Get variables for comparison
				featureFlagVariables := featureFlag.GetVariables()
				featureFlag2Variables := featureFlag2.GetVariables()

				if testData.Expectation.ShouldReturnSameVariation != nil {
					if *testData.Expectation.ShouldReturnSameVariation {
						assert.Equal(t, featureFlagVariables, featureFlag2Variables, "The feature flag variables should be equal!")
					} else {
						assert.NotEqual(t, featureFlagVariables, featureFlag2Variables, "The feature flag variables should not be equal!")
					}
				}
			}
		})
	}
}
