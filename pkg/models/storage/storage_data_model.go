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

package storage

// StorageData represents data to be stored for VWO
type StorageData struct {
	FeatureKey            string `json:"featureKey"`
	FeatureID             int    `json:"featureId"`
	User                  string `json:"user"`
	RolloutID             int    `json:"rolloutId,omitempty"`
	RolloutKey            string `json:"rolloutKey,omitempty"`
	RolloutVariationID    int    `json:"rolloutVariationId,omitempty"`
	ExperimentID          int    `json:"experimentId,omitempty"`
	ExperimentKey         string `json:"experimentKey,omitempty"`
	ExperimentVariationID int    `json:"experimentVariationId,omitempty"`
}

// GetFeatureKey returns the feature key
func (sd *StorageData) GetFeatureKey() string {
	return sd.FeatureKey
}

// GetFeatureID returns the feature ID
func (sd *StorageData) GetFeatureID() int {
	return sd.FeatureID
}

// GetUser returns the user
func (sd *StorageData) GetUser() string {
	return sd.User
}

// GetRolloutID returns the rollout ID
func (sd *StorageData) GetRolloutID() int {
	return sd.RolloutID
}

// GetRolloutKey returns the rollout key
func (sd *StorageData) GetRolloutKey() string {
	return sd.RolloutKey
}

// GetRolloutVariationID returns the rollout variation ID
func (sd *StorageData) GetRolloutVariationID() int {
	return sd.RolloutVariationID
}

// GetExperimentID returns the experiment ID
func (sd *StorageData) GetExperimentID() int {
	return sd.ExperimentID
}

// GetExperimentKey returns the experiment key
func (sd *StorageData) GetExperimentKey() string {
	return sd.ExperimentKey
}

// GetExperimentVariationID returns the experiment variation ID
func (sd *StorageData) GetExperimentVariationID() int {
	return sd.ExperimentVariationID
}

// SetFeatureKey sets the feature key
func (sd *StorageData) SetFeatureKey(featureKey string) {
	sd.FeatureKey = featureKey
}

// SetFeatureID sets the feature ID
func (sd *StorageData) SetFeatureID(featureID int) {
	sd.FeatureID = featureID
}

// SetUser sets the user
func (sd *StorageData) SetUser(user string) {
	sd.User = user
}

// SetRolloutID sets the rollout ID
func (sd *StorageData) SetRolloutID(rolloutID int) {
	sd.RolloutID = rolloutID
}

// SetRolloutKey sets the rollout key
func (sd *StorageData) SetRolloutKey(rolloutKey string) {
	sd.RolloutKey = rolloutKey
}

// SetRolloutVariationID sets the rollout variation ID
func (sd *StorageData) SetRolloutVariationID(rolloutVariationID int) {
	sd.RolloutVariationID = rolloutVariationID
}

// SetExperimentID sets the experiment ID
func (sd *StorageData) SetExperimentID(experimentID int) {
	sd.ExperimentID = experimentID
}

// SetExperimentKey sets the experiment key
func (sd *StorageData) SetExperimentKey(experimentKey string) {
	sd.ExperimentKey = experimentKey
}

// SetExperimentVariationID sets the experiment variation ID
func (sd *StorageData) SetExperimentVariationID(experimentVariationID int) {
	sd.ExperimentVariationID = experimentVariationID
}
