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
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
)

// RetryConfig represents the retry configuration for network requests
type RetryConfig struct {
	ShouldRetry       bool `json:"shouldRetry"`
	MaxRetries        int  `json:"maxRetries"`
	InitialDelay      int  `json:"initialDelay"` // in seconds
	BackoffMultiplier int  `json:"backoffMultiplier"`
}

// NewRetryConfig creates a new RetryConfig with default values
func NewRetryConfig() *RetryConfig {
	return &RetryConfig{
		ShouldRetry:       true,
		MaxRetries:        3,
		InitialDelay:      2,
		BackoffMultiplier: 2,
	}
}

// NewRetryConfigFromMap creates a RetryConfig from a map of options
func NewRetryConfigFromMap(options map[string]interface{}) *RetryConfig {
	config := NewRetryConfig()

	if shouldRetry, ok := options[enums.RetryConfigShouldRetry.GetValue()].(bool); ok {
		config.ShouldRetry = shouldRetry
	}

	if maxRetries, ok := options[enums.RetryConfigMaxRetries.GetValue()].(int); ok {
		config.MaxRetries = maxRetries
	}

	if initialDelay, ok := options[enums.RetryConfigInitialDelay.GetValue()].(int); ok {
		config.InitialDelay = initialDelay
	}

	if backoffMultiplier, ok := options[enums.RetryConfigBackoffMultiplier.GetValue()].(int); ok {
		config.BackoffMultiplier = backoffMultiplier
	}

	return config
}

// GetRetryDelay calculates the delay for a specific retry attempt using exponential backoff
func (r *RetryConfig) GetRetryDelay(attempt int) int {
	if attempt <= 0 {
		return 0
	}

	// Calculate delay: initialDelay * (backoffMultiplier ^ (attempt - 1))
	delay := r.InitialDelay
	for i := 1; i < attempt; i++ {
		delay *= r.BackoffMultiplier
	}

	return delay
}

// IsRetryable checks if a retry should be attempted based on the current attempt number
func (r *RetryConfig) IsRetryable(attempt int) bool {
	return r.ShouldRetry && attempt <= r.MaxRetries
}
