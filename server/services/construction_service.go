package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/errors"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/rpc"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/elastos/Elastos.ELA/common"
	contract2 "github.com/elastos/Elastos.ELA/core/contract"
	pg "github.com/elastos/Elastos.ELA/core/contract/program"
	elatypes "github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/crypto"
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
	ctx context.Context,
	request *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	txUnsignedBytes, err := hex.DecodeString(request.UnsignedTransaction)
	if err != nil {
		log.Printf("decode tx from hexstring err: %s\n", err.Error())
		return nil, errors.DecodeTransactionFailed
	}

	var txn elatypes.Transaction
	err = txn.Deserialize(bytes.NewReader(txUnsignedBytes))
	if err != nil {
		log.Printf("deserialize tx err: %s\n", err.Error())
		return nil, errors.DeserializeTransactionFailed
	}

	for _, sign := range request.Signatures {
		if err := checkCurveType(sign.PublicKey.CurveType); err != nil {
			log.Printf("invalid curve type")
			return nil, err
		}

		if sign.SignatureType != base.MainnetSignatureType {
			log.Printf("invalid signature type")
			return nil, errors.InvalidSignatureType
		}

		pubkey, err := crypto.DecodePoint(sign.PublicKey.Bytes)
		if err != nil {
			log.Printf("decode pubkey err: %s\n", err.Error())
			return nil, errors.DecodePublicKeyFailed
		}

		redeemScript, err := contract2.CreateStandardRedeemScript(pubkey)
		if err != nil {
			log.Printf("pubkey to code err: %s\n", err.Error())
			return nil, errors.InvalidPublicKey
		}

		var txProgram = &pg.Program{
			Code:      redeemScript,
			Parameter: sign.Bytes,
		}
		txn.Programs = append(txn.Programs, txProgram)
	}

	buf := new(bytes.Buffer)
	if err := txn.Serialize(buf); err != nil {
		log.Printf("tx serialize err: %s\n", err.Error())
		return nil, errors.InvalidTransaction
	}

	return &types.ConstructionCombineResponse{
		SignedTransaction: common.BytesToHexString(buf.Bytes()),
	}, nil
}

func (s *ConstructionAPIServicer) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	if err := checkCurveType(request.PublicKey.CurveType); err != nil {
		return nil, err
	}

	pkBytes := request.PublicKey.Bytes
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

	return &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address:    addr,
			SubAccount: nil,
			Metadata:   nil,
		},
		Metadata: nil,
	}, nil
}

func (s *ConstructionAPIServicer) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	txBytes, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		log.Printf("decode tx from hexstring err: %s\n", err.Error())
		return nil, errors.DecodeTransactionFailed
	}

	var txn elatypes.Transaction
	err = txn.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		log.Printf("deserialize tx err: %s\n", err.Error())
		return nil, errors.DecodeTransactionFailed
	}

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txn.Hash().String(),
		},
		Metadata: nil,
	}, nil
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
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	return &types.ConstructionPreprocessResponse{}, nil
}

func (s *ConstructionAPIServicer) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	if !CheckNetwork(request.NetworkIdentifier) {
		log.Printf("unsupport network")
		return nil, errors.UnsupportNetwork
	}

	buf, err := hex.DecodeString(request.SignedTransaction)
	if err != nil {
		log.Printf("decode tx from hexstring err: %s\n", err.Error())
		return nil, errors.DecodeTransactionFailed
	}

	var txn elatypes.Transaction
	err = txn.Deserialize(bytes.NewReader(buf))
	if err != nil {
		log.Printf("tx deserialize err: %s\n", err.Error())
		return nil, errors.DeserializeTransactionFailed
	}

	txHash, err := rpc.PublishTransaction(request.SignedTransaction, config.Parameters.MainNodeRPC)
	if err != nil {
		log.Printf("publishtransaction err: %s\n", err.Error())
		return nil, errors.PublishTransactionFailed
	}

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txHash,
		},
	}, nil
}
