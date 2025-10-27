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

// BatchFlushCallback is a function type for batch flush callbacks
type BatchFlushCallback func(err error, events string)

// BatchEventData represents batch event configuration
type BatchEventData struct {
	EventsPerRequest    int                `json:"eventsPerRequest"`
	RequestTimeInterval int                `json:"requestTimeInterval"` // in seconds
	FlushCallback       BatchFlushCallback `json:"flushCallback"`
}

// NewBatchEventData creates a new BatchEventData with default values
func NewBatchEventData(options map[string]interface{}) *BatchEventData {
	batchEventData := &BatchEventData{}

	if eventsPerRequest, ok := options[enums.BatchEventEventsPerRequest.GetValue()].(int); ok {
		batchEventData.EventsPerRequest = eventsPerRequest
	}

	if requestTimeInterval, ok := options[enums.BatchEventRequestTimeInterval.GetValue()].(int); ok {
		batchEventData.RequestTimeInterval = requestTimeInterval
	}

	if flushCallback, ok := options[enums.BatchEventFlushCallback.GetValue()].(func(err error, events string)); ok {
		batchEventData.FlushCallback = flushCallback
	}

	return batchEventData
}

// GetEventsPerRequest returns the events per request
func (b *BatchEventData) GetEventsPerRequest() int {
	return b.EventsPerRequest
}

// SetEventsPerRequest sets the events per request
func (b *BatchEventData) SetEventsPerRequest(value int) {
	b.EventsPerRequest = value
}

// GetRequestTimeInterval returns the request time interval
func (b *BatchEventData) GetRequestTimeInterval() int {
	return b.RequestTimeInterval
}

// SetRequestTimeInterval sets the request time interval
func (b *BatchEventData) SetRequestTimeInterval(value int) {
	b.RequestTimeInterval = value
}

// GetFlushCallback returns the flush callback
func (b *BatchEventData) GetFlushCallback() BatchFlushCallback {
	return b.FlushCallback
}

// SetFlushCallback sets the flush callback
func (b *BatchEventData) SetFlushCallback(callback BatchFlushCallback) {
	b.FlushCallback = callback
}
