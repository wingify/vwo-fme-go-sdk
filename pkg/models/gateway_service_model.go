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

package models

// GatewayService represents the gateway service configuration
type GatewayService struct {
	URL  string `json:"url"`
	Port int    `json:"port,omitempty"`
}

// GetURL returns the gateway service URL
func (gs *GatewayService) GetURL() string {
	return gs.URL
}

// GetPort returns the gateway service port
func (gs *GatewayService) GetPort() int {
	return gs.Port
}
