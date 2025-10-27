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

package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	log "github.com/wingify/vwo-fme-go-sdk/pkg/log_messages"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models"
	settingsModel "github.com/wingify/vwo-fme-go-sdk/pkg/models/settings"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/manager"
	networkModels "github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"
)

// BatchEventQueue manages the batching of events
type BatchEventQueue struct {
	batchQueue          []map[string]interface{}
	eventsPerRequest    int
	requestTimeInterval int
	ticker              *time.Ticker
	stopChan            chan bool
	isBatchProcessing   bool
	accountID           int
	sdkKey              string
	flushCallback       models.FlushCallback
	logManager          interfaces.LoggerServiceInterface
	settings            *settingsModel.Settings
	mutex               sync.Mutex
	networkManager      *manager.NetworkManager
	isInitialized       bool
}

// NewBatchEventQueue creates a new batch event queue
func NewBatchEventQueue(
	eventsPerRequest int,
	requestTimeInterval int,
	flushCallback models.FlushCallback,
	accountID int,
	sdkKey string,
	logManager interfaces.LoggerServiceInterface,
) *BatchEventQueue {
	queue := &BatchEventQueue{
		batchQueue:          make([]map[string]interface{}, 0),
		eventsPerRequest:    eventsPerRequest,
		requestTimeInterval: requestTimeInterval,
		flushCallback:       flushCallback,
		accountID:           accountID,
		sdkKey:              sdkKey,
		logManager:          logManager,
		stopChan:            make(chan bool, 1),
		isInitialized:       true,
	}

	queue.createNewBatchTimer()
	logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["BATCH_EVENT_QUEUE_INITIALIZED"], map[string]interface{}{
		"eventsPerRequest":    strconv.Itoa(eventsPerRequest),
		"requestTimeInterval": strconv.Itoa(requestTimeInterval),
	}))

	return queue
}

// IsInitialized checks if the batch event queue is initialized
func (batchEventQueue *BatchEventQueue) IsInitialized() bool {
	if batchEventQueue == nil {
		return false
	}
	if batchEventQueue.isInitialized {
		return true
	}
	return false
}

// SetSettings sets the settings for the batch event queue
func (batchEventQueue *BatchEventQueue) SetSettings(settings *settingsModel.Settings) {
	batchEventQueue.settings = settings
}

// Enqueue adds an event to the batch queue
func (batchEventQueue *BatchEventQueue) Enqueue(eventData map[string]interface{}) {
	batchEventQueue.mutex.Lock()
	defer batchEventQueue.mutex.Unlock()

	batchEventQueue.batchQueue = append(batchEventQueue.batchQueue, eventData)
	batchEventQueue.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["EVENT_ADDED_TO_QUEUE"], map[string]interface{}{
		"queueSize": strconv.Itoa(len(batchEventQueue.batchQueue)),
	}))

	// If batch size reaches the limit, trigger flush
	if len(batchEventQueue.batchQueue) >= batchEventQueue.eventsPerRequest {
		batchEventQueue.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["QUEUE_REACHED_MAX_CAPACITY"], nil))
		go batchEventQueue.flush(false)
	}
}

// createNewBatchTimer initializes the batch timer
func (batchEventQueue *BatchEventQueue) createNewBatchTimer() {
	batchEventQueue.ticker = time.NewTicker(time.Duration(batchEventQueue.requestTimeInterval) * time.Second)

	go func() {
		for {
			select {
			case <-batchEventQueue.ticker.C:
				batchEventQueue.flush(false)
			case <-batchEventQueue.stopChan:
				return
			}
		}
	}()

	batchEventQueue.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["BATCH_TIMER_INITIALIZED"], map[string]interface{}{
		"interval": strconv.Itoa(batchEventQueue.requestTimeInterval),
	}))
}

// FlushAndClearInterval flushes the queue and stops the timer
func (batchEventQueue *BatchEventQueue) FlushAndClearInterval() bool {
	if batchEventQueue.ticker != nil {
		batchEventQueue.ticker.Stop()
		select {
		case batchEventQueue.stopChan <- true:
		default:
		}
		batchEventQueue.ticker = nil
		batchEventQueue.stopChan = nil
	}
	return batchEventQueue.flush(true)
}

// flush processes the batch queue
func (batchEventQueue *BatchEventQueue) flush(manual bool) bool {
	batchEventQueue.mutex.Lock()

	if len(batchEventQueue.batchQueue) == 0 {
		batchEventQueue.mutex.Unlock()
		batchEventQueue.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["BATCH_QUEUE_EMPTY"], nil))
		return false
	}

	// Create a snapshot of events to send
	eventsToSend := make([]map[string]interface{}, len(batchEventQueue.batchQueue))
	copy(eventsToSend, batchEventQueue.batchQueue)
	batchEventQueue.batchQueue = make([]map[string]interface{}, 0)

	batchEventQueue.mutex.Unlock()

	batchEventQueue.logManager.Debug(log.BuildMessage(log.DebugLogMessagesEnum["EVENT_BATCH_BEFORE_FLUSHING"], map[string]interface{}{
		"timer": func() string {
			if manual {
				return "Timer will be cleared and registered again"
			} else {
				return ""
			}
		}(),
		"length": strconv.Itoa(len(eventsToSend)),
		"manually": func() string {
			if manual {
				return "manually"
			} else {
				return ""
			}
		}(),
		"accountId": strconv.Itoa(batchEventQueue.accountID),
	}))

	// Send the batch events and handle the result
	isSentSuccessfully := batchEventQueue.sendBatchEvents(eventsToSend)
	if isSentSuccessfully {
		batchEventQueue.logManager.Info(log.BuildMessage(log.InfoLogMessagesEnum["EVENT_BATCH_After_FLUSHING"], map[string]interface{}{
			"length": strconv.Itoa(len(eventsToSend)),
			"manually": func() string {
				if manual {
					return "manually"
				} else {
					return ""
				}
			}(),
		}))
	} else {
		// Re-enqueue events in case of failure for retry logic
		batchEventQueue.mutex.Lock()
		batchEventQueue.batchQueue = append(eventsToSend, batchEventQueue.batchQueue...)
		batchEventQueue.mutex.Unlock()
		batchEventQueue.logManager.Error("BATCH_FLUSH_FAILED", nil, map[string]interface{}{
			"an":        enums.ApiFlushEvents,
			"accountId": strconv.Itoa(batchEventQueue.accountID),
		})
	}

	batchEventQueue.isBatchProcessing = false
	return isSentSuccessfully
}

// sendBatchEvents sends the batch events to VWO server
func (batchEventQueue *BatchEventQueue) sendBatchEvents(events []map[string]interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			eventsJSON, _ := json.Marshal(events)
			if batchEventQueue.flushCallback != nil {
				batchEventQueue.flushCallback(fmt.Sprintf("%v", r), string(eventsJSON))
			}
			batchEventQueue.logManager.Error("ERROR_SENDING_BATCH_EVENTS", map[string]interface{}{"err": fmt.Sprintf("%v", r)}, map[string]interface{}{
				"an":        enums.ApiFlushEvents,
				"accountId": strconv.Itoa(batchEventQueue.accountID),
			})
		}
	}()

	// Send the batch request using the network utility
	isSentSuccessfully := batchEventQueue.SendPostBatchRequest(events, batchEventQueue.accountID, batchEventQueue.sdkKey, batchEventQueue.flushCallback)
	return isSentSuccessfully
}

func (batchEventQueue *BatchEventQueue) SendPostBatchRequest(payload interface{}, accountID int, sdkKey string, flushCallback func(err string, events string)) bool {
	// Create the batch payload
	batchPayload := map[string]interface{}{
		"ev": payload,
	}

	// Create the query parameters
	query := map[string]string{
		"a":   fmt.Sprintf("%d", accountID),
		"env": sdkKey,
	}

	url := constants.HostName
	if batchEventQueue.settings.GetCollectionPrefix() != "" {
		url = url + "/" + batchEventQueue.settings.GetCollectionPrefix()
	}

	// Create the request model
	requestModel := networkModels.NewRequestModel(
		url,
		enums.ApiMethodPost.GetValue(),
		enums.BatchEvents.GetURL(),
		query,
		batchPayload,
		map[string]string{
			"Authorization": sdkKey,
			"Content-Type":  "application/json",
		},
		constants.HTTPSProtocol,
		0,
		"",
	)

	// Send the request
	batchEventQueue.networkManager.Post(requestModel, func(err string, events string) {
		if flushCallback != nil {
			flushCallback(err, events)
		}
	})
	return true
}

// GetBatchQueue returns the current batch queue
func (batchEventQueue *BatchEventQueue) GetBatchQueue() []map[string]interface{} {
	batchEventQueue.mutex.Lock()
	defer batchEventQueue.mutex.Unlock()
	return batchEventQueue.batchQueue
}

func (batchEventQueue *BatchEventQueue) SetNetworkManager(networkManager *manager.NetworkManager) {
	batchEventQueue.networkManager = networkManager
}
