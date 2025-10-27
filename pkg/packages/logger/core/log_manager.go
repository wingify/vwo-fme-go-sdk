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
	"fmt"
	"time"

	"github.com/google/uuid"
	globalEnums "github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log_messages "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/transports"
	"github.com/wingify/vwo-fme-go-sdk/pkg/utils"
)

// LogManager provides logging functionality with support for multiple transports
type LogManager struct {
	transportManager *LogTransportManager
	config           map[string]interface{}
	name             string
	requestID        string
	level            string
	prefix           string
	dateTimeFormat   func() string
	transports       []map[string]interface{}
	settingsManager  interfaces.SettingsManagerInterface
}

// NewLogManager creates a new LogManager instance
func NewLogManager(config map[string]interface{}) *LogManager {
	// Initialize configuration with defaults or provided values
	if config == nil {
		config = make(map[string]interface{})
	}

	name := enums.LogManagerDefaultName.GetValue()
	if n, ok := config[enums.LogManagerConfigName.GetValue()].(string); ok && n != "" {
		name = n
	}

	level := enums.LogLevelError.String()
	if l, ok := config[enums.LogManagerConfigLevel.GetValue()].(string); ok && l != "" {
		level = l
	}

	prefix := enums.LogManagerDefaultPrefix.GetValue()
	if p, ok := config[enums.LogManagerConfigPrefix.GetValue()].(string); ok && p != "" {
		prefix = p
	}

	dateTimeFormat := func() string {
		now := time.Now()
		return now.Format(enums.LogManagerDateTimeFormat.GetValue())
	}
	if dtf, ok := config[enums.LogManagerConfigDateTimeFormat.GetValue()].(func() string); ok && dtf != nil {
		dateTimeFormat = dtf
	}

	lm := &LogManager{
		transportManager: NewLogTransportManager(config),
		config:           config,
		name:             name,
		requestID:        uuid.New().String(),
		level:            level,
		prefix:           prefix,
		dateTimeFormat:   dateTimeFormat,
		transports:       make([]map[string]interface{}, 0),
	}

	// Update config with actual values
	config[enums.LogManagerConfigName.GetValue()] = lm.name
	config[enums.LogManagerConfigRequestID.GetValue()] = lm.requestID
	config[enums.LogManagerConfigLevel.GetValue()] = lm.level
	config[enums.LogManagerConfigPrefix.GetValue()] = lm.prefix
	config[enums.LogManagerConfigDateTimeFormat.GetValue()] = lm.dateTimeFormat

	lm.handleTransports()
	return lm
}

func (lm *LogManager) SetSettingsManager(settingsManager interfaces.SettingsManagerInterface) {
	lm.settingsManager = settingsManager
}

// handleTransports handles the initialization and setup of transports based on configuration
func (lm *LogManager) handleTransports() {
	// Check for transports array
	if transportsList, ok := lm.config[enums.LogManagerConfigTransports.GetValue()].([]map[string]interface{}); ok && len(transportsList) > 0 {
		lm.AddTransports(transportsList)
	} else if transport, ok := lm.config[enums.LogManagerConfigTransport.GetValue()].(map[string]interface{}); ok && len(transport) > 0 {
		lm.AddTransport(transport)
	} else {
		// Add default ConsoleTransport if no other transport is specified
		defaultTransport := transports.NewConsoleTransport(map[string]interface{}{
			enums.LogManagerConfigLevel.GetValue(): lm.config[enums.LogManagerConfigLevel.GetValue()],
		})
		lm.AddTransport(map[string]interface{}{
			enums.LogManagerTransportTrace.GetValue(): defaultTransport.Trace,
			enums.LogManagerTransportDebug.GetValue(): defaultTransport.Debug,
			enums.LogManagerTransportInfo.GetValue():  defaultTransport.Info,
			enums.LogManagerTransportWarn.GetValue():  defaultTransport.Warn,
			enums.LogManagerTransportError.GetValue(): defaultTransport.Error,
		})
	}
}

// AddTransport adds a single transport to the LogManager
func (lm *LogManager) AddTransport(transport map[string]interface{}) {
	lm.transportManager.AddTransport(transport)
	lm.transports = append(lm.transports, transport)
}

// AddTransports adds multiple transports to the LogManager
func (lm *LogManager) AddTransports(transports []map[string]interface{}) {
	for _, transport := range transports {
		lm.AddTransport(transport)
	}
}

// Trace logs a trace message
func (lm *LogManager) Trace(message string) {
	lm.transportManager.Trace(message)
}

// Debug logs a debug message
func (lm *LogManager) Debug(message string) {
	lm.transportManager.Debug(message)
}

// Info logs an informational message
func (lm *LogManager) Info(message string) {
	lm.transportManager.Info(message)
}

// Warn logs a warning message
func (lm *LogManager) Warn(message string) {
	lm.transportManager.Warn(message)
}

// Error logs an error message
func (lm *LogManager) Error(template string, templateData map[string]interface{}, extraData map[string]interface{}, shouldSendDebugEvent ...bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error sending debug event: ", r)
		}
	}()
	message := log_messages.BuildMessage(log_messages.ErrorLogMessagesEnum[template], templateData)
	lm.transportManager.Error(message)

	// create debug event props
	debugEventProps := make(map[string]interface{})
	// debugEventProps should contain all extraData
	for key, value := range extraData {
		debugEventProps[key] = value
	}
	debugEventProps[globalEnums.DebugPropMessageType.GetValue()] = template
	debugEventProps[globalEnums.DebugPropMessage.GetValue()] = message
	debugEventProps[globalEnums.DebugPropLogLevel.GetValue()] = enums.LogLevelError.String()
	debugEventProps[globalEnums.DebugPropCategory.GetValue()] = enums.LogLevelError.String()

	if len(shouldSendDebugEvent) == 0 {
		utils.SendDebugEventToVWO(lm.settingsManager, debugEventProps)
	}
}

// Log delegates the logging of messages to the appropriate transports
func (lm *LogManager) Log(level enums.LogLevel, message string) {
	lm.transportManager.Log(level, message)
}
