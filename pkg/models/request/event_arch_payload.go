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
	"encoding/json"
)

// EventArchPayload represents the payload structure for event arch APIs
type EventArchPayload struct {
	D *EventArchData `json:"d"`
}

// EventArchData represents the data structure within the event arch payload
type EventArchData struct {
	MsgID     string   `json:"msgId"`
	VisID     string   `json:"visId"`
	SessionID int64    `json:"sessionId"`
	Event     *Event   `json:"event"`
	Visitor   *Visitor `json:"visitor"`
	VisitorUA string   `json:"visitor_ua,omitempty"`
	VisitorIP string   `json:"visitor_ip,omitempty"`
}

// Event represents an event in the event arch payload
type Event struct {
	Props *Props `json:"props"`
	Name  string `json:"name"`
	Time  int64  `json:"time"`
}

// Props represents properties of an event
type Props struct {
	SDKName              string                 `json:"vwo_sdkName"`
	SDKVersion           string                 `json:"vwo_sdkVersion"`
	EnvKey               string                 `json:"vwo_envKey"`
	Variation            string                 `json:"variation,omitempty"`
	ID                   int                    `json:"id,omitempty"`
	IsFirst              int                    `json:"isFirst,omitempty"`
	IsCustomEvent        bool                   `json:"isCustomEvent,omitempty"`
	VWOMeta              map[string]interface{} `json:"vwoMeta,omitempty"`
	Product              string                 `json:"product,omitempty"`
	Data                 map[string]interface{} `json:"data,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for Props
func (p *Props) MarshalJSON() ([]byte, error) {
	// Create a map to hold all properties
	result := make(map[string]interface{})

	// Add all the standard fields
	if p.SDKName != "" {
		result["vwo_sdkName"] = p.SDKName
	}
	if p.SDKVersion != "" {
		result["vwo_sdkVersion"] = p.SDKVersion
	}
	if p.EnvKey != "" {
		result["vwo_envKey"] = p.EnvKey
	}
	if p.Variation != "" {
		result["variation"] = p.Variation
	}
	if p.ID != 0 {
		result["id"] = p.ID
	}
	if p.IsFirst != 0 {
		result["isFirst"] = p.IsFirst
	}
	if p.IsCustomEvent {
		result["isCustomEvent"] = p.IsCustomEvent
	}
	if len(p.VWOMeta) > 0 {
		result["vwoMeta"] = p.VWOMeta
	}
	if p.Product != "" {
		result["product"] = p.Product
	}
	if len(p.Data) > 0 {
		result["data"] = p.Data
	}

	// Add additional properties at the same level
	if p.AdditionalProperties != nil {
		for key, value := range p.AdditionalProperties {
			result[key] = value
		}
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for Props
func (p *Props) UnmarshalJSON(data []byte) error {
	// Create a temporary map to hold all JSON data
	temp := make(map[string]interface{})
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Initialize AdditionalProperties if nil
	if p.AdditionalProperties == nil {
		p.AdditionalProperties = make(map[string]interface{})
	}

	// Process each field
	for key, value := range temp {
		switch key {
		case "vwo_sdkName":
			if str, ok := value.(string); ok {
				p.SDKName = str
			}
		case "vwo_sdkVersion":
			if str, ok := value.(string); ok {
				p.SDKVersion = str
			}
		case "vwo_envKey":
			if str, ok := value.(string); ok {
				p.EnvKey = str
			}
		case "variation":
			if str, ok := value.(string); ok {
				p.Variation = str
			}
		case "id":
			if num, ok := value.(float64); ok {
				p.ID = int(num)
			}
		case "isFirst":
			if num, ok := value.(float64); ok {
				p.IsFirst = int(num)
			}
		case "isCustomEvent":
			if b, ok := value.(bool); ok {
				p.IsCustomEvent = b
			}
		case "vwoMeta":
			if m, ok := value.(map[string]interface{}); ok {
				p.VWOMeta = m
			}
		case "product":
			if str, ok := value.(string); ok {
				p.Product = str
			}
		case "data":
			if m, ok := value.(map[string]interface{}); ok {
				p.Data = m
			}
		default:
			// Add unknown fields to AdditionalProperties
			p.AdditionalProperties[key] = value
		}
	}

	return nil
}

// Visitor represents visitor information in the event arch payload
type Visitor struct {
	Props map[string]interface{} `json:"props"`
}
