package services

import (
	"context"
	"log"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/errors"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/rpc"
	"github.com/elastos/Elastos.ELA/common"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// AccountAPIService implements the server.AccountAPIServicer interface.
type AccountAPIService struct {
	network *types.NetworkIdentifier
}

// NewAccounAPIService creates a new instance of a AccountAPIService.
func NewAccounAPIService(network *types.NetworkIdentifier) server.AccountAPIServicer {
	return &AccountAPIService{
		network: network,
	}
}

func (s *AccountAPIService) AccountBalance(
	ctx context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	if request.AccountIdentifier == nil {
		return nil, errors.InvalidAccountIdentifier
	}

	balance, err := rpc.GetReceivedByAddress(request.AccountIdentifier.Address, config.Parameters.MainNodeRPC)
	if err != nil {
		errStr := err.Error()
		log.Printf("err: %s\n", errStr)
		return nil, errors.GetAddressBalanceFailed
	}
	b, err := common.StringToFixed64(balance)
	if err != nil {
		return nil, errors.InvalidAmount
	}
	log.Printf("Address %s balance: %v\n", request.AccountIdentifier.Address, balance)

	var amountSlice []*types.Amount

	amount := &types.Amount{
		Value: GetSelaString(*b),
		Currency: &types.Currency{
			Symbol:   base.MainnetCurrencySymbol,
			Decimals: 8,
			Metadata: nil,
		},
	}
	amountSlice = append(amountSlice, amount)

	currentHeight, err := rpc.GetCurrentHeight(config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetCurrentHeight err: %s\n", err.Error())
		return nil, errors.GetCurrentBlockFailed
	}

	blockInfo, err := rpc.GetBlockByHeight(currentHeight, config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetBlockByHeight err: %s\n", err.Error())
		return nil, errors.BlockNotExist
	}

	if request.BlockIdentifier == nil {
		return &types.AccountBalanceResponse{
			BlockIdentifier: &types.BlockIdentifier{
				// This is also known as the block height.
				Index: int64(currentHeight),
				Hash:  blockInfo.Hash,
			},
			Balances: amountSlice,
			Metadata: nil,
		}, nil
	}
	// todo support blockIdentifier
	return &types.AccountBalanceResponse{
		BlockIdentifier: &types.BlockIdentifier{
			// This is also known as the block height.
			Index: int64(currentHeight),
			Hash:  blockInfo.Hash,
		},
		Balances: amountSlice,
		Metadata: nil,
	}, nil
}

// BlockTransaction implements the /block/transaction endpoint.
func (s *AccountAPIService) AccountCoins(
	ctx context.Context,
	request *types.AccountCoinsRequest,
) (*types.AccountCoinsResponse, *types.Error) {

	if request.IncludeMempool != false {
		return nil, errors.InvalidIncludeMempool
	}

	if request.NetworkIdentifier == nil {
		return nil, errors.NoNetworkIdentifier
	}

	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}
	if request.AccountIdentifier == nil {
		return nil, errors.InvalidAccountIdentifier
	}

	currentHeight, err := rpc.GetCurrentHeight(config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetCurrentHeight err: %s\n", err.Error())
		return nil, errors.GetCurrentBlockFailed
	}

	blockInfo, err := rpc.GetBlockByHeight(currentHeight, config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetBlockByHeight err: %s\n", err.Error())
		return nil, errors.BlockNotExist
	}

	var addresses []string
	addresses = append(addresses, request.AccountIdentifier.Address)
	utxoInfoSlice, err := rpc.GetUnspentUtxo(addresses, config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetUnspentUtxo err: %s\n", err.Error())
		return nil, errors.GetUnspentUtxoFailed
	}

	coinsSlice := make([]*types.Coin, 0, len(utxoInfoSlice))
	for _, utxoInfo := range utxoInfoSlice {
		value, err := common.StringToFixed64(utxoInfo.Amount)
		if err != nil {
			return nil, errors.InvalidAmount
		}
		coin := &types.Coin{
			CoinIdentifier: &types.CoinIdentifier{
				Identifier: GetCoinIdentifierByHashStr(utxoInfo.Txid, uint16(utxoInfo.VOut)),
			},
			Amount: &types.Amount{
				Value: GetSelaString(*value),
				Currency: &types.Currency{
					Symbol:   base.MainnetCurrencySymbol,
					Decimals: 8,
					Metadata: nil,
				},
			},
		}
		coinsSlice = append(coinsSlice, coin)
	}
	return &types.AccountCoinsResponse{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(blockInfo.Height),
			Hash:  blockInfo.Hash,
		},
		Coins:    coinsSlice,
		Metadata: nil,
	}, nil
}
