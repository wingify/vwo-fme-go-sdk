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

package log_messages

import (
	_ "embed"
	"encoding/json"
)

//go:embed debug-messages.json
var debugMessagesJSON []byte

//go:embed info-messages.json
var infoMessagesJSON []byte

//go:embed error-messages.json
var errorMessagesJSON []byte

//go:embed trace-messages.json
var traceMessagesJSON []byte

//go:embed warn-messages.json
var warnMessagesJSON []byte

// DebugLogMessagesEnum contains all debug log message templates
var DebugLogMessagesEnum map[string]string

// InfoLogMessagesEnum contains all info log message templates
var InfoLogMessagesEnum map[string]string

// ErrorLogMessagesEnum contains all error log message templates
var ErrorLogMessagesEnum map[string]string

// TraceLogMessagesEnum contains all trace log message templates
var TraceLogMessagesEnum map[string]string

// WarnLogMessagesEnum contains all warn log message templates
var WarnLogMessagesEnum map[string]string

func init() {
	// Load debug messages
	DebugLogMessagesEnum = make(map[string]string)
	if err := json.Unmarshal(debugMessagesJSON, &DebugLogMessagesEnum); err != nil {
		DebugLogMessagesEnum = make(map[string]string)
	}

	// Load info messages
	InfoLogMessagesEnum = make(map[string]string)
	if err := json.Unmarshal(infoMessagesJSON, &InfoLogMessagesEnum); err != nil {
		InfoLogMessagesEnum = make(map[string]string)
	}

	// Load error messages
	ErrorLogMessagesEnum = make(map[string]string)
	if err := json.Unmarshal(errorMessagesJSON, &ErrorLogMessagesEnum); err != nil {
		ErrorLogMessagesEnum = make(map[string]string)
	}

	// Load trace messages
	TraceLogMessagesEnum = make(map[string]string)
	if err := json.Unmarshal(traceMessagesJSON, &TraceLogMessagesEnum); err != nil {
		TraceLogMessagesEnum = make(map[string]string)
	}

	// Load warn messages
	WarnLogMessagesEnum = make(map[string]string)
	if err := json.Unmarshal(warnMessagesJSON, &WarnLogMessagesEnum); err != nil {
		WarnLogMessagesEnum = make(map[string]string)
	}
}
