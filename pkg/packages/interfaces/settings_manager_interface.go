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
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/manager"
)

// SettingsManagerInterface defines the contract for settings manager
// This interface is used to break import cycles and provide abstraction for settings management
type SettingsManagerInterface interface {
	// Getter methods
	GetAccountID() string
	GetSDKKey() string
	GetProtocol() string
	GetHostname() string
	GetPort() int
	GetSettingsObject() *settings.Settings
	GetIsGatewayServiceProvided() bool
	GetLoggerService() LoggerServiceInterface
	GetNetworkManager() *manager.NetworkManager

	// Additional getter methods for direct field access
	GetSettingsFetchTime() int64
	GetIsSettingsProvidedInInit() bool
	GetStartTimeForInit() int64
	GetIsSettingsValidOnInit() bool

	// Settings management methods
	GetSettings(forceFetch bool) string
	FetchSettings(isViaWebhook bool) (string, error)
	SetSettings(settingsObj *settings.Settings, settingsString string)
	SetSettingsValidOnInit(valid bool)
	SetNetworkManager(networkManager *manager.NetworkManager)
}
