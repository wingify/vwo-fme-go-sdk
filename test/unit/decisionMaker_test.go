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
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/decision_maker"
)

func TestDecisionMaker(t *testing.T) {

	t.Run("TestGenerateBucketValue", func(t *testing.T) {
		hashValue := uint32(2147483647) // Example hash value
		maxValue := 100
		multiplier := 1
		expectedBucketValue := int(float64(maxValue)*float64(hashValue)/float64(1<<32) + 1)
		bucketValue := decision_maker.GenerateBucketValue(hashValue, maxValue, multiplier)
		assert.Equal(t, expectedBucketValue, bucketValue)
	})

	t.Run("TestGetBucketValueForUser", func(t *testing.T) {
		userID := "user123"

		bucketValue := decision_maker.GetBucketValueForUser(userID)
		assert.GreaterOrEqual(t, bucketValue, 1)
		assert.LessOrEqual(t, bucketValue, 100)
	})

	t.Run("TestCalculateBucketValue", func(t *testing.T) {
		str := "testString"

		bucketValue := decision_maker.CalculateBucketValue(str)
		assert.GreaterOrEqual(t, bucketValue, 1)
		assert.LessOrEqual(t, bucketValue, 10000)
	})

	t.Run("TestGenerateHashValue", func(t *testing.T) {
		hashKey := "key123"

		hashValue := decision_maker.GenerateHashValue(hashKey)
		assert.NotEqual(t, int64(0), hashValue)

		// Test that same key produces same hash
		hashValue2 := decision_maker.GenerateHashValue(hashKey)
		assert.Equal(t, hashValue, hashValue2)

		// Test that different keys produce different hashes
		hashValue3 := decision_maker.GenerateHashValue("differentKey")
		assert.NotEqual(t, hashValue, hashValue3)
	})

	t.Run("TestConsistency", func(t *testing.T) {
		userID := "testUser"

		// Test that the same user gets the same bucket value consistently
		bucketValue1 := decision_maker.GetBucketValueForUser(userID)
		bucketValue2 := decision_maker.GetBucketValueForUser(userID)
		assert.Equal(t, bucketValue1, bucketValue2)
	})

	t.Run("TestRange", func(t *testing.T) {
		userID := "testUser"
		maxValue := 50

		bucketValue := decision_maker.GetBucketValueForUserTest(userID, maxValue)
		assert.GreaterOrEqual(t, bucketValue, 1)
		assert.LessOrEqual(t, bucketValue, maxValue)
	})
}
