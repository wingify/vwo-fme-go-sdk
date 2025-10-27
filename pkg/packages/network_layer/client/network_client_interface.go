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

package client

import "github.com/wingify/vwo-fme-go-sdk/pkg/packages/network_layer/models"

// NetworkClientInterface defines the interface for network operations
type NetworkClientInterface interface {
	// GET sends a GET request to the server
	GET(request *models.RequestModel) *models.ResponseModel

	// POST sends a POST request to the server
	POST(request *models.RequestModel) *models.ResponseModel
}
