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

package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	settingsModel "github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
)

// UUID namespaces for UUID v5
var (
	dnsNamespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	urlNamespace = uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	seedURL      = "https://vwo.com"
)

// GetRandomUUID generates a random UUID based on an SDK key
func GetRandomUUID(sdkKey string) string {
	// Generate a namespace based on the SDK key using DNS namespace
	namespace := GenerateUUID(sdkKey, dnsNamespace)
	// Generate a random UUID (UUIDv4)
	randomUUID := uuid.New()
	// Generate a UUIDv5 using the random UUID and the namespace
	uuidv5 := GenerateUUID(randomUUID.String(), namespace)

	return uuidv5.String()
}

// GetUUID generates a UUID for a user based on their userId and accountId
func GetUUID(userID string, accountID string) string {
	// Generate a namespace UUID based on SEED_URL using URL namespace
	vwoNamespace := GenerateUUID(seedURL, urlNamespace)

	// Ensure userId and accountId are strings
	userIDStr := userID
	if userIDStr == "" {
		userIDStr = ""
	}
	accountIDStr := accountID
	if accountIDStr == "" {
		accountIDStr = ""
	}

	// Generate a namespace UUID based on the accountId
	userIDNamespace := GenerateUUID(accountIDStr, vwoNamespace)

	// Generate a UUID based on the userId and the previously generated namespace
	uuidForUserIDAccountID := GenerateUUID(userIDStr, userIDNamespace)

	// Remove all dashes from the UUID and convert it to uppercase
	desiredUUID := strings.ToUpper(strings.ReplaceAll(uuidForUserIDAccountID.String(), "-", ""))
	return desiredUUID
}

// GenerateUUID generates a UUID v5 based on a name and a namespace
func GenerateUUID(name string, namespace uuid.UUID) uuid.UUID {
	// Get namespace bytes
	namespaceBytes := uuidToBytes(namespace)
	nameBytes := []byte(name)

	// Combine namespace and name bytes
	combined := append(namespaceBytes, nameBytes...)

	// Generate SHA-1 hash
	hash := sha1.Sum(combined)

	// Set version to 5 (name-based using SHA-1)
	hash[6] = (hash[6] & 0x0f) | 0x50 // Version 5
	hash[8] = (hash[8] & 0x3f) | 0x80 // IETF variant

	// Convert hash to UUID
	return bytesToUUID(hash[:16])
}

// uuidToBytes converts a UUID to a byte array
func uuidToBytes(u uuid.UUID) []byte {
	bytes := make([]byte, 16)
	// UUID in Go is already [16]byte, so we can directly copy it
	copy(bytes, u[:])
	return bytes
}

// bytesToUUID converts a byte array to a UUID
func bytesToUUID(bytes []byte) uuid.UUID {
	var u uuid.UUID
	copy(u[:], bytes)
	return u
}

// GetUUIDFromBits converts most significant bits and least significant bits to UUID
func GetUUIDFromBits(msb, lsb uint64) uuid.UUID {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes[0:8], msb)
	binary.BigEndian.PutUint64(bytes[8:16], lsb)
	return bytesToUUID(bytes)
}

// IsWebUuid reports whether id is a web-generated UUID.
// It checks that id looks like a web-generated context ID: D or J + 32 hex chars = 33 chars total.
func IsWebUuid(id string) bool {
	return constants.WebUUIDRegex.MatchString(id)
}

// getUUIDFromContext generates a UUID from context ID, taking into account
// web connectivity settings and optional useIdForWeb flag on the context.
func GetUUIDFromContext(context map[string]interface{}, apiName string, accountID int, processedSettings *settingsModel.Settings) (uuid string, isWebUuid bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error getting UUID from context: %v", r)
			uuid = ""
			isWebUuid = false
			return
		}
	}()

	userId, ok := context[enums.ContextID.GetValue()].(string)
	if !ok {
		return "", false, fmt.Errorf("context ID is not a string")
	}

	if processedSettings.GetIsWebConnectivityEnabled() != false {
		// if web connectivity is enabled, check if context ID is a valid web UUID
		if IsWebUuid(userId) {
			// if context ID is a valid web UUID, set it as uuid
			return userId, true, nil
		}

		// get useIdForWeb from context
		useIdForWeb, ok := context[enums.ContextUseIdForWeb.GetValue()].(bool)
		if !ok {
			useIdForWeb = false
		}

		// if context useIdForWeb is true and context ID is not a valid web UUID, return error
		if useIdForWeb {
			return "", false, fmt.Errorf("UUID passed in context id is not a valid UUID")
		}

		// if context useIdForWeb is false, fallback to server‑side UUID derivation
		return GetUUID(userId, strconv.Itoa(accountID)), false, nil
	}

	// if web connectivity is disabled, fallback to server-side UUID derivation
	return GetUUID(userId, strconv.Itoa(accountID)), false, nil
}
