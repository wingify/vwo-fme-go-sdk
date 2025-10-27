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
	"math"
	"reflect"
	"time"
)

// IsObject checks if value is an object (map or struct)
func IsObject(val interface{}) bool {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.Map || kind == reflect.Struct
}

// IsArray checks if value is an array or slice
func IsArray(val interface{}) bool {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

// IsNull checks if value is nil
func IsNull(val interface{}) bool {
	return val == nil
}

// IsUndefined checks if value is undefined (nil in Go)
func IsUndefined(val interface{}) bool {
	return val == nil
}

// IsDefined checks if value is defined (not nil)
func IsDefined(val interface{}) bool {
	return val != nil
}

// IsNumber checks if value is a number
func IsNumber(val interface{}) bool {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind >= reflect.Int && kind <= reflect.Float64
}

// IsInteger checks if value is an integer
func IsInteger(val interface{}) bool {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind >= reflect.Int && kind <= reflect.Uint64
}

// IsString checks if value is a string
func IsString(val interface{}) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.String
}

// IsBoolean checks if value is a boolean
func IsBoolean(val interface{}) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Bool
}

// IsNaN checks if value is NaN (Not a Number)
func IsNaN(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case float32:
		return math.IsNaN(float64(v))
	case float64:
		return math.IsNaN(v)
	default:
		return false
	}
}

// IsDate checks if value is a time.Time
func IsDate(val interface{}) bool {
	if val == nil {
		return false
	}
	_, ok := val.(time.Time)
	return ok
}

// IsFunction checks if value is a function
func IsFunction(val interface{}) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Func
}

// GetType returns the type of the value as a string
func GetType(val interface{}) string {
	if IsNull(val) {
		return "Null"
	}
	if IsUndefined(val) {
		return "Undefined"
	}
	if IsNaN(val) {
		return "NaN"
	}
	if IsArray(val) {
		return "Array"
	}
	if IsObject(val) {
		return "Object"
	}
	if IsInteger(val) {
		return "Integer"
	}
	if IsNumber(val) {
		return "Number"
	}
	if IsString(val) {
		return "String"
	}
	if IsBoolean(val) {
		return "Boolean"
	}
	if IsDate(val) {
		return "Date"
	}
	if IsFunction(val) {
		return "Function"
	}
	return "Unknown Type"
}
