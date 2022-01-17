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

// BlockAPIService implements the server.BlockAPIServicer interface.
type BlockAPIService struct {
	network *types.NetworkIdentifier
}

// NewBlockAPIService creates a new instance of a BlockAPIService.
func NewBlockAPIService(network *types.NetworkIdentifier) server.BlockAPIServicer {
	return &BlockAPIService{
		network: network,
	}
}

// Block implements the /block endpoint.
func (s *BlockAPIService) Block(
	ctx context.Context,
	request *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	var block *base.BlockInfo
	if request.BlockIdentifier.Index != nil {
		var err error
		block, err = rpc.GetBlockByHeight(uint32(*request.BlockIdentifier.Index), config.Parameters.MainNodeRPC)
		if err != nil {
			return nil, errors.BlockNotExist
		}

		if request.BlockIdentifier.Hash != nil && block.Hash != *request.BlockIdentifier.Hash {
			return nil, errors.BlockNotExist
		}
	} else if request.BlockIdentifier.Hash != nil {
		var err error
		block, err = rpc.GetBlockByHash(*request.BlockIdentifier.Hash, config.Parameters.MainNodeRPC)
		if err != nil {
			return nil, errors.BlockNotExist
		}
	} else {
		height, err := rpc.GetCurrentHeight(config.Parameters.MainNodeRPC)
		if err != nil {
			return nil, errors.GetCurrentHeightFailed
		}
		block, err = rpc.GetBlockByHeight(height, config.Parameters.MainNodeRPC)
		if err != nil {
			return nil, errors.BlockNotExist
		}
	}

	rsBlock, rsError := GetRosettaBlockByBlockInfo(block, &base.MainnetStatus)
	if rsError != nil {
		return nil, rsError
	}

	return &types.BlockResponse{Block: rsBlock}, nil
}

// BlockTransaction implements the /block/transaction endpoint.
func (s *BlockAPIService) BlockTransaction(
	ctx context.Context,
	request *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	tx, err := rpc.GetTransaction(request.TransactionIdentifier.Hash, config.Parameters.MainNodeRPC)
	if err != nil {
		return nil, errors.TransactionNotExist
	}

	rstx, e := GetRosettaTransaction(tx, &base.MainnetStatus)
	if e != nil {
		return nil, e
	}

	return &types.BlockTransactionResponse{Transaction: rstx}, nil
}
