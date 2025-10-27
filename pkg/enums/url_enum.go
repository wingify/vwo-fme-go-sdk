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

// UrlEnum represents URL endpoints
type UrlEnum string

const (
	Events         UrlEnum = "/events/t"
	AttributeCheck UrlEnum = "/check-attribute"
	GetUserData    UrlEnum = "/get-user-details"
	BatchEvents    UrlEnum = "/server-side/batch-events-v2"
)

// GetURL returns the URL string value
func (u UrlEnum) GetURL() string {
	return string(u)
}
