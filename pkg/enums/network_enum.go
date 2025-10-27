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

// NetworkEnum represents network request keys
type NetworkEnum string

const (
	NetworkHostname   NetworkEnum = "hostname"
	NetworkAgent      NetworkEnum = "agent"
	NetworkScheme     NetworkEnum = "scheme"
	NetworkProtocol   NetworkEnum = "protocol"
	NetworkPort       NetworkEnum = "port"
	NetworkHeaders    NetworkEnum = "headers"
	NetworkMethod     NetworkEnum = "method"
	NetworkBody       NetworkEnum = "body"
	NetworkPath       NetworkEnum = "path"
	NetworkTimeout    NetworkEnum = "timeout"
	NetworkAPIVersion NetworkEnum = "api-version"
	NetworkSDKName    NetworkEnum = "sn"
	NetworkSDKVersion NetworkEnum = "sv"
	NetworkURL        NetworkEnum = "url"
)

// GetValue returns the string value of the network enum
func (n NetworkEnum) GetValue() string {
	return string(n)
}
