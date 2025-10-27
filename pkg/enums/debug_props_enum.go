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

// DebugPropsEnum represents debug properties keys
type DebugPropsEnum string

const (
	DebugPropAPI         DebugPropsEnum = "an"
	DebugPropFeatureKey  DebugPropsEnum = "fk"
	DebugPropUUID        DebugPropsEnum = "uuid"
	DebugPropSessionID   DebugPropsEnum = "sId"
	DebugPropMessage     DebugPropsEnum = "msg"
	DebugPropMessageType DebugPropsEnum = "msg_t"
	DebugPropLogLevel    DebugPropsEnum = "lt"
	DebugPropCategory    DebugPropsEnum = "cg"
	DebugPropEventID     DebugPropsEnum = "eventId"
	DebugPropAccountID   DebugPropsEnum = "a"
	DebugPropProduct     DebugPropsEnum = "product"
	DebugPropSDKName     DebugPropsEnum = "sn"
	DebugPropSDKVersion  DebugPropsEnum = "sv"
)

// GetValue returns the string value of the debug props enum
func (d DebugPropsEnum) GetValue() string {
	return string(d)
}
