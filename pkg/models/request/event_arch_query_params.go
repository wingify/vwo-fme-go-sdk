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

package request

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
)

// SettingsQueryParams represents query parameters for settings API
type SettingsQueryParams struct {
	I string // API key
	R string // Random value
	A string // Account ID
}

// NewSettingsQueryParams creates a new SettingsQueryParams instance
func NewSettingsQueryParams(apiKey string, random string, accountID string) *SettingsQueryParams {
	return &SettingsQueryParams{
		I: apiKey,
		R: random,
		A: accountID,
	}
}

// GetQueryParams returns the query parameters as a map
func (s *SettingsQueryParams) GetQueryParams() map[string]string {
	return map[string]string{
		"i": s.I,
		"r": s.R,
		"a": s.A,
	}
}

// RequestQueryParams represents query parameters for event arch requests
type RequestQueryParams struct {
	EN        string  // Event name
	A         string  // Account ID
	Env       string  // SDK key
	ETime     int64   // Event time
	Random    float64 // Random value
	P         string  // Platform
	VisitorUA string  // Visitor user agent
	VisitorIP string  // Visitor IP address
	SN        string  // SDK name
	SV        string  // SDK version
}

// NewRequestQueryParams creates a new RequestQueryParams instance
func NewRequestQueryParams(eventName string, accountID string, sdkKey string, visitorUserAgent string, ipAddress string) *RequestQueryParams {
	return &RequestQueryParams{
		EN:        eventName,
		A:         accountID,
		Env:       sdkKey,
		ETime:     time.Now().UnixNano() / 1e6,
		Random:    rand.Float64(),
		P:         "FS",
		VisitorUA: visitorUserAgent,
		VisitorIP: ipAddress,
		SN:        constants.SDKName,
		SV:        constants.SDKVersion,
	}
}

// GetQueryParams returns the query parameters as a map
func (r *RequestQueryParams) GetQueryParams() map[string]string {
	params := map[string]string{
		"en":     r.EN,
		"a":      r.A,
		"env":    r.Env,
		"eTime":  fmt.Sprintf("%d", r.ETime),
		"random": fmt.Sprintf("%.16f", r.Random),
		"p":      r.P,
		"sn":     r.SN,
		"sv":     r.SV,
	}

	if r.VisitorUA != "" {
		params["visitor_ua"] = r.VisitorUA
	}
	if r.VisitorIP != "" {
		params["visitor_ip"] = r.VisitorIP
	}

	return params
}
