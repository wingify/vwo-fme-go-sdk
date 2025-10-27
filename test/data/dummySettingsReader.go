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

package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DummySettingsReader reads settings from JSON files
type DummySettingsReader struct {
	SettingsMap map[string]string
}

// NewDummySettingsReader creates a new DummySettingsReader and loads settings
func NewDummySettingsReader() *DummySettingsReader {
	reader := &DummySettingsReader{
		SettingsMap: make(map[string]string),
	}
	reader.loadSettings()
	return reader
}

// loadSettings loads all settings files from the settings directory
func (r *DummySettingsReader) loadSettings() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Try different possible paths for the settings directory
	possiblePaths := []string{
		filepath.Join(wd, "test/data/settings"),
		filepath.Join(wd, "data/settings"),
		"test/data/settings",
		"data/settings",
		// If we're in a subdirectory, go up to find the correct path
		filepath.Join(wd, "..", "data", "settings"),
		filepath.Join(wd, "..", "..", "test", "data", "settings"),
	}

	var settingsDir string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			settingsDir = path
			break
		}
	}

	if settingsDir == "" {
		panic(fmt.Sprintf("Could not find settings directory. Tried paths: %v", possiblePaths))
	}

	files, err := os.ReadDir(settingsDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			settingsName := file.Name()[:len(file.Name())-5] // Remove .json extension
			filePath := filepath.Join(settingsDir, file.Name())

			data, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}

			// Validate JSON
			var jsonData interface{}
			if err := json.Unmarshal(data, &jsonData); err != nil {
				panic(err)
			}

			r.SettingsMap[settingsName] = string(data)
		}
	}
}
