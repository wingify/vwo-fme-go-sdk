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

package enums

// SegmentOperatorValue represents the different operators used in segmentation
type SegmentOperatorValue string

const (
	SegmentOperatorAND             SegmentOperatorValue = "and"
	SegmentOperatorNOT             SegmentOperatorValue = "not"
	SegmentOperatorOR              SegmentOperatorValue = "or"
	SegmentOperatorCustomVariable  SegmentOperatorValue = "custom_variable"
	SegmentOperatorUser            SegmentOperatorValue = "user"
	SegmentOperatorCountry         SegmentOperatorValue = "country"
	SegmentOperatorRegion          SegmentOperatorValue = "region"
	SegmentOperatorCity            SegmentOperatorValue = "city"
	SegmentOperatorOperatingSystem SegmentOperatorValue = "os"
	SegmentOperatorDeviceType      SegmentOperatorValue = "device_type"
	SegmentOperatorBrowserAgent    SegmentOperatorValue = "browser_string"
	SegmentOperatorUA              SegmentOperatorValue = "ua"
	SegmentOperatorDevice          SegmentOperatorValue = "device"
	SegmentOperatorFeatureID       SegmentOperatorValue = "featureId"
	SegmentOperatorIP              SegmentOperatorValue = "ip_address"
	SegmentOperatorBrowserVersion  SegmentOperatorValue = "browser_version"
	SegmentOperatorOSVersion       SegmentOperatorValue = "os_version"
)

// String returns the string value of the operator
func (s SegmentOperatorValue) String() string {
	return string(s)
}

// FromValue converts a string to SegmentOperatorValue
func SegmentOperatorValueFromString(value string) (SegmentOperatorValue, bool) {
	operators := map[string]SegmentOperatorValue{
		"and":             SegmentOperatorAND,
		"not":             SegmentOperatorNOT,
		"or":              SegmentOperatorOR,
		"custom_variable": SegmentOperatorCustomVariable,
		"user":            SegmentOperatorUser,
		"country":         SegmentOperatorCountry,
		"region":          SegmentOperatorRegion,
		"city":            SegmentOperatorCity,
		"os":              SegmentOperatorOperatingSystem,
		"device_type":     SegmentOperatorDeviceType,
		"browser_string":  SegmentOperatorBrowserAgent,
		"ua":              SegmentOperatorUA,
		"device":          SegmentOperatorDevice,
		"featureId":       SegmentOperatorFeatureID,
		"ip_address":      SegmentOperatorIP,
		"browser_version": SegmentOperatorBrowserVersion,
		"os_version":      SegmentOperatorOSVersion,
	}

	op, exists := operators[value]
	return op, exists
}
