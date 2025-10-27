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

// SegmentOperandValue represents different operand types used in segmentation
type SegmentOperandValue int

const (
	SegmentOperandLowerValue SegmentOperandValue = iota + 1
	SegmentOperandStartingEndingStarValue
	SegmentOperandStartingStarValue
	SegmentOperandEndingStarValue
	SegmentOperandRegexValue
	SegmentOperandEqualValue
	SegmentOperandGreaterThanValue
	SegmentOperandGreaterThanEqualToValue
	SegmentOperandLessThanValue
	SegmentOperandLessThanEqualToValue
)

// GetValue returns the integer value of the operand
func (s SegmentOperandValue) GetValue() int {
	return int(s)
}
