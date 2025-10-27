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

// HTTPMethodEnum represents HTTP method constants
type HTTPMethodEnum string

const (
	HTTPMethodGET  HTTPMethodEnum = "GET"
	HTTPMethodPOST HTTPMethodEnum = "POST"
)

// GetValue returns the string value of the HTTP method enum
func (h HTTPMethodEnum) GetValue() string {
	return string(h)
}

// NetworkOptionsEnum represents network options keys
type NetworkOptionsEnum string

const (
	NetworkOptionsHostname NetworkOptionsEnum = "hostname"
	NetworkOptionsPath     NetworkOptionsEnum = "path"
	NetworkOptionsScheme   NetworkOptionsEnum = "scheme"
	NetworkOptionsPort     NetworkOptionsEnum = "port"
	NetworkOptionsHeaders  NetworkOptionsEnum = "headers"
	NetworkOptionsBody     NetworkOptionsEnum = "body"
)

// GetValue returns the string value of the network options enum
func (n NetworkOptionsEnum) GetValue() string {
	return string(n)
}

// HTTPHeaderEnum represents HTTP header constants
type HTTPHeaderEnum string

const (
	HTTPHeaderContentType HTTPHeaderEnum = "Content-Type"
)

// GetValue returns the string value of the HTTP header enum
func (h HTTPHeaderEnum) GetValue() string {
	return string(h)
}

// ContentTypeEnum represents content type constants
type ContentTypeEnum string

const (
	ContentTypeApplicationJSON ContentTypeEnum = "application/json"
)

// GetValue returns the string value of the content type enum
func (c ContentTypeEnum) GetValue() string {
	return string(c)
}

// NetworkClientErrorEnum represents network client error message constants
type NetworkClientErrorEnum string

const (
	NetworkClientErrorFailedToCreateGETRequest  NetworkClientErrorEnum = "failed to create GET request"
	NetworkClientErrorGETRequestFailed          NetworkClientErrorEnum = "GET request failed"
	NetworkClientErrorFailedToReadResponse      NetworkClientErrorEnum = "failed to read response body"
	NetworkClientErrorInvalidResponse           NetworkClientErrorEnum = "invalid response %s, Status Code: %d, Response : %s"
	NetworkClientErrorFailedToMarshalBody       NetworkClientErrorEnum = "failed to marshal request body"
	NetworkClientErrorFailedToCreatePOSTRequest NetworkClientErrorEnum = "failed to create POST request"
	NetworkClientErrorPOSTRequestFailed         NetworkClientErrorEnum = "POST request failed"
	NetworkClientErrorRequestFailed             NetworkClientErrorEnum = "request failed. Status Code: %d, Response: %s"
)

// GetValue returns the string value of the network client error enum
func (n NetworkClientErrorEnum) GetValue() string {
	return string(n)
}
