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

package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/logger/enums"
)

// ANSI color codes for terminal output
const (
	Bold      = "\x1b[1m"
	Cyan      = "\x1b[36m"
	Green     = "\x1b[32m"
	LightBlue = "\x1b[94m"
	Red       = "\x1b[31m"
	Reset     = "\x1b[0m"
	White     = "\x1b[30m"
	Yellow    = "\x1b[33m"
)

// LogMessageBuilder handles formatting of log messages
type LogMessageBuilder struct {
	loggerConfig    map[string]interface{}
	transportConfig map[string]interface{}
	prefix          string
	dateTimeFormat  func() string
}

// NewLogMessageBuilder creates a new LogMessageBuilder instance
func NewLogMessageBuilder(loggerConfig, transportConfig map[string]interface{}) *LogMessageBuilder {
	// Set the prefix, defaulting to an empty string if not provided
	prefix := ""
	if p, ok := transportConfig[enums.LogManagerConfigPrefix.GetValue()].(string); ok && p != "" {
		prefix = p
	} else if p, ok := loggerConfig[enums.LogManagerConfigPrefix.GetValue()].(string); ok && p != "" {
		prefix = p
	}

	// Set the date and time format, defaulting to the logger's format if the transport's format is not provided
	dateTimeFormat := func() string {
		return time.Now().Format(time.RFC3339)
	}

	if dtf, ok := transportConfig[enums.LogManagerConfigDateTimeFormat.GetValue()].(func() string); ok && dtf != nil {
		dateTimeFormat = dtf
	} else if dtf, ok := loggerConfig[enums.LogManagerConfigDateTimeFormat.GetValue()].(func() string); ok && dtf != nil {
		dateTimeFormat = dtf
	}

	return &LogMessageBuilder{
		loggerConfig:    loggerConfig,
		transportConfig: transportConfig,
		prefix:          prefix,
		dateTimeFormat:  dateTimeFormat,
	}
}

// FormatMessage formats a log message combining level, prefix, date/time, and the actual message
func (lmb *LogMessageBuilder) FormatMessage(level string, message string) string {
	return fmt.Sprintf("[%s]: %s %s %s",
		lmb.getFormattedLevel(level),
		lmb.getFormattedPrefix(lmb.prefix),
		lmb.getFormattedDateTime(),
		message)
}

// getFormattedPrefix formats the prefix with ANSI colors if enabled
func (lmb *LogMessageBuilder) getFormattedPrefix(prefix string) string {
	// Default to false if not specified (colors disabled by default)
	isAnsiColorEnabled := false
	if val, ok := lmb.loggerConfig[enums.LogManagerConfigIsAnsiColorEnabled.GetValue()].(bool); ok {
		isAnsiColorEnabled = val
	}

	if isAnsiColorEnabled {
		return fmt.Sprintf("%s%s%s%s", Bold, Green, prefix, Reset)
	}
	return prefix
}

// getFormattedLevel formats the log level with appropriate ANSI colors
func (lmb *LogMessageBuilder) getFormattedLevel(level string) string {
	upperCaseLevel := strings.ToUpper(level)

	// Default to false if not specified (colors disabled by default)
	isAnsiColorEnabled := false
	if val, ok := lmb.loggerConfig[enums.LogManagerConfigIsAnsiColorEnabled.GetValue()].(bool); ok {
		isAnsiColorEnabled = val
	}

	if isAnsiColorEnabled {
		switch enums.ParseLogLevel(level) {
		case enums.LogLevelTrace:
			return fmt.Sprintf("%s%s%s%s", Bold, White, upperCaseLevel, Reset)
		case enums.LogLevelDebug:
			return fmt.Sprintf("%s%s%s%s", Bold, LightBlue, upperCaseLevel, Reset)
		case enums.LogLevelInfo:
			return fmt.Sprintf("%s%s%s%s", Bold, Cyan, upperCaseLevel, Reset)
		case enums.LogLevelWarn:
			return fmt.Sprintf("%s%s%s%s", Bold, Yellow, upperCaseLevel, Reset)
		case enums.LogLevelError:
			return fmt.Sprintf("%s%s%s%s", Bold, Red, upperCaseLevel, Reset)
		default:
			return upperCaseLevel
		}
	}

	return upperCaseLevel
}

// getFormattedDateTime returns the current date and time formatted according to the specified format
func (lmb *LogMessageBuilder) getFormattedDateTime() string {
	return lmb.dateTimeFormat()
}
