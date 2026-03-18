/**
 * Copyright 2025 Wingify Software Pvt. Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package e2e

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingify/vwo-fme-go-sdk"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
)

const mockSettingsFile = `{
    "version": 1,
    "sdkKey": "abcdef",
    "accountId": 123456,
    "campaigns": [
        {
            "segments": {},
            "status": "RUNNING",
            "variations": [
                {
                    "weight": 100,
                    "segments": {},
                    "id": 1,
                    "variables": [
                        {
                            "id": 1,
                            "type": "string",
                            "value": "def",
                            "key": "kaus"
                        }
                    ],
                    "name": "Rollout-rule-1"
                }
            ],
            "type": "FLAG_ROLLOUT",
            "isAlwaysCheckSegment": false,
            "isForcedVariationEnabled": false,
            "name": "featureOne : Rollout",
            "key": "featureOne_rolloutRule1",
            "percentTraffic": 0,
            "id": 1
        },
        {
            "segments": {},
            "status": "RUNNING",
            "key": "featureOne_testingRule1",
            "type": "FLAG_TESTING",
            "isAlwaysCheckSegment": false,
            "name": "featureOne : Testing rule 1",
            "isForcedVariationEnabled": true,
            "id": 2, 
            "percentTraffic": 100,
            "variations": [
                {
                    "weight": 50,
                    "segments": {},
                    "id": 1,
                    "variables": [
                        {
                            "id": 1,
                            "type": "string",
                            "value": "def",
                            "key": "kaus"
                        }
                    ],
                    "name": "Default"
                },
                {
                    "weight": 50,
                    "segments": {},
                    "id": 2,
                    "variables": [
                        {
                            "id": 1,
                            "type": "string",
                            "value": "var1",
                            "key": "kaus"
                        }
                    ],
                    "name": "Variation-1"
                },
                {
                    "weight": 0.0001,
                    "segments": {
                        "or": [
                            {
                                "user": "forcedWingify"
                            }
                        ]
                    },
                    "id": 3,
                    "variables": [
                        {
                            "id": 1,
                            "type": "string",
                            "value": "var2",
                            "key": "kaus"
                        }
                    ],
                    "name": "Variation-2"
                }
            ]
        }
    ],
    "features": [
        {
            "impactCampaign": {},
            "rules": [
                {
                    "type": "FLAG_TESTING",
                    "ruleKey": "testingRule1",
                    "campaignId": 2
                }
            ],
            "status": "ON",
            "key": "featureOne",
            "metrics": [
                {
                    "type": "CUSTOM_GOAL",
                    "identifier": "e1",
                    "id": 1
                }
            ],
            "type": "FEATURE_FLAG",
            "name": "featureOne",
            "id": 1
        }
    ]
}`
const mockSettingsWithSameSaltFile = `{
    "features": [{
      "key": "feature1",
      "name": "Feature1",
      "metrics": [{
        "id": 1,
        "type": "REVENUE_TRACKING",
        "identifier": "custom1",
        "mca": -1
      }],
      "rules": [{
          "campaignId": 2,
          "type": "FLAG_TESTING",
          "ruleKey": "testingRule1"
        }
      ],
      "type": "FEATURE_FLAG",
      "impactCampaign": {},
      "id": 1,
      "status": "ON"
    },
    {
        "key": "feature2",
        "name": "Feature2",
        "metrics": [{
          "id": 1,
          "type": "REVENUE_TRACKING",
          "identifier": "custom1",
          "mca": -1
        }],
        "rules": [
          {
            "campaignId": 4,
            "type": "FLAG_TESTING",
            "ruleKey": "testingRule1"
          }
        ],
        "type": "FEATURE_FLAG",
        "impactCampaign": {},
        "id": 2,
        "status": "ON"
      }
    ],
    "version": 1,
    "accountId": 12345,
    "sdkKey": "000000000000_MASKED_000000000000",
    "campaigns": [{
        "key": "feature1_rolloutRule1",
        "name": "feature1_rolloutRule1",
        "id": 1,
        "segments": {},
        "isForcedVariationEnabled": false,
        "variations": [{
          "variables": [{
              "key": "int",
              "id": 1,
              "value": 10,
              "type": "integer"
            },
            {
              "key": "float",
              "id": 2,
              "value": 20.01,
              "type": "double"
            },
            {
              "key": "string",
              "id": 3,
              "value": "test",
              "type": "string"
            },
            {
              "key": "boolean",
              "id": 4,
              "value": false,
              "type": "boolean"
            },
            {
              "key": "json",
              "id": 5,
              "value": "{\"name\": \"varun\"}",
              "type": "json"
            }
          ],
          "id": 1,
          "salt": "rolloutSalt",
          "segments": {},
          "weight": 100,
          "name": "Rollout-rule-1"
        }],
        "type": "FLAG_ROLLOUT",
        "status": "RUNNING"
      },
      {
        "key": "feature1_testingRule1",
        "name": "feature1_testingRule1",
        "id": 2,
        "segments": {},
        "salt": "testingSalt",
        "isForcedVariationEnabled": false,
        "variations": [{
            "weight": 50,
            "id": 1,
            "variables": [{
                "key": "int",
                "id": 1,
                "value": 10,
                "type": "integer"
              },
              {
                "key": "float",
                "id": 2,
                "value": 20.01,
                "type": "double"
              },
              {
                "key": "string",
                "id": 3,
                "value": "test",
                "type": "string"
              },
              {
                "key": "boolean",
                "id": 4,
                "value": false,
                "type": "boolean"
              },
              {
                "key": "json",
                "id": 5,
                "value": "{\"name\": \"varun\"}",
                "type": "json"
              }
            ],
            "name": "Default"
          },
          {
            "weight": 50,
            "id": 2,
            "variables": [{
                "key": "int",
                "id": 1,
                "value": 11,
                "type": "integer"
              },
              {
                "key": "float",
                "id": 2,
                "value": 20.02,
                "type": "double"
              },
              {
                "key": "string",
                "id": 3,
                "value": "test_variation",
                "type": "string"
              },
              {
                "key": "boolean",
                "id": 4,
                "value": true,
                "type": "boolean"
              },
              {
                "key": "json",
                "id": 5,
                "value": {
                  "variation": 1,
                  "name": "VWO"
                },
                "type": "json"
              }
            ],
            "name": "Variation-1"
          }
        ],
        "percentTraffic": 100,
        "type": "FLAG_TESTING",
        "status": "RUNNING"
      },
      {
        "key": "feature2_rolloutRule1",
        "name": "feature2_rolloutRule1",
        "id": 3,
        "segments": {},
        "isForcedVariationEnabled": false,
        "variations": [{
          "variables": [{
              "key": "int",
              "id": 1,
              "value": 10,
              "type": "integer"
            },
            {
              "key": "float",
              "id": 2,
              "value": 20.01,
              "type": "double"
            },
            {
              "key": "string",
              "id": 3,
              "value": "test",
              "type": "string"
            },
            {
              "key": "boolean",
              "id": 4,
              "value": false,
              "type": "boolean"
            },
            {
              "key": "json",
              "id": 5,
              "value": "{\"name\": \"varun\"}",
              "type": "json"
            }
          ],
          "id": 1,
          "salt": "rolloutSalt",
          "segments": {},
          "weight": 100,
          "name": "Rollout-rule-1"
        }],
        "type": "FLAG_ROLLOUT",
        "status": "RUNNING"
      },
      {
        "key": "feature2_testingRule1",
        "name": "feature2_testingRule1",
        "id": 4,
        "segments": {},
        "salt": "testingSalt",
        "isForcedVariationEnabled": false,
        "variations": [{
            "weight": 50,
            "id": 1,
            "variables": [{
                "key": "int",
                "id": 1,
                "value": 10,
                "type": "integer"
              },
              {
                "key": "float",
                "id": 2,
                "value": 20.01,
                "type": "double"
              },
              {
                "key": "string",
                "id": 3,
                "value": "test",
                "type": "string"
              },
              {
                "key": "boolean",
                "id": 4,
                "value": false,
                "type": "boolean"
              },
              {
                "key": "json",
                "id": 5,
                "value": "{\"name\": \"varun\"}",
                "type": "json"
              }
            ],
            "name": "Default"
          },
          {
            "weight": 50,
            "id": 2,
            "variables": [{
                "key": "int",
                "id": 1,
                "value": 11,
                "type": "integer"
              },
              {
                "key": "float",
                "id": 2,
                "value": 20.02,
                "type": "double"
              },
              {
                "key": "string",
                "id": 3,
                "value": "test_variation",
                "type": "string"
              },
              {
                "key": "boolean",
                "id": 4,
                "value": true,
                "type": "boolean"
              },
              {
                "key": "json",
                "id": 5,
                "value": {
                  "variation": 1,
                  "name": "VWO"
                },
                "type": "json"
              }
            ],
            "name": "Variation-1"
          }
        ],
        "percentTraffic": 100,
        "type": "FLAG_TESTING",
        "status": "RUNNING"
      }
    ]
  }`

// setupVWOClient initializes the VWO Client with the mock settings
func setupVWOClient(t *testing.T) *vwo.VWOClient {
    options := map[string]interface{}{
        enums.OptionSDKKey.GetValue():   "abcdef",
        enums.OptionAccountID.GetValue(): 123456,
        enums.OptionSettings.GetValue():  mockSettingsFile,
    }

    vwoClient, err := vwo.Init(options)
    assert.NoError(t, err)
    assert.NotNil(t, vwoClient)
    return vwoClient
}

func TestCustomBucketingSeed(t *testing.T) {
    vwoClient := setupVWOClient(t)

    // Case 1: Standard bucketing (no custom seed)
    // Scenario: Two different users ('user1', 'user3') with NO bucketing seed.
    // Expected: They should be bucketed into different variations based on their User IDs.
    t.Run("getFlag without bucketing seed: should assign different variations to users with different user IDs", func(t *testing.T) {
        context1 := map[string]interface{}{
            enums.ContextID.GetValue(): "user1",
        }
        context2 := map[string]interface{}{
            enums.ContextID.GetValue(): "user3",
        }

        user1Flag, err1 := vwoClient.GetFlag("featureOne", context1)
        user2Flag, err2 := vwoClient.GetFlag("featureOne", context2)

        assert.NoError(t, err1)
        assert.NoError(t, err2)
        assert.NotEqual(t, user1Flag.GetVariables(), user2Flag.GetVariables())
    })

    // Case 2: Bucketing Seed Provided
    // Scenario: Two different users ('user1', 'user3') are provided with the SAME bucketingSeed.
    // Expected: Since the seed is identical, they MUST get the same variation.
    t.Run("getFlag with bucketing seed: should assign same variation to different users with same bucketing seed", func(t *testing.T) {
        sameBucketingSeed := "common-seed-123"

        context1 := map[string]interface{}{
            enums.ContextID.GetValue():            "user1",
            enums.ContextBucketingSeed.GetValue(): sameBucketingSeed,
        }
        context2 := map[string]interface{}{
            enums.ContextID.GetValue():            "user3",
            enums.ContextBucketingSeed.GetValue(): sameBucketingSeed,
        }

        user1Flag, err1 := vwoClient.GetFlag("featureOne", context1)
        user2Flag, err2 := vwoClient.GetFlag("featureOne", context2)

        assert.NoError(t, err1)
        assert.NoError(t, err2)
        assert.Equal(t, user1Flag.GetVariables(), user2Flag.GetVariables())
    })



    // Case 3: Different Seeds
    // Scenario: The SAME User ID is used, but with DIFFERENT bucketing seeds.
    // Expected: The SDK should bucket based on the seed. Since we use seeds known to produce different results, the outcomes should differ.
    t.Run("getFlag with bucketing seed: should assign different variations to users with different bucketing seeds", func(t *testing.T) {
        context1 := map[string]interface{}{
            enums.ContextID.GetValue():            "sameId",
            enums.ContextBucketingSeed.GetValue(): "user1",
        }
        context2 := map[string]interface{}{
            enums.ContextID.GetValue():            "sameId",
            enums.ContextBucketingSeed.GetValue(): "user3",
        }

        user1Flag, err1 := vwoClient.GetFlag("featureOne", context1)
        user2Flag, err2 := vwoClient.GetFlag("featureOne", context2)

        assert.NoError(t, err1)
        assert.NoError(t, err2)
        assert.NotEqual(t, user1Flag.GetVariables(), user2Flag.GetVariables())
    })

    // Case 4: Empty String Seed
    // Scenario: bucketingSeed is provided but it's an empty string.
    // Expected: Empty string is falsy, so it should fall back to userId. Different users should get different variations.
    t.Run("getFlag with bucketing seed: should fallback to userId when bucketingSeed is empty string", func(t *testing.T) {
        context1 := map[string]interface{}{
            enums.ContextID.GetValue():            "user1",
            enums.ContextBucketingSeed.GetValue(): "",
        }
        context2 := map[string]interface{}{
            enums.ContextID.GetValue():            "user3",
            enums.ContextBucketingSeed.GetValue(): "",
        }

        user1Flag, err1 := vwoClient.GetFlag("featureOne", context1)
        user2Flag, err2 := vwoClient.GetFlag("featureOne", context2)
        //No error fallback to userId and different users should get different variations
        assert.NoError(t, err1)
        assert.NoError(t, err2)
        t.Logf("User1 variables: %v", user1Flag.GetVariables())
        t.Logf("User2 variables: %v", user2Flag.GetVariables())
        assert.NotEqual(t, user1Flag.GetVariables(), user2Flag.GetVariables())
    })

    // Case 5: Custom Salt and Bucketing Seed Combinations
    t.Run("getFlag with custom salt and bucketing seed combinations", func(t *testing.T) {
        options := map[string]interface{}{
            enums.OptionSDKKey.GetValue():   "abcdef",
            enums.OptionAccountID.GetValue(): 123456,
            enums.OptionSettings.GetValue():  mockSettingsWithSameSaltFile,
        }

        vwoClientSalt, err := vwo.Init(options)
        assert.NoError(t, err)
        assert.NotNil(t, vwoClientSalt)
        //No bucketing seed, custom salt present - 10 users, randomly distributed, but each user getting same variation in both flags
        t.Run("No bucketing seed, custom salt present - 10 users, randomly distributed, but each user getting same variation in both flags", func(t *testing.T) {
            for i := 1; i <= 10; i++ {
                userID := "user" + strconv.Itoa(i)
                contextUser := map[string]interface{}{
                    enums.ContextID.GetValue(): userID,
                }

                flag1, _ := vwoClientSalt.GetFlag("feature1", contextUser)
                flag2, _ := vwoClientSalt.GetFlag("feature2", contextUser)
                //Both feature flag should have same variation
                assert.Equal(t, flag1.GetVariables(), flag2.GetVariables())
            }
        })
        //Bucketing seed present, salt present - 10 users, all users getting same variation in both flags
        t.Run("Bucketing seed present, salt present - 10 users, all users getting same variation in both flags", func(t *testing.T) {
            commonBucketingSeed := "common_seed_456"
            variationsAssigned := make(map[string]bool)

            for i := 1; i <= 10; i++ {
                userID := "user" + strconv.Itoa(i)
                contextUser := map[string]interface{}{
                    enums.ContextID.GetValue():            userID,
                    enums.ContextBucketingSeed.GetValue(): commonBucketingSeed,
                }
                //Both feature flag should have same variation
                flag1, _ := vwoClientSalt.GetFlag("feature1", contextUser)
                flag2, _ := vwoClientSalt.GetFlag("feature2", contextUser)

                assert.Equal(t, flag1.GetVariables(), flag2.GetVariables())

                variationsAssigned[fmt.Sprintf("%v", flag1.GetVariables())] = true
            }

            // Since the bucketing seed is the exact same for all 10 users, they MUST all get the same variation
            assert.Equal(t, 1, len(variationsAssigned))
        })
    })
    //Case 6: Forced variation (whitelisting) and bucketing seed
    t.Run("getFlag with forced variation (whitelisting) and bucketing seed", func(t *testing.T) {
        // In MOCK_SETTINGS_FILE, 'forcedWingify' is whitelisted to Variation-2 (value: 'var2').
        t.Run("should return forced variation for whitelisted user without bucketing seed", func(t *testing.T) {
            context1 := map[string]interface{}{
                enums.ContextID.GetValue(): "forcedWingify",
            }
            forcedUserFlag, err := vwoClient.GetFlag("featureOne", context1)
            assert.NoError(t, err)
            assert.Equal(t, "var2", forcedUserFlag.GetVariable("kaus", ""))
        })

        // Even with a bucketing seed, forcedWingify must still get the forced variation (Variation-2, value: 'var2')
        t.Run("should still return forced variation for whitelisted user when bucketing seed is present", func(t *testing.T) {
            context1 := map[string]interface{}{
                enums.ContextID.GetValue():            "forcedWingify",
                enums.ContextBucketingSeed.GetValue(): "some-seed-xyz",
            }
            forcedUserFlag, err := vwoClient.GetFlag("featureOne", context1)
            assert.NoError(t, err)
            assert.Equal(t, "var2", forcedUserFlag.GetVariable("kaus", ""))
        })
    })
}