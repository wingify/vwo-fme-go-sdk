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

package utils

import (
	"net/url"
	"strings"

	"github.com/wingify/vwo-fme-go-sdk/pkg/constants"
	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
	"github.com/wingify/vwo-fme-go-sdk/pkg/models/user"
	"github.com/wingify/vwo-fme-go-sdk/pkg/packages/interfaces"
)

// CreateAndSendImpressionForVariationShown creates and sends an impression for a variation shown event.
// This function constructs the necessary properties and payload for the event
// and uses the NetworkUtil to send a POST API request.
func CreateAndSendImpressionForVariationShown(
	serviceContainer interfaces.ServiceContainerInterface,
	campaignID int,
	variationID int,
	context *user.VWOContext,
	featureKey string,
) {
	// Get base properties for the event
	properties := GetEventsBaseProperties(
		serviceContainer.GetSettingsManager(),
		enums.VWOVariationShown.GetValue(),
		EncodeURIComponent(context.GetUserAgent()),
		context.GetIPAddress(),
	)

	// Construct payload data for tracking the user
	payload := GetTrackUserPayloadData(
		serviceContainer,
		enums.VWOVariationShown.GetValue(),
		campaignID,
		variationID,
		context,
	)

	// get campaign key and variation name
	campaignKeyWithVariationName := GetCampaignKeyFromCampaignID(serviceContainer.GetSettings(), campaignID)
	variationName := GetVariationNameFromCampaignIdAndVariationId(serviceContainer.GetSettings(), campaignID, variationID)

	campaignKey := campaignKeyWithVariationName // default campaign key

	if campaignKeyWithVariationName == featureKey {
		campaignKey = constants.IMPACT_ANALYSIS
	} else {
		// split campaignKeyWithVariationName by featureKey_ and get the first part
		campaignKey = strings.Split(campaignKeyWithVariationName, featureKey+"_")[1]
	}

	// get campaign type
	campaignType := GetCampaignTypeFromCampaignId(serviceContainer.GetSettings(), campaignID)

	// Check if batch event queue is available
	if serviceContainer.GetBatchEventQueue().IsInitialized() {
		// Enqueue the event to the batch queue for future processing
		serviceContainer.GetBatchEventQueue().Enqueue(payload)
	} else {
		// Send the event immediately if batch event queue is not available
		SendPostAPIRequest(serviceContainer, properties, payload, context, map[string]interface{}{
			"campaignKey":   campaignKey,
			"variationName": variationName,
			"featureKey":    featureKey,
			"campaignType":  campaignType,
		})
	}
}

// EncodeURIComponent encodes a string to be URL-safe
func EncodeURIComponent(value string) string {
	return url.QueryEscape(value)
}
