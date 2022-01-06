package rpc

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/servers"
)

type Response struct {
	ID      int64       `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	*Error  `json:"error"`
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type ArbitratorGroupInfo struct {
	OnDutyArbitratorIndex int
	Arbitrators           []string
}

func GetCurrentHeight(config *config.RpcConfig) (uint32, error) {
	result, err := CallAndUnmarshal("getblockcount", nil, config)
	if err != nil {
		return 0, err
	}
	if count, ok := result.(float64); ok && count >= 1 {
		return uint32(count) - 1, nil
	}
	return 0, errors.New("[GetCurrentHeight] invalid count")
}

func GetBlockByHeight(height uint32, config *config.RpcConfig) (*base.BlockInfo, error) {
	resp, err := CallAndUnmarshal("getblockbyheight", Param("height", height), config)
	if err != nil {
		return nil, err
	}
	block := &base.BlockInfo{}
	Unmarshal(&resp, block)

	return block, nil
}

func GetBlockByHash(hash *common.Uint256, config *config.RpcConfig) (*base.BlockInfo, error) {
	hashBytes, err := common.HexStringToBytes(hash.String())
	if err != nil {
		return nil, err
	}
	reversedHashBytes := common.BytesReverse(hashBytes)
	reversedHashStr := common.BytesToHexString(reversedHashBytes)

	resp, err := CallAndUnmarshal("getblock",
		Param("blockhash", reversedHashStr).Add("verbosity", 2), config)
	if err != nil {
		return nil, err
	}
	block := &base.BlockInfo{}
	if err := Unmarshal(&resp, block); err != nil {
		return nil, err
	}

	return block, nil
}

func GetReferenceAddress(txid string, index int, config *config.RpcConfig) (string, error) {
	parameter := make(map[string]interface{})
	parameter["txid"] = txid
	parameter["index"] = index
	result, err := CallAndUnmarshal("getreferenceaddress", parameter, config)
	if err != nil {
		return "", err
	}
	if a, ok := result.(string); ok {
		return a, nil
	}
	return "", errors.New("invalid data type")
}

func GetAmountByInputs(inputs []*types.Input, config *config.RpcConfig) (common.Fixed64, error) {
	buf := new(bytes.Buffer)
	if err := common.WriteVarUint(buf, uint64(len(inputs))); err != nil {
		return 0, err
	}
	for _, input := range inputs {
		if err := input.Serialize(buf); err != nil {
			return 0, err
		}
	}
	parameter := make(map[string]interface{})
	parameter["inputs"] = common.BytesToHexString(buf.Bytes())
	result, err := CallAndUnmarshal("getamountbyinputs", parameter, config)
	if err != nil {
		return 0, err
	}
	if a, ok := result.(string); ok {
		amount, err := common.StringToFixed64(a)
		if err != nil {
			return 0, err
		}
		return *amount, nil
	}
	return 0, errors.New("get amount by inputs failed")
}

func GetUnspentUtxo(addresses []string, config *config.RpcConfig) ([]base.UTXOInfo, error) {
	parameter := make(map[string]interface{})
	parameter["addresses"] = addresses
	result, err := CallAndUnmarshal("listunspent", parameter, config)
	if err != nil {
		return nil, err
	}

	var utxoInfos []base.UTXOInfo
	if err := Unmarshal(&result, &utxoInfos); err != nil {
		return nil, err
	}
	return utxoInfos, nil
}

func GetReceivedByAddress(address string, config *config.RpcConfig) (string, error) {
	balanceInterface, err := CallAndUnmarshal("getreceivedbyaddress", Param("address", address), config)
	if err != nil {
		return "", err
	}
	balance := balanceInterface.(string)
	return balance, nil
}

func GetNodeState(config *config.RpcConfig) (*servers.ServerInfo, error) {
	result, err := CallAndUnmarshal("getnodestate", nil, config)
	if err != nil {
		log.Printf("getnodestate err: %s\n", err.Error())
		return nil, err
	}

	serverInfo := &servers.ServerInfo{}
	if err := Unmarshal(&result, serverInfo); err != nil {
		return nil, err
	}

	return serverInfo, nil
}

func GetNeighbors(config *config.RpcConfig) ([]string, error) {
	result, err := CallAndUnmarshal("getneighbors", nil, config)
	if err != nil {
		log.Printf("getneighbors err: %s\n", err.Error())
		return []string{}, err
	}

	var neighbors []string
	if err := Unmarshal(&result, &neighbors); err != nil {
		return []string{}, err
	}

	return neighbors, nil
}

func GetMempool(config *config.RpcConfig) ([]string, error) {
	result, err := CallAndUnmarshal("getrawmempool", nil, config)
	if err != nil {
		log.Printf("getrawmempool err: %s\n", err.Error())
		return []string{}, err
	}

	var txHashes []string
	if err := Unmarshal(&result, &txHashes); err != nil {
		log.Printf("Unmarshal txHashes from mempool err: %s\n", err.Error())
		return []string{}, err
	}

	return txHashes, nil
}

func GetMempoolAll(config *config.RpcConfig) ([]*servers.TransactionContextInfo, error) {
	parameter := make(map[string]interface{})
	parameter["state"] = "all"

	result, err := CallAndUnmarshal("getrawmempool", parameter, config)
	if err != nil {
		log.Printf("getrawmempool all err: %s", err.Error())
		return []*servers.TransactionContextInfo{}, err
	}

	var txContextInfo []*servers.TransactionContextInfo
	if err := Unmarshal(&result, &txContextInfo); err != nil {
		log.Printf("Unmarshal TransactionContextInfo from mempool err: %s\n", err.Error())
		return []*servers.TransactionContextInfo{}, err
	}

	return txContextInfo, nil
}

func post(url string, contentType string, user string, pass string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	auth := user + ":" + pass
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Content-Type", contentType)

	client := *http.DefaultClient
	client.Timeout = time.Minute
	return client.Do(req)
}

func Call(method string, params map[string]interface{}, config *config.RpcConfig) ([]byte, error) {
	url := "http://" + config.IpAddress + ":" + strconv.Itoa(config.HttpJsonPort)
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"params": params,
	})
	if err != nil {
		return nil, err
	}

	resp, err := post(url, "application/json", config.User, config.Pass, strings.NewReader(string(data)))
	if err != nil {
		log.Printf("POST requset err:%s", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func CallAndUnmarshal(method string, params map[string]interface{}, config *config.RpcConfig) (interface{}, error) {
	body, err := Call(method, params, config)
	if err != nil {
		return nil, err
	}

	resp := Response{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return string(body), nil
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return resp.Result, nil
}

func CallAndUnmarshalResponse(method string, params map[string]interface{}, config *config.RpcConfig) (Response, error) {
	body, err := Call(method, params, config)
	if err != nil {
		return Response{}, err
	}

	resp := Response{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

func Unmarshal(result interface{}, target interface{}) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

func GetTransaction(tx string, config *config.RpcConfig) (*types.Transaction, error) {
	param := make(map[string]interface{})
	param["txid"] = tx
	resp, err := CallAndUnmarshalResponse("getrawtransaction", param,
		config)
	if err != nil {
		return nil, errors.New("[MoniterFailedDepositTransfer] Unable to call getfaileddeposittransactions rpc " + err.Error())
	}
	rawTx, ok := resp.Result.(string)
	if !ok {
		return nil, errors.New("[MoniterFailedDepositTransfer] Getrawtransaction rpc result not correct ")
	}
	buf, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, errors.New("[MoniterFailedDepositTransfer] Invalid data from GetSmallCrossTransferTxs " + err.Error())
	}
	var txn types.Transaction
	err = txn.Deserialize(bytes.NewReader(buf))
	if err != nil {
		return nil, errors.New("[MoniterFailedDepositTransfer] Decode transaction error " + err.Error())
	}
	return &txn, nil
}
