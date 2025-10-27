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

package decorators

import (
	"fmt"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// StorageDecorator handles storage operations with validation and logging
type StorageDecorator struct{}

// NewStorageDecorator creates a new StorageDecorator instance
func NewStorageDecorator() *StorageDecorator {
	return &StorageDecorator{}
}

// GetFeatureFromStorage retrieves feature data from storage
func (sd *StorageDecorator) GetFeatureFromStorage(
	featureKey string,
	context *user.VWOContext,
	storageService interfaces.StorageServiceInterface,
	serviceContainer interfaces.ServiceContainerInterface,
) map[string]interface{} {
	// Use defer with recover
	defer func() {
		if r := recover(); r != nil {
			serviceContainer.GetLoggerService().Error("ERROR_READING_DATA_FROM_STORAGE", map[string]interface{}{"err": fmt.Sprint(r)}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		}
	}()

	data, err := storageService.GetDataInStorage(featureKey, context)
	if err != nil {
		serviceContainer.GetLoggerService().Error("ERROR_READING_DATA_FROM_STORAGE", map[string]interface{}{"err": err.Error()}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return nil
	}
	return data
}

// SetDataInStorage stores data in storage with validation
func (sd *StorageDecorator) SetDataInStorage(
	data map[string]interface{},
	storageService interfaces.StorageServiceInterface,
	serviceContainer interfaces.ServiceContainerInterface,
) *campaign.Variation {
	// Extract and validate featureKey
	featureKey, ok := data[enums.StorageFeatureKey.GetValue()].(string)
	if !ok || featureKey == "" {
		serviceContainer.GetLoggerService().Error("ERROR_STORING_DATA_IN_STORAGE", map[string]interface{}{"key": "featureKey"}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return nil
	}

	// Extract and validate userId
	userIDValue, exists := data[enums.StorageUserID.GetValue()]
	if !exists {
		serviceContainer.GetLoggerService().Error("ERROR_STORING_DATA_IN_STORAGE", map[string]interface{}{"key": "Context or Context.id"}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return nil
	}
	userID := fmt.Sprint(userIDValue)
	if userID == "" {
		serviceContainer.GetLoggerService().Error("ERROR_STORING_DATA_IN_STORAGE", map[string]interface{}{"key": "Context or Context.id"}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return nil
	}

	// Extract variation data
	rolloutKey, _ := data[enums.DecisionRolloutKey.GetValue()].(string)
	experimentKey, _ := data[enums.DecisionExperimentKey.GetValue()].(string)
	rolloutVariationID, _ := data[enums.DecisionRolloutVariationID.GetValue()].(int)
	experimentVariationID, _ := data[enums.DecisionExperimentVariationID.GetValue()].(int)

	// Validate rollout data
	if rolloutKey != "" && rolloutVariationID == 0 {
		serviceContainer.GetLoggerService().Error("ERROR_STORING_DATA_IN_STORAGE", map[string]interface{}{"key": "Variation:(rolloutKey or rolloutVariationId)"}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return nil
	}

	// Validate experiment data
	if experimentKey != "" && experimentVariationID == 0 {
		serviceContainer.GetLoggerService().Error("ERROR_STORING_DATA_IN_STORAGE", map[string]interface{}{"key": "Variation:(experimentKey or experimentVariationID)"}, serviceContainer.GetDebuggerService().GetStandardDebugProps())
		return nil
	}

	// Store the data
	storageService.SetDataInStorage(data)

	// Return a new Variation instance
	return &campaign.Variation{}
}
