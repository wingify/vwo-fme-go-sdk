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

package enums

// LogManagerConfigEnum represents log manager configuration keys
type LogManagerConfigEnum string

const (
	LogManagerConfigName                LogManagerConfigEnum = "name"
	LogManagerConfigLevel               LogManagerConfigEnum = "level"
	LogManagerConfigPrefix              LogManagerConfigEnum = "prefix"
	LogManagerConfigDateTimeFormat      LogManagerConfigEnum = "dateTimeFormat"
	LogManagerConfigIsAnsiColorEnabled  LogManagerConfigEnum = "isAnsiColorEnabled"
	LogManagerConfigIsAlwaysNewInstance LogManagerConfigEnum = "isAlwaysNewInstance"
	LogManagerConfigTransports          LogManagerConfigEnum = "transports"
	LogManagerConfigTransport           LogManagerConfigEnum = "transport"
	LogManagerConfigRequestID           LogManagerConfigEnum = "requestId"
)

// GetValue returns the string value of the log manager config enum
func (l LogManagerConfigEnum) GetValue() string {
	return string(l)
}

// LogManagerDefaultsEnum represents log manager default values
type LogManagerDefaultsEnum string

const (
	LogManagerDefaultName   LogManagerDefaultsEnum = "VWO Logger"
	LogManagerDefaultPrefix LogManagerDefaultsEnum = "VWO-SDK"
)

// GetValue returns the string value of the log manager defaults enum
func (l LogManagerDefaultsEnum) GetValue() string {
	return string(l)
}

// LogManagerTransportEnum represents log manager transport method names
type LogManagerTransportEnum string

const (
	LogManagerTransportTrace LogManagerTransportEnum = "trace"
	LogManagerTransportDebug LogManagerTransportEnum = "debug"
	LogManagerTransportInfo  LogManagerTransportEnum = "info"
	LogManagerTransportWarn  LogManagerTransportEnum = "warn"
	LogManagerTransportError LogManagerTransportEnum = "error"
)

// GetValue returns the string value of the log manager transport enum
func (l LogManagerTransportEnum) GetValue() string {
	return string(l)
}

// LogManagerDateTimeFormatEnum represents date time format constants
type LogManagerDateTimeFormatEnum string

const (
	LogManagerDateTimeFormat LogManagerDateTimeFormatEnum = "2006-01-02T15:04:05.000Z07:00"
)

// GetValue returns the string value of the log manager date time format enum
func (l LogManagerDateTimeFormatEnum) GetValue() string {
	return string(l)
}
