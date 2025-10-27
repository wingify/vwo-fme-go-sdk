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

package core

import (
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/manager"
	segmentationCore "github.com/wingify/vwo-fme-go-sdk/pkg/packages/segmentation_evaluator/core"
	"github.com/wingify/vwo-fme-go-sdk/pkg/services"
)

// ServiceContainer holds all the service instances
// Ensure ServiceContainer implements ServiceContainerInterface
var _ interfaces.ServiceContainerInterface = (*ServiceContainer)(nil)

type ServiceContainer struct {
	userId              string
	loggerService       interfaces.LoggerServiceInterface
	settingsManager     *services.SettingsManager
	hooksManager        *services.HooksManager
	options             *models.VWOInitOptions
	batchEventQueue     *services.BatchEventQueue
	segmentationManager interfaces.SegmentationManagerInterface
	settings            *settings.Settings
	debuggerService     *services.DebuggerService
	networkManager      *manager.NetworkManager
}

// NewServiceContainer creates a new ServiceContainer instance
func NewServiceContainer(
	userId string,
	loggerService interfaces.LoggerServiceInterface,
	settingsManager *services.SettingsManager,
	options *models.VWOInitOptions,
	batchEventQueue *services.BatchEventQueue,
	settings *settings.Settings,
	networkManager *manager.NetworkManager,
) *ServiceContainer {
	// Create segmentation manager internally
	segmentationManager := segmentationCore.NewSegmentationManager(loggerService)
	debuggerService := services.NewDebuggerService()
	return &ServiceContainer{
		userId:              userId,
		loggerService:       loggerService,
		settingsManager:     settingsManager,
		hooksManager:        services.NewHooksManager(getIntegrationsCallback(options.GetIntegrations())),
		options:             options,
		batchEventQueue:     batchEventQueue,
		segmentationManager: segmentationManager,
		settings:            settings,
		debuggerService:     debuggerService,
		networkManager:      networkManager,
	}
}

// GetUserId returns the user ID
func (sc *ServiceContainer) GetUserId() string {
	return sc.userId
}

// GetLoggerService returns the logger service instance
func (sc *ServiceContainer) GetLoggerService() interfaces.LoggerServiceInterface {
	return sc.loggerService
}

// GetSettingsManager returns the settings manager instance
func (sc *ServiceContainer) GetSettingsManager() interfaces.SettingsManagerInterface {
	return sc.settingsManager
}

// GetHooksManager returns the hooks manager instance
func (sc *ServiceContainer) GetHooksManager() interfaces.HooksManagerInterface {
	return sc.hooksManager
}

// GetVWOInitOptions returns the VWO init options instance
func (sc *ServiceContainer) GetVWOInitOptions() *models.VWOInitOptions {
	return sc.options
}

// GetBatchEventQueue returns the batch event queue instance
func (sc *ServiceContainer) GetBatchEventQueue() interfaces.BatchEventQueueInterface {
	return sc.batchEventQueue
}

// GetSegmentationManager returns the segmentation manager instance
func (sc *ServiceContainer) GetSegmentationManager() interfaces.SegmentationManagerInterface {
	return sc.segmentationManager
}

// GetSettings returns the settings instance
func (sc *ServiceContainer) GetSettings() *settings.Settings {
	return sc.settings
}

// GetBaseUrl returns the base URL for API requests
func (sc *ServiceContainer) GetBaseUrl() string {
	baseUrl := sc.settingsManager.Hostname

	if sc.settingsManager.IsGatewayServiceProvided {
		return baseUrl
	}

	if sc.settings.GetCollectionPrefix() != "" {
		return baseUrl + "/" + sc.settings.GetCollectionPrefix()
	}

	return baseUrl
}

// getIntegrationsCallback extracts the callback function from IntegrationOptions
func getIntegrationsCallback(integrations *models.IntegrationOptions) func(map[string]interface{}) {
	if integrations != nil {
		return integrations.Callback
	}
	return nil
}

// GetDebuggerService returns the debugger service instance
func (sc *ServiceContainer) GetDebuggerService() interfaces.DebuggerServiceInterface {
	return sc.debuggerService
}

// GetNetworkManager returns the network manager instance
func (sc *ServiceContainer) GetNetworkManager() *manager.NetworkManager {
	return sc.networkManager
}
