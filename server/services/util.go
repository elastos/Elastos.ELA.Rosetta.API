package services

import (
	"strconv"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/errors"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/rpc"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/elastos/Elastos.ELA/common"
	types2 "github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/servers"
)

func GetRosettaTimestamp(time uint32) int64 {
	return int64(time * 1000)
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

func GetOperations(tx *types2.Transaction) ([]*types.Operation, *types.Error) {
	operations := make([]*types.Operation, 0)
	for i, input := range tx.Inputs {
		referTransactionHash := input.Previous.TxID
		referTransaction, err := rpc.GetTransaction(input.Previous.TxID.String(), config.Parameters.MainNodeRPC)
		if err != nil {
			return nil, errors.TransactionNotExist
		}
		addr, err := referTransaction.Outputs[input.Previous.Index].ProgramHash.ToAddress()
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
			Status:            &base.MainnetStatus,
			Account: &types.AccountIdentifier{
				Address:    addr,
				SubAccount: nil,
				Metadata:   nil,
			},
			Amount: &types.Amount{
				Value: GetSelaString(referTransaction.Outputs[input.Previous.Index].Value),
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
				CoinAction: "coin_spent",
			},
			Metadata: nil,
		})
	}

	for i, output := range tx.Outputs {
		addr, err := output.ProgramHash.ToAddress()
		if err != nil {
			return nil, errors.EncodeToAddressFailed
		}

		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        int64(len(tx.Inputs) + i),
				NetworkIndex: &base.MainnetNetworkIndex,
			},
			RelatedOperations: nil,
			Type:              base.MainnetNextworkType,
			Status:            &base.MainnetStatus,
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
					Identifier: GetCoinIdentifier(tx.Hash(), uint16(i)),
				},
				CoinAction: "coin_created",
			},
			Metadata: nil,
		})
	}

	return operations, nil
}

func GetOperationsByTxInfo(tx *servers.TransactionInfo) ([]*types.Operation, *types.Error) {
	operations := make([]*types.Operation, 0)
	for i, input := range tx.Inputs {
		referTransactionHash := input.TxID
		referTransaction, err := rpc.GetTransaction(input.TxID, config.Parameters.MainNodeRPC)
		if err != nil {
			return nil, errors.TransactionNotExist
		}
		addr, err := referTransaction.Outputs[input.VOut].ProgramHash.ToAddress()
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
			Status:            &base.MainnetStatus,
			Account: &types.AccountIdentifier{
				Address:    addr,
				SubAccount: nil,
				Metadata:   nil,
			},
			Amount: &types.Amount{
				Value: GetSelaString(referTransaction.Outputs[input.VOut].Value),
				Currency: &types.Currency{
					Symbol:   base.MainnetCurrencySymbol,
					Decimals: base.MainnetCurrencyDecimal,
					Metadata: nil,
				},
				Metadata: nil,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: GetCoinIdentifierByHashStr(referTransactionHash, input.VOut),
				},
				CoinAction: "coin_spent",
			},
			Metadata: nil,
		})
	}

	for i, output := range tx.Outputs {
		addr := output.Address
		amount, err := common.StringToFixed64(output.Value)
		if err != nil {
			return nil, errors.InvalidTransaction
		}

		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        int64(len(tx.Inputs) + i),
				NetworkIndex: &base.MainnetNetworkIndex,
			},
			RelatedOperations: nil,
			Type:              base.MainnetNextworkType,
			Status:            &base.MainnetStatus,
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
				CoinAction: "coin_created",
			},
			Metadata: nil,
		})
	}

	return operations, nil
}

func GetRosettaTransaction(tx *types2.Transaction) (*types.Transaction, *types.Error) {
	operations, e := GetOperations(tx)
	if e != nil {
		return nil, e
	}
	return &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: tx.Hash().String(),
		},
		Operations:          operations,
		RelatedTransactions: nil,
		Metadata:            nil,
	}, nil
}

func GetRosettaTransactionByTxInfo(tx *servers.TransactionInfo) (*types.Transaction, *types.Error) {
	operations, e := GetOperationsByTxInfo(tx)
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

func GetRosettaBlock(block *base.BlockInfo) (*types.Block, *types.Error) {
	var txs []*types.Transaction
	for _, t := range block.Tx {
		tx, ok := t.(*types2.Transaction)
		if !ok {
			return nil, errors.InvalidTransaction
		}
		rstx, e := GetRosettaTransaction(tx)
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
