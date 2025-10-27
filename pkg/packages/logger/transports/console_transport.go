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

package transports

import (
	"fmt"
	"os"

	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/enums"
)

// ConsoleTransport implements Logger interface to provide logging functionality
type ConsoleTransport struct {
	config map[string]interface{}
	level  string
}

// NewConsoleTransport creates a new ConsoleTransport instance
func NewConsoleTransport(config map[string]interface{}) *ConsoleTransport {
	level := enums.LogLevelError.String()
	if l, ok := config[enums.LogManagerConfigLevel.GetValue()].(string); ok && l != "" {
		level = l
	}

	return &ConsoleTransport{
		config: config,
		level:  level,
	}
}

// Trace logs a trace message
func (ct *ConsoleTransport) Trace(message string) {
	ct.consoleLog(enums.LogLevelTrace.String(), message)
}

// Debug logs a debug message
func (ct *ConsoleTransport) Debug(message string) {
	ct.consoleLog(enums.LogLevelDebug.String(), message)
}

// Info logs an info message
func (ct *ConsoleTransport) Info(message string) {
	ct.consoleLog(enums.LogLevelInfo.String(), message)
}

// Warn logs a warning message
func (ct *ConsoleTransport) Warn(message string) {
	ct.consoleLog(enums.LogLevelWarn.String(), message)
}

// Error logs an error message
func (ct *ConsoleTransport) Error(message string) {
	ct.consoleLog(enums.LogLevelError.String(), message)
}

// Log implements the LogTransport interface
func (ct *ConsoleTransport) Log(level enums.LogLevel, message string) {
	ct.consoleLog(level.String(), message)
}

// consoleLog logs messages to the console based on the log level
func (ct *ConsoleTransport) consoleLog(level, message string) {
	// Check if the message should be logged based on the configured level
	if !ct.shouldLog(level) {
		return
	}

	// Use appropriate console output based on level
	switch level {
	case enums.LogLevelError.String():
		fmt.Fprintln(os.Stderr, message)
	case enums.LogLevelWarn.String():
		fmt.Fprintln(os.Stdout, message)
	default:
		fmt.Fprintln(os.Stdout, message)
	}
}

// shouldLog determines if a message should be logged based on the configured level
func (ct *ConsoleTransport) shouldLog(level string) bool {
	configLevel := enums.ParseLogLevel(ct.level)
	messageLevel := enums.ParseLogLevel(level)
	return messageLevel.GetLevel() >= configLevel.GetLevel()
}
