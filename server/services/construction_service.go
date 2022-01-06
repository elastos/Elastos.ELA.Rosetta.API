package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionAPIServicer implements the server.ConstructionAPIServicer interface.
type ConstructionAPIServicer struct {
	network *types.NetworkIdentifier
}

// NewConstructionAPIServicer creates a new instance of a ConstructionAPIServicer.
func NewConstructionAPIServicer(network *types.NetworkIdentifier) server.ConstructionAPIServicer {
	return &ConstructionAPIServicer{
		network: network,
	}
}

func (s *ConstructionAPIServicer) ConstructionCombine(
	context.Context,
	*types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	return nil, nil
}

func (s *ConstructionAPIServicer) ConstructionDerive(
	context.Context,
	*types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	return nil, nil
}

func (s *ConstructionAPIServicer) ConstructionHash(
	context.Context,
	*types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	return nil, nil
}

func (s *ConstructionAPIServicer) ConstructionMetadata(
	context.Context,
	*types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	return nil, nil
}

func (s *ConstructionAPIServicer) ConstructionParse(
	context.Context,
	*types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	return nil, nil
}

func (s *ConstructionAPIServicer) ConstructionPayloads(
	context.Context,
	*types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	return nil, nil
}

func (s *ConstructionAPIServicer) ConstructionPreprocess(
	context.Context,
	*types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	return nil, nil
}

func (s *ConstructionAPIServicer) ConstructionSubmit(
	context.Context,
	*types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	return nil, nil
}
