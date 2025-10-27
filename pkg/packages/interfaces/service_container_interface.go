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

package interfaces

import (
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/campaign"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/manager"
)

// ServiceContainerInterface defines the contract for service container
// This interface is used to break import cycles between core and segmentation_evaluator packages
type ServiceContainerInterface interface {
	GetLoggerService() LoggerServiceInterface
	GetSettings() *settings.Settings
	GetBaseUrl() string
	GetVWOInitOptions() *models.VWOInitOptions
	GetSettingsManager() SettingsManagerInterface
	GetHooksManager() HooksManagerInterface
	GetSegmentationManager() SegmentationManagerInterface
	GetDebuggerService() DebuggerServiceInterface
	GetNetworkManager() *manager.NetworkManager
	GetBatchEventQueue() BatchEventQueueInterface
}

// SegmentationManagerInterface defines the contract for segmentation manager
// This interface is used to break import cycles between core and segmentation_evaluator packages
type SegmentationManagerInterface interface {
	SetContextualData(serviceContainer ServiceContainerInterface, feature *campaign.Feature, context *user.VWOContext)
	ValidateSegmentation(dsl interface{}, properties map[string]interface{}) bool
}
