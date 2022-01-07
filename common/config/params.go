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
