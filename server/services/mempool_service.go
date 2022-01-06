// Copyright (c) 2017-2022 The Elastos Foundation
//
// The MIT License(MIT)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

	//txContextInfo, err := rpc.GetMempoolAll(config.Parameters.MainNodeRPC)
	//if err != nil {
	//	log.Printf("GetMempoolAll err: %s\n", err.Error())
	//	return nil, errors.GetMempoolFailed
	//}

	return &types.MempoolTransactionResponse{}, nil
}
