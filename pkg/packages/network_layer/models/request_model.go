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

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
)

// RequestModel represents an HTTP request with all necessary configuration
type RequestModel struct {
	URL       string
	Method    string
	Scheme    string
	Port      int
	Path      string
	Query     map[string]string
	Timeout   int
	Body      map[string]interface{}
	Headers   map[string]string
	EventName string
}

// NewRequestModel creates a new RequestModel with the specified parameters
func NewRequestModel(
	baseURL string,
	method string,
	path string,
	query map[string]string,
	body map[string]interface{},
	headers map[string]string,
	scheme string,
	port int,
	eventName string,
) *RequestModel {
	if method == "" {
		method = "GET"
	}
	if scheme == "" {
		scheme = "http"
	}

	return &RequestModel{
		URL:       baseURL,
		Method:    method,
		Path:      path,
		Query:     query,
		Body:      body,
		Headers:   headers,
		Scheme:    scheme,
		Port:      port,
		Timeout:   -1, // Default timeout not set
		EventName: eventName,
	}
}

// GetOptions returns the network options as a map
func (r *RequestModel) GetOptions() map[string]interface{} {
	options := make(map[string]interface{})

	// Build query parameters
	var queryParams strings.Builder
	if r.Query != nil {
		for key, value := range r.Query {
			if queryParams.Len() > 0 {
				queryParams.WriteString("&")
			}
			queryParams.WriteString(fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
		}
	}

	// Set hostname
	options[enums.NetworkHostname.GetValue()] = r.URL
	options[enums.NetworkAgent.GetValue()] = false

	// Set scheme
	if r.Scheme != "" {
		options[enums.NetworkScheme.GetValue()] = r.Scheme
	}

	// Set port if not default
	if r.Port != 0 && r.Port != 80 {
		options[enums.NetworkPort.GetValue()] = r.Port
	}

	// Set headers
	if r.Headers != nil {
		options[enums.NetworkHeaders.GetValue()] = r.Headers
	}

	// Set method
	if r.Method != "" {
		options[enums.NetworkMethod.GetValue()] = r.Method
	}

	// Set body and related headers
	if r.Body != nil {
		bodyJSON, err := json.Marshal(r.Body)
		if err == nil {
			if r.Headers == nil {
				r.Headers = make(map[string]string)
			}
			r.Headers["Content-Type"] = "application/json"
			r.Headers["Content-Length"] = fmt.Sprintf("%d", len(bodyJSON))
			options[enums.NetworkHeaders.GetValue()] = r.Headers
			options[enums.NetworkBody.GetValue()] = r.Body
		}
	}

	// Set path with query parameters
	if r.Path != "" {
		combinedPath := r.Path
		if queryParams.Len() > 0 {
			combinedPath += "?" + queryParams.String()
		}
		options[enums.NetworkPath.GetValue()] = combinedPath
	}

	// Set timeout
	if r.Timeout > 0 {
		options[enums.NetworkTimeout.GetValue()] = r.Timeout
	}

	return options
}

func (r *RequestModel) GetEventName() string {
	return r.EventName
}
