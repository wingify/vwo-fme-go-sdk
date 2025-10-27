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

package user

import (
	"time"

	"github.com/wingify/vwo-fme-go-sdk/pkg/enums"
)

// VWOContext represents the user context for VWO operations
type VWOContext struct {
	ID                          string                 `json:"id"`
	UserAgent                   string                 `json:"userAgent,omitempty"`
	IPAddress                   string                 `json:"ipAddress,omitempty"`
	CustomVariables             map[string]interface{} `json:"customVariables,omitempty"`
	VariationTargetingVariables map[string]interface{} `json:"variationTargetingVariables,omitempty"`
	PostSegmentationVariables   []string               `json:"postSegmentationVariables,omitempty"`
	VWO                         *ContextVWO            `json:"_vwo,omitempty"`
	SessionId                   int64                  `json:"sessionId,omitempty"`
	UUID                        string                 `json:"uuid,omitempty"`
}

// NewVWOContext creates a new VWOContext from a map
func NewVWOContext(context map[string]interface{}) *VWOContext {
	vwoContext := &VWOContext{}

	if id, ok := context[enums.ContextID.GetValue()]; ok {
		if idStr, ok := id.(string); ok {
			vwoContext.ID = idStr
		}
	}

	if userAgent, ok := context[enums.ContextUserAgent.GetValue()]; ok {
		if userAgentStr, ok := userAgent.(string); ok {
			vwoContext.UserAgent = userAgentStr
		}
	}

	if ipAddress, ok := context[enums.ContextIPAddress.GetValue()]; ok {
		if ipAddressStr, ok := ipAddress.(string); ok {
			vwoContext.IPAddress = ipAddressStr
		}
	}

	if customVariables, ok := context[enums.ContextCustomVariables.GetValue()]; ok {
		if customVarsMap, ok := customVariables.(map[string]interface{}); ok {
			vwoContext.CustomVariables = customVarsMap
		}
	}

	if variationTargetingVariables, ok := context[enums.ContextVariationTargetingVariables.GetValue()]; ok {
		if vtVarsMap, ok := variationTargetingVariables.(map[string]interface{}); ok {
			vwoContext.VariationTargetingVariables = vtVarsMap
		}
	}

	if postSegmentationVariables, ok := context[enums.ContextPostSegmentationVariables.GetValue()]; ok {
		if psvArray, ok := postSegmentationVariables.([]string); ok {
			vwoContext.PostSegmentationVariables = psvArray
		}
	}

	if vwo, ok := context[enums.ContextVWO.GetValue()]; ok {
		if vwoMap, ok := vwo.(map[string]interface{}); ok {
			vwoContext.VWO = NewContextVWO(vwoMap)
		}
	}

	if sessionId, ok := context[enums.ContextSessionID.GetValue()]; ok {
		if sessionIdInt, ok := sessionId.(int64); ok {
			vwoContext.SessionId = sessionIdInt
		}
	} else {
		vwoContext.SessionId = time.Now().Unix()
	}

	return vwoContext
}

// GetID returns the user ID
func (c *VWOContext) GetID() string {
	return c.ID
}

// GetUserAgent returns the user agent
func (c *VWOContext) GetUserAgent() string {
	return c.UserAgent
}

// GetIPAddress returns the IP address
func (c *VWOContext) GetIPAddress() string {
	return c.IPAddress
}

// GetCustomVariables returns custom variables
func (c *VWOContext) GetCustomVariables() map[string]interface{} {
	return c.CustomVariables
}

// SetCustomVariables sets custom variables
func (c *VWOContext) SetCustomVariables(customVariables map[string]interface{}) {
	c.CustomVariables = customVariables
}

// GetVariationTargetingVariables returns variation targeting variables
func (c *VWOContext) GetVariationTargetingVariables() map[string]interface{} {
	return c.VariationTargetingVariables
}

// SetVariationTargetingVariables sets variation targeting variables
func (c *VWOContext) SetVariationTargetingVariables(variationTargetingVariables map[string]interface{}) {
	c.VariationTargetingVariables = variationTargetingVariables
}

// GetVWO returns the VWO context
func (c *VWOContext) GetVWO() *ContextVWO {
	return c.VWO
}

// SetVWO sets the VWO context
func (c *VWOContext) SetVWO(vwo *ContextVWO) {
	c.VWO = vwo
}

// GetPostSegmentationVariables returns post segmentation variables
func (c *VWOContext) GetPostSegmentationVariables() []string {
	return c.PostSegmentationVariables
}

// SetPostSegmentationVariables sets post segmentation variables
func (c *VWOContext) SetPostSegmentationVariables(postSegmentationVariables []string) {
	c.PostSegmentationVariables = postSegmentationVariables
}

// GetSessionId returns the session ID
func (c *VWOContext) GetSessionId() int64 {
	return c.SessionId
}

// GetUUID returns the UUID
func (c *VWOContext) GetUUID() string {
	return c.UUID
}

// SetSessionId sets the session ID
func (c *VWOContext) SetSessionId(sessionId int64) {
	c.SessionId = sessionId
}

// SetUUID sets the UUID
func (c *VWOContext) SetUUID(uuid string) {
	c.UUID = uuid
}
