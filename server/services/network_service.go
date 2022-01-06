// Copyright 2020 Coinbase, Inc.
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

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/errors"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/rpc"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type NetworkAPIService struct {
	network *types.NetworkIdentifier
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(network *types.NetworkIdentifier) server.NetworkAPIServicer {
	return &NetworkAPIService{
		network: network,
	}
}

// NetworkList implements the /network/list endpoint
func (s *NetworkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			s.network,
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *NetworkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
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

	genesisBlockInfo, err := rpc.GetBlockByHeight(0, config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetBlockByHeight err: %s\n", err.Error())
		return nil, errors.BlockNotExist
	}

	neighbors, err := rpc.GetNeighbors(config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("GetNeighbors err: %s\n", err.Error())
		return nil, errors.GetNeighborsFailed
	}

	peers := make([]*types.Peer, 0, len(neighbors))
	for _, n := range neighbors {
		peer := types.Peer{
			PeerID: n,
		}
		peers = append(peers, &peer)
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: int64(currentHeight),
			Hash:  blockInfo.Hash,
		},
		CurrentBlockTimestamp: int64(blockInfo.Time),
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: int64(genesisBlockInfo.Height),
			Hash:  genesisBlockInfo.Hash,
		},
		Peers: peers,
	}, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *NetworkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	nodeState, err := rpc.GetNodeState(config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("get node state err: %s\n", err.Error())
		return nil, errors.GetNodeStateFailed
	}

	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: "1.4.10",
			NodeVersion:    nodeState.Compile,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     "Success",
					Successful: true,
				},
				{
					Status:     "Reverted",
					Successful: false,
				},
			},
			OperationTypes: []string{
				"Transfer",
			},
			Errors:                  errors.APIErrorMap,
			HistoricalBalanceLookup: false,
			CallMethods:             []string{},
			BalanceExemptions:       []*types.BalanceExemption{},
			MempoolCoins:            true,
		},
	}, nil
}
