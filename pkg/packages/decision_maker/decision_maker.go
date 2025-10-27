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

package decision_maker

import (
	"math"

	"github.com/spaolacci/murmur3"
	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
)

// GenerateBucketValue generates a bucket value based on hash value, max value, and multiplier
func GenerateBucketValue(hashValue uint32, maxValue int, multiplier int) int {
	ratio := float64(hashValue) / math.Pow(2, 32)
	multipliedValue := (float64(maxValue)*ratio + 1) * float64(multiplier)
	return int(math.Floor(multipliedValue))
}

// GetBucketValueForUser validates user ID and generates a bucket value
func GetBucketValueForUser(userID string) int {
	if userID == "" {
		return 0
	}
	hashValue := GenerateHashValue(userID)
	return GenerateBucketValue(hashValue, constants.MAX_CAMPAIGN_VALUE, 1)
}

func GetBucketValueForUserTest(userID string, maxValue int) int {
	if userID == "" {
		return 0
	}
	hashValue := GenerateHashValue(userID)
	return GenerateBucketValue(hashValue, maxValue, 1)
}

// CalculateBucketValue calculates bucket value for a given string with multiplier and max value
func CalculateBucketValue(str string) int {
	hashValue := GenerateHashValue(str)
	return GenerateBucketValue(hashValue, constants.MAX_TRAFFIC_VALUE, 1)
}

// GenerateHashValue generates a hash value using MurmurHash3
func GenerateHashValue(hashKey string) uint32 {
	// Use MurmurHash3 32-bit with seed value of 1
	hasher := murmur3.New32WithSeed(uint32(constants.SeedValue))
	hasher.Write([]byte(hashKey))
	return hasher.Sum32()
}
