package services

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/errors"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/rpc"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/elastos/Elastos.ELA/common"
	contract2 "github.com/elastos/Elastos.ELA/core/contract"
	types2 "github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/crypto"
	"github.com/elastos/Elastos.ELA/servers"
)

func GetRosettaTimestamp(time uint32) int64 {
	return int64(time) * 1000
}

func GetSelaString(value common.Fixed64) string {
	return strconv.Itoa(int(value))
}

func GetCoinIdentifier(hash common.Uint256, index uint16) string {
	return hash.String() + ":" + strconv.Itoa(int(index))
}

func GetCoinIdentifierByHashStr(hash string, index uint16) string {
	return hash + ":" + strconv.Itoa(int(index))
}

func CheckNetwork(network *types.NetworkIdentifier) bool {
	if network.Blockchain == base.BlockChainName && network.Network == config.Parameters.ActiveNet {
		return true
	}

	return false
}

func GetOperations(tx *types2.Transaction, status *string) ([]*types.Operation, *types.Error) {
	operations := make([]*types.Operation, 0)

	for i, output := range tx.Outputs {
		addr, err := output.ProgramHash.ToAddress()
		if err != nil {
			return nil, errors.EncodeToAddressFailed
		}

		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        int64(i),
				NetworkIndex: &base.MainnetNetworkIndex,
			},
			RelatedOperations: nil,
			Type:              base.MainnetNextworkType,
			Status:            status,
			Account: &types.AccountIdentifier{
				Address:    addr,
				SubAccount: nil,
				Metadata:   nil,
			},
			Amount: &types.Amount{
				Value: GetSelaString(output.Value),
				Currency: &types.Currency{
					Symbol:   base.MainnetCurrencySymbol,
					Decimals: base.MainnetCurrencyDecimal,
					Metadata: nil,
				},
				Metadata: nil,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: GetCoinIdentifierByHashStr(common.ToReversedString(tx.Hash()), uint16(i)),
				},
				CoinAction: types.CoinCreated,
			},
			Metadata: nil,
		})
	}

	if tx.TxType != types2.CoinBase {
		outpusCount := len(tx.Outputs)
		for i, input := range tx.Inputs {
			referTransactionHash := input.Previous.TxID
			referTransaction, err := rpc.GetTransaction(common.ToReversedString(input.Previous.TxID), config.Parameters.MainNodeRPC)
			if err != nil {
				log.Println("get transaction error:", err, "hash:", common.ToReversedString(input.Previous.TxID))
				return nil, errors.TransactionNotExist
			}
			addr, err := referTransaction.Outputs[input.Previous.Index].ProgramHash.ToAddress()
			if err != nil {
				return nil, errors.EncodeToAddressFailed
			}

			operations = append(operations, &types.Operation{
				OperationIdentifier: &types.OperationIdentifier{
					Index:        int64(outpusCount + i),
					NetworkIndex: &base.MainnetNetworkIndex,
				},
				RelatedOperations: nil,
				Type:              base.MainnetNextworkType,
				Status:            status,
				Account: &types.AccountIdentifier{
					Address:    addr,
					SubAccount: nil,
					Metadata:   nil,
				},
				Amount: &types.Amount{
					Value: GetSelaString(-referTransaction.Outputs[input.Previous.Index].Value),
					Currency: &types.Currency{
						Symbol:   base.MainnetCurrencySymbol,
						Decimals: base.MainnetCurrencyDecimal,
						Metadata: nil,
					},
					Metadata: nil,
				},
				CoinChange: &types.CoinChange{
					CoinIdentifier: &types.CoinIdentifier{
						Identifier: GetCoinIdentifier(referTransactionHash, input.Previous.Index),
					},
					CoinAction: types.CoinSpent,
				},
				Metadata: nil,
			})
		}
	}

	return operations, nil
}

func GetReversedString(txid string) (*string, error) {
	hex, err := common.FromReversedString(txid)
	if err != nil {
		return nil, err
	}
	var hash common.Uint256
	err = hash.Deserialize(bytes.NewReader(hex))
	if err != nil {
		return nil, err
	}
	reversedTxid := hash.String()
	return &reversedTxid, nil
}

func GetOperationsByTxInfo(tx *servers.TransactionInfo, status *string) ([]*types.Operation, *types.Error) {
	operations := make([]*types.Operation, 0)

	// record output index first, then record input index
	for i, output := range tx.Outputs {
		addr := output.Address
		amount, err := common.StringToFixed64(output.Value)
		if err != nil {
			return nil, errors.InvalidTransaction
		}

		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        int64(i),
				NetworkIndex: &base.MainnetNetworkIndex,
			},
			RelatedOperations: nil,
			Type:              base.MainnetNextworkType,
			Status:            status,
			Account: &types.AccountIdentifier{
				Address:    addr,
				SubAccount: nil,
				Metadata:   nil,
			},
			Amount: &types.Amount{
				Value: GetSelaString(*amount),
				Currency: &types.Currency{
					Symbol:   base.MainnetCurrencySymbol,
					Decimals: base.MainnetCurrencyDecimal,
					Metadata: nil,
				},
				Metadata: nil,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: GetCoinIdentifierByHashStr(tx.Hash, uint16(i)),
				},
				CoinAction: types.CoinCreated,
			},
			Metadata: nil,
		})
	}

	if tx.TxType != types2.CoinBase {
		outpusCount := len(tx.Outputs)
		for i, input := range tx.Inputs {
			referTransaction, err := rpc.GetTransaction(input.TxID, config.Parameters.MainNodeRPC)
			if err != nil {
				log.Println("get transaction error:", err, "hash:", input.TxID)
				return nil, errors.TransactionNotExist
			}
			addr, err := referTransaction.Outputs[input.VOut].ProgramHash.ToAddress()
			if err != nil {
				return nil, errors.EncodeToAddressFailed
			}

			operations = append(operations, &types.Operation{
				OperationIdentifier: &types.OperationIdentifier{
					Index:        int64(outpusCount + i),
					NetworkIndex: &base.MainnetNetworkIndex,
				},
				RelatedOperations: nil,
				Type:              base.MainnetNextworkType,
				Status:            status,
				Account: &types.AccountIdentifier{
					Address:    addr,
					SubAccount: nil,
					Metadata:   nil,
				},
				Amount: &types.Amount{
					Value: GetSelaString(-referTransaction.Outputs[input.VOut].Value),
					Currency: &types.Currency{
						Symbol:   base.MainnetCurrencySymbol,
						Decimals: base.MainnetCurrencyDecimal,
						Metadata: nil,
					},
					Metadata: nil,
				},
				CoinChange: &types.CoinChange{
					CoinIdentifier: &types.CoinIdentifier{
						Identifier: GetCoinIdentifierByHashStr(input.TxID, input.VOut),
					},
					CoinAction: types.CoinSpent,
				},
				Metadata: nil,
			})
		}
	}

	return operations, nil
}

func GetRosettaTransaction(tx *types2.Transaction, status *string) (*types.Transaction, *types.Error) {
	operations, e := GetOperations(tx, status)
	if e != nil {
		return nil, e
	}

	return &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: common.ToReversedString(tx.Hash()),
		},
		Operations:          operations,
		RelatedTransactions: nil,
		Metadata:            nil,
	}, nil
}

func GetRosettaTransactionByTxInfo(tx *servers.TransactionInfo, status *string) (*types.Transaction, *types.Error) {
	operations, e := GetOperationsByTxInfo(tx, status)
	if e != nil {
		return nil, e
	}
	return &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: tx.Hash,
		},
		Operations:          operations,
		RelatedTransactions: nil,
		Metadata:            nil,
	}, nil
}

func GetRosettaBlock(block *types2.Block, status *string) (*types.Block, *types.Error) {
	var txs []*types.Transaction
	for _, t := range block.Transactions {
		rstx, e := GetRosettaTransaction(t, status)
		if e != nil {
			return nil, e
		}
		txs = append(txs, rstx)
	}

	var previousBlockIndex int64
	if block.Height > 1 {
		previousBlockIndex = int64(block.Height - 1)
	}
	return &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(block.Height),
			Hash:  common.ToReversedString(block.Hash()),
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: previousBlockIndex,
			Hash:  common.ToReversedString(block.Previous),
		},
		Timestamp:    GetRosettaTimestamp(block.Timestamp),
		Transactions: txs,
		Metadata:     nil,
	}, nil
}

func GetRosettaBlockByBlockInfo(block *base.BlockInfo, status *string) (*types.Block, *types.Error) {
	txs := make([]*types.Transaction, 0)
	for _, t := range block.Tx {
		bytes, _ := json.Marshal(t)
		var txInfo servers.TransactionContextInfo
		err := json.Unmarshal(bytes, &txInfo)
		if err != nil {
			log.Printf("invalid transaction context %v", t)
			return nil, errors.InvalidTransaction
		}
		rstx, e := GetRosettaTransactionByTxInfo(txInfo.TransactionInfo, status)
		if e != nil {
			return nil, e
		}
		txs = append(txs, rstx)
	}

	var previousBlockIndex int64
	if block.Height > 1 {
		previousBlockIndex = int64(block.Height - 1)
	}
	return &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(block.Height),
			Hash:  block.Hash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: previousBlockIndex,
			Hash:  block.PreviousBlockHash,
		},
		Timestamp:    GetRosettaTimestamp(block.Time),
		Transactions: txs,
		Metadata:     nil,
	}, nil
}

func checkCurveType(curveType types.CurveType) *types.Error {
	if curveType != base.MainnetCurveType {
		return errors.InvalidCurveType
	}
	return nil
}

func publicKeyToAddress(pkBytes []byte) (*string, *types.Error) {
	pk, err := crypto.DecodePoint(pkBytes)
	if err != nil {
		return nil, errors.InvalidCurveType
	}
	contract, err := contract2.CreateStandardContract(pk)
	if err != nil {
		return nil, errors.InvalidPublicKey
	}

	addr, err := contract.ToProgramHash().ToAddress()
	if err != nil {
		return nil, errors.InvalidPublicKey
	}

	return &addr, nil
}

func getPositiveAmountFromString(value string) (common.Fixed64, *types.Error) {
	amount, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.InvalidOperationAmount
	}
	var positiveAmount common.Fixed64
	if amount < 0 {
		positiveAmount = common.Fixed64(-amount)
	} else {
		positiveAmount = common.Fixed64(amount)
	}

	return positiveAmount, nil
}

func getParameterBySignature(signature []byte) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(len(signature)))
	buf.Write(signature)
	return buf.Bytes()
}
