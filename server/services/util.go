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
)

func GetSelaString(value common.Fixed64) string {
	return strconv.Itoa(int(value))
}

func GetCoinIdentifier(hash common.Uint256, index uint16) string {
	return hash.String() + ":" + strconv.Itoa(int(index))
}

func GetOperations(tx *types2.Transaction) ([]*types.Operation, *types.Error) {
	operations := make([]*types.Operation, 0)
	for i, input := range tx.Inputs {
		referTransactionHash := input.Previous.TxID
		referTransaction, err := rpc.GetTransactionByHash(input.Previous.TxID.String(), config.Parameters.MainNodeRPC)
		if err != nil {
			return nil, errors.TransactionNotExist
		}
		addr, err := referTransaction.Outputs[input.Previous.Index].ProgramHash.ToAddress()
		if err != nil {
			return nil, errors.EncodeToAddress
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
			return nil, errors.EncodeToAddress
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
