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
	"reflect"

	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/enums"
)

// LogTransport defines the interface for log transport
type LogTransport interface {
	Log(level enums.LogLevel, message string)
}

// LogTransportManager manages logging transports and delegates logging messages to them
type LogTransportManager struct {
	transports []map[string]interface{}
	config     map[string]interface{}
}

// NewLogTransportManager creates a new LogTransportManager instance
func NewLogTransportManager(config map[string]interface{}) *LogTransportManager {
	return &LogTransportManager{
		transports: make([]map[string]interface{}, 0),
		config:     config,
	}
}

// AddTransport adds a new transport to the manager
func (ltm *LogTransportManager) AddTransport(transport map[string]interface{}) {
	ltm.transports = append(ltm.transports, transport)
}

// ShouldLog determines if the log should be processed based on the transport and configuration levels
func (ltm *LogTransportManager) ShouldLog(transportLevel, configLevel string) bool {
	// Default to the most specific level available
	targetLevel := enums.ParseLogLevel(transportLevel)
	desiredLevel := enums.ParseLogLevel(configLevel)

	// If configLevel is empty, use the config level
	if configLevel == "" {
		if l, ok := ltm.config[enums.LogManagerConfigLevel.GetValue()].(string); ok && l != "" {
			desiredLevel = enums.ParseLogLevel(l)
		} else {
			desiredLevel = enums.LogLevelError
		}
	}

	return targetLevel.GetLevel() >= desiredLevel.GetLevel()
}

// Trace logs a message at TRACE level
func (ltm *LogTransportManager) Trace(message string) {
	ltm.Log(enums.LogLevelTrace, message)
}

// Debug logs a message at DEBUG level
func (ltm *LogTransportManager) Debug(message string) {
	ltm.Log(enums.LogLevelDebug, message)
}

// Info logs a message at INFO level
func (ltm *LogTransportManager) Info(message string) {
	ltm.Log(enums.LogLevelInfo, message)
}

// Warn logs a message at WARN level
func (ltm *LogTransportManager) Warn(message string) {
	ltm.Log(enums.LogLevelWarn, message)
}

// Error logs a message at ERROR level
func (ltm *LogTransportManager) Error(message string) {
	ltm.Log(enums.LogLevelError, message)
}

// Log delegates the logging of messages to the appropriate transports
func (ltm *LogTransportManager) Log(level enums.LogLevel, message string) {
	for _, transport := range ltm.transports {
		logMessageBuilder := logger.NewLogMessageBuilder(ltm.config, transport)
		formattedMessage := logMessageBuilder.FormatMessage(level.String(), message)

		// Get transport level
		transportLevel := level.String()
		if tLevel, ok := transport[enums.LogManagerConfigLevel.GetValue()].(string); ok && tLevel != "" {
			transportLevel = tLevel
		}

		if ltm.ShouldLog(level.String(), transportLevel) {
			// Check if custom log handler is available
			if logHandler, ok := transport["log"]; ok && logHandler != nil {
				// Check if it's a function
				if reflect.TypeOf(logHandler).Kind() == reflect.Func {
					// Call the function with level and message
					args := []reflect.Value{
						reflect.ValueOf(level.String()),
						reflect.ValueOf(message),
					}
					reflect.ValueOf(logHandler).Call(args)
				} else if logTransport, ok := logHandler.(LogTransport); ok {
					// Use LogTransport interface
					logTransport.Log(level, message)
				}
			} else {
				// Use the default log method based on level
				if loggerImpl, ok := transport[level.String()]; ok && loggerImpl != nil {
					if reflect.TypeOf(loggerImpl).Kind() == reflect.Func {
						args := []reflect.Value{reflect.ValueOf(formattedMessage)}
						reflect.ValueOf(loggerImpl).Call(args)
					}
				}
			}
		}
	}
}
