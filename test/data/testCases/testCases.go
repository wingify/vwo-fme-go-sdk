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

package testCases

// TestCases represents the collection of all test cases
type TestCases struct {
	GetFlagWithoutStorage []TestData `json:"GETFLAG_WITHOUT_STORAGE"`
	GetFlagWithSalt       []TestData `json:"GETFLAG_WITH_SALT"`
	GetFlagMegRandom      []TestData `json:"GETFLAG_MEG_RANDOM"`
	GetFlagMegAdvance     []TestData `json:"GETFLAG_MEG_ADVANCE"`
	GetFlagWithStorage    []TestData `json:"GETFLAG_WITH_STORAGE"`
}
