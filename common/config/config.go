package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// DefaultConfigFilename indicates the file name of config.
	DefaultConfigFilename = "./config.json"
)

var (
	Version    string
	Parameters ConfigParams
)

type RpcConfiguration struct {
	User        string   `json:"User"`
	Pass        string   `json:"Pass"`
	WhiteIPList []string `json:"WhiteIPList"`
}

type Configuration struct {
	ActiveNet string          `json:"ActiveNet"`
	Version   uint32          `json:"Version"`
	MainNode  *MainNodeConfig `json:"MainNode"`
}

type RpcConfig struct {
	IpAddress    string `json:"IpAddress"`
	HttpJsonPort int    `json:"HttpJsonPort"`
	User         string `json:"User"`
	Pass         string `json:"Pass"`
}

type MainNodeConfig struct {
	Rpc *RpcConfig `json:"Rpc"`
}

type ConfigFile struct {
	ConfigFile Configuration `json:"Configuration"`
}

type ConfigParams struct {
	*Configuration
}

func Initialize() {
	file, e := ioutil.ReadFile(DefaultConfigFilename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		return
	}
	i := ConfigFile{}
	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	e = json.Unmarshal(file, &i)
	var config ConfigFile
	switch strings.ToLower(i.ConfigFile.ActiveNet) {
	case "testnet", "test":
		config = testnet
	case "regnet", "reg":
		config = regnet
	default:
		config = mainnet
	}

	Parameters.Configuration = &(config.ConfigFile)

	e = json.Unmarshal(file, &config)
	if e != nil {
		fmt.Printf("Unmarshal json file erro %v", e)
		os.Exit(1)
	}

	e = json.Unmarshal(file, &config)
	if e != nil {
		fmt.Printf("Unmarshal json file erro %v", e)
		os.Exit(1)
	}

	var out bytes.Buffer
	err := json.Indent(&out, file, "", "")
	if err != nil {
		fmt.Printf("Config file error: %v\n", e)
		os.Exit(1)
	}

	if Parameters.Configuration.MainNode == nil {
		fmt.Printf("Need to set main node in config file\n")
		return
	}
}
