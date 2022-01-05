package config

var (
	regnet = ConfigFile{
		ConfigFile: Configuration{
			Version: 0,
			MainNode: &MainNodeConfig{
				Rpc: &RpcConfig{
					User:         "",
					Pass:         "",
					HttpJsonPort: 22336,
					IpAddress:    "127.0.0.1",
				},
			},
		},
	}

	testnet = ConfigFile{
		ConfigFile: Configuration{
			Version: 0,
			MainNode: &MainNodeConfig{
				Rpc: &RpcConfig{
					User:         "",
					Pass:         "",
					HttpJsonPort: 21336,
					IpAddress:    "127.0.0.1",
				},
			},
		},
	}

	mainnet = ConfigFile{
		ConfigFile: Configuration{
			ActiveNet: "mainnet",
			Version:   0,
			MainNode: &MainNodeConfig{
				Rpc: &RpcConfig{
					User:         "",
					Pass:         "",
					HttpJsonPort: 20336,
					IpAddress:    "127.0.0.1",
				},
			},
		},
	}
)
