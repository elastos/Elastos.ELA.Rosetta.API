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

package config

var (
	regnet = ConfigFile{
		ConfigFile: Configuration{
			ActiveNet:  "regnet",
			Version:    0,
			ServerPort: 1234,
			MainNodeRPC: &RpcConfig{
				User:         "",
				Pass:         "",
				HttpJsonPort: 22336,
				IpAddress:    "127.0.0.1",
			},
		},
	}

	testnet = ConfigFile{
		ConfigFile: Configuration{
			ActiveNet:  "testnet",
			Version:    0,
			ServerPort: 1234,
			MainNodeRPC: &RpcConfig{
				User:         "",
				Pass:         "",
				HttpJsonPort: 21336,
				IpAddress:    "127.0.0.1",
			},
		},
	}

	mainnet = ConfigFile{
		ConfigFile: Configuration{
			ActiveNet:  "mainnet",
			Version:    0,
			ServerPort: 1234,
			MainNodeRPC: &RpcConfig{
				User:         "",
				Pass:         "",
				HttpJsonPort: 20336,
				IpAddress:    "127.0.0.1",
			},
		},
	}
)
