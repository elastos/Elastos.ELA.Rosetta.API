package base

import (
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/core/types/payload"
)

var MainnetNetworkIndex int64 = 0
var MainnetStatus = "Success"

const BlockChainName = "Elastos"
const MainnetNextworkType = "Transfer"
const MainnetCurrencySymbol = "ELA"
const MainnetCurrencyDecimal int32 = 8
const MainnetCurveType = "secp256r1"

type ProgramInfo struct {
	Code      string
	Parameter string
}

type TxoutputMap struct {
	Key   common.Uint256
	Txout []OutputInfo
}

type AmountMap struct {
	Key   common.Uint256
	Value common.Fixed64
}

type AttributeInfo struct {
	Usage types.AttributeUsage `json:"usage"`
	Data  string               `json:"data"`
}

type InputInfo struct {
	TxID     string `json:"txid"`
	VOut     uint16 `json:"vout"`
	Sequence uint32 `json:"sequence"`
}

type OutputInfo struct {
	Value      string `json:"value"`
	Index      uint32 `json:"n"`
	Address    string `json:"address"`
	OutputLock uint32 `json:"outputlock"`
}

type BlockInfo struct {
	Hash              string        `json:"hash"`
	Confirmations     uint32        `json:"confirmations"`
	StrippedSize      uint32        `json:"strippedsize"`
	Size              uint32        `json:"size"`
	Weight            uint32        `json:"weight"`
	Height            uint32        `json:"height"`
	Version           uint32        `json:"version"`
	VersionHex        string        `json:"versionhex"`
	MerkleRoot        string        `json:"merkleroot"`
	Tx                []interface{} `json:"tx"`
	Time              uint32        `json:"time"`
	MedianTime        uint32        `json:"mediantime"`
	Nonce             uint32        `json:"nonce"`
	Bits              uint32        `json:"bits"`
	Difficulty        string        `json:"difficulty"`
	ChainWork         string        `json:"chainwork"`
	PreviousBlockHash string        `json:"previousblockhash"`
	NextBlockHash     string        `json:"nextblockhash"`
	AuxPow            string        `json:"auxpow"`
}

type PayloadInfo interface {
}

type RegisterAssetInfo struct {
	Asset      *payload.Asset
	Amount     string
	Controller string
}

type CoinbaseInfo struct {
	CoinbaseData string
}

type TransferAssetInfo struct {
}

type UTXOInfo struct {
	AssetId       string `json:"assetid"`
	Txid          string `json:"txid"`
	VOut          uint32 `json:"vout"`
	Address       string `json:"address"`
	Amount        string `json:"amount"`
	Confirmations uint32 `json:"confirmations"`
	OutputLock    uint32 `json:"OutputLock"`
}
