package services

import (
	"context"
	"log"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/errors"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/rpc"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// MempoolAPIService implements the server.MempoolAPIServicer interface.
type MempoolAPIService struct {
	network *types.NetworkIdentifier
}

// NewMempoolAPIService creates a new instance of a MempoolAPIService.
func NewMempoolAPIService(network *types.NetworkIdentifier) server.MempoolAPIServicer {
	return &MempoolAPIService{
		network: network,
	}
}

func (s *MempoolAPIService) Mempool(ctx context.Context, request *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	txHashes, err := rpc.GetMempool(config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetMempool err: %s\n", err.Error())
		return nil, errors.GetMempoolFailed
	}

	txIdentifiers := make([]*types.TransactionIdentifier, 0)
	for _, txHash := range txHashes {
		txIdentifiers = append(txIdentifiers, &types.TransactionIdentifier{
			Hash: txHash,
		})
	}

	return &types.MempoolResponse{
		TransactionIdentifiers: txIdentifiers,
	}, nil
}

func (s *MempoolAPIService) MempoolTransaction(
	ctx context.Context,
	request *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	txContextInfo, err := rpc.GetMempoolAll(config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetMempoolAll err: %s\n", err.Error())
		return nil, errors.GetMempoolFailed
	}

	for _, tx := range txContextInfo {
		if tx.Hash == request.TransactionIdentifier.Hash {
			rstx, e := GetRosettaTransactionByTxInfo(tx.TransactionInfo)
			if e != nil {
				return nil, e
			}
			return &types.MempoolTransactionResponse{
				Transaction: rstx,
				Metadata:    nil,
			}, nil
		}
	}

	return &types.MempoolTransactionResponse{}, nil
}
