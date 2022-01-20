// Copyright (c) 2017-2022 The Elastos Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"
	"log"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
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

	if request.TransactionIdentifier == nil {
		return nil, errors.NoTransactionIdentifier
	}

	for _, tx := range txContextInfo {
		if tx.Hash == request.TransactionIdentifier.Hash {
			rstx, e := GetRosettaTransactionByTxInfo(tx.TransactionInfo, &base.MainnetDefaultStatus)
			if e != nil {
				return nil, e
			}
			return &types.MempoolTransactionResponse{
				Transaction: rstx,
				Metadata:    nil,
			}, nil
		}
	}

	return nil, errors.TransactionNotExist
}
