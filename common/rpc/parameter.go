// Copyright (c) 2017-2022 The Elastos Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import "strconv"

type Parameter map[string]interface{}

func Param(key string, value interface{}) Parameter {
	return Parameter{}.Add(key, value)
}

func (param Parameter) Add(key string, value interface{}) Parameter {
	switch value.(type) {
	case int:
		value = strconv.Itoa(value.(int))
	case int8:
		value = strconv.FormatInt(int64(value.(int8)), 10)
	case int16:
		value = strconv.FormatInt(int64(value.(int16)), 10)
	case int32:
		value = strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		value = strconv.FormatInt(value.(int64), 10)
	case uint:
		value = strconv.FormatUint(uint64(value.(uint)), 10)
	case uint8:
		value = strconv.FormatUint(uint64(value.(uint8)), 10)
	case uint16:
		value = strconv.FormatUint(uint64(value.(uint16)), 10)
	case uint32:
		value = strconv.FormatUint(uint64(value.(uint32)), 10)
	case uint64:
		value = strconv.FormatUint(value.(uint64), 10)
	}
	param[key] = value.(string)
	return param
}
