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

package vwo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	settingsModel "github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
)

// Init initializes the VWO FME client with the provided options
func Init(options map[string]interface{}) (vwoClientInstance *VWOClient, err error) {
	// handle panic and return error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to initialize VWO FME client: %v", r)
		}
	}()

	// start time for init
	startTimeForInit := time.Now().UnixNano() / 1e6
	// Validate required parameters
	if options[enums.OptionSDKKey.GetValue()] == nil || options[enums.OptionSDKKey.GetValue()] == "" {
		return nil, fmt.Errorf(log.ErrorLogMessagesEnum["INVALID_SDK_KEY_IN_OPTIONS"])
	}

	if options[enums.OptionAccountID.GetValue()] == nil || options[enums.OptionAccountID.GetValue()] == 0 {
		return nil, fmt.Errorf(log.ErrorLogMessagesEnum["INVALID_ACCOUNT_ID_IN_OPTIONS"])
	}

	// Convert map to VWOInitOptions using the factory function
	vwoOptions := models.NewVWOInitOptions(options)

	if vwoOptions == nil {
		return nil, fmt.Errorf(log.ErrorLogMessagesEnum["INVALID_OPTIONS"])
	}

	// Create builder and setup services
	builder := &vwoBuilder{
		options: vwoOptions,
	}

	builder.SetLogger().
		SetSettingsManager().
		SetNetworkManager().
		SetStorage().
		InitBatching().
		InitPolling()

	// Check if settings were provided in options
	builder.settingsManager.StartTimeForInit = startTimeForInit
	if vwoOptions.Settings != "" {
		// Parse and validate the provided settings
		builder.originalSettings = vwoOptions.Settings
		builder.settingsManager.IsSettingsProvidedInInit = true

		// Parse the settings JSON string into Settings object
		var settingsObj settingsModel.Settings
		err := json.Unmarshal([]byte(vwoOptions.Settings), &settingsObj)
		if err != nil {
			return nil, fmt.Errorf("failed to parse provided settings: %v", err)
		}

		// Set the parsed settings
		builder.settings = &settingsObj
		builder.settingsManager.SetSettings(&settingsObj, vwoOptions.Settings)

		vwoClientInstance = builder.Build(builder.settings)
	} else {
		// Fetch settings from server
		builder.GetSettings(false)
		vwoClientInstance = builder.Build(builder.settings)
	}

	return vwoClientInstance, nil
}
