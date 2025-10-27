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
	"os"
	"path/filepath"

	"github.com/wingify/vwo-fme-go-sdk/test/data/testCases"
)

// TestDataReader reads test cases from JSON files
type TestDataReader struct {
	TestCases *testCases.TestCases
}

// NewTestDataReader creates a new TestDataReader and loads test cases
func NewTestDataReader() *TestDataReader {
	reader := &TestDataReader{}
	reader.TestCases = reader.readTestCases("test/data/testCases")
	return reader
}

// readTestCases reads the test cases from a JSON file located in the specified folder
func (r *TestDataReader) readTestCases(folderPath string) *testCases.TestCases {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Try different possible paths for the test cases file
	possiblePaths := []string{
		filepath.Join(wd, "test/data/testCases", "index.json"),
		filepath.Join(wd, "data/testCases", "index.json"),
		"test/data/testCases/index.json",
		"data/testCases/index.json",
		filepath.Join(folderPath, "index.json"),
		// If we're in a subdirectory, go up to find the correct path
		filepath.Join(wd, "..", "data", "testCases", "index.json"),
		filepath.Join(wd, "..", "..", "test", "data", "testCases", "index.json"),
	}

	var indexPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			indexPath = path
			break
		}
	}

	if indexPath == "" {
		return nil
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		panic(err)
	}

	var testCases testCases.TestCases
	if err := json.Unmarshal(data, &testCases); err != nil {
		panic(err)
	}

	return &testCases
}
