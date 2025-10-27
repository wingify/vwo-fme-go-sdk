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

package constants

import (
	"time"
)

const (
	// Platform identifier
	Platform = "server"

	// Traffic and campaign constants
	MaxTrafficPercent = 100
	MaxTrafficValue   = 10000
	StatusRunning     = "RUNNING"

	// Seed and event constants
	SeedValue                  = 1
	MaxEventsPerRequest        = 5000
	DefaultRequestTimeInterval = 600 // 10 * 60(secs) = 600 secs i.e. 10 minutes
	DefaultEventsPerRequest    = 100

	// SDK information
	SDKName    = "vwo-fme-go-sdk"
	SDKVersion = "1.3.0"

	// Settings constants
	SettingsExpiry  = 10000000
	SettingsTimeout = 50000

	// Network constants
	NetworkTimeout = 30 * time.Second

	// API endpoints
	HostName                = "dev.visualwebsiteoptimizer.com"
	SettingsEndpoint        = "/server-side/v2-settings"
	WebhookSettingsEndpoint = "/server-side/v2-pull"

	// Environment and protocol
	VWOFsEnvironment = "vwo_fs_environment"
	HTTPSProtocol    = "https"
	HTTPProtocol     = "http"

	// Algorithm and meta constants
	RandomAlgo                  = 1
	VWOMetaMegKey               = "_vwo_meta_meg_"
	VariationTargetingUserIDKey = "_vwoUserId"

	// Polling constants
	DefaultPollInterval = 600000 // 10 minutes
	FME                 = "fme"

	// Gateway Service constants
	BaseURL                = "" // Will be set from settings
	EndpointAttributeCheck = "/check-attribute"
	EndpointGetUserData    = "/get-user-details"
	QueryParamUserAgent    = "userAgent"
	QueryParamIPAddress    = "ipAddress"

	// decision maker constants
	MAX_TRAFFIC_VALUE  = 10000
	MAX_CAMPAIGN_VALUE = 100

	// Debugger constants
	NETWORK_CALL_EXCEPTION                 = "NETWORK_CALL_EXCEPTION"
	FLAG_DECISION                          = "FLAG_DECISION"
	POLLING                                = "polling"
	NETWORK_CALL_SUCCESS_WITH_RETRIES      = "NETWORK_CALL_SUCCESS_WITH_RETRIES"
	NETWORK_CALL_FAILURE_AFTER_MAX_RETRIES = "NETWORK_CALL_FAILURE_AFTER_MAX_RETRIES"
	IMPACT_ANALYSIS                        = "IMPACT_ANALYSIS"
)

// NonRetryable events
var NonRetryableEvents = []string{"vwo_sdkDebug", "vwo_sdkUsageStats", "vwo_fmeSdkInit"}
