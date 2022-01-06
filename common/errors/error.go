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

package errors

import "github.com/coinbase/rosetta-sdk-go/types"

const (
	DefaultError int32 = iota
	TransactionNotExistError
	ReferTransactionNotExistError
	GetCurrentBlockError
	BlockNotExistError
	GetNeighborsError
	GetNodeStateError
	GetMempoolError
	DecodeAddressError
	EncodeToAddressError
	UnsupportNetworkError
)

var (
	TransactionNotExist = &types.Error{
		Code:      TransactionNotExistError,
		Message:   "failed to get transaction by rpc",
		Retriable: false,
	}

	ReferTransactionNotExist = &types.Error{
		Code:      ReferTransactionNotExistError,
		Message:   "failed to get transaction reference by rpc",
		Retriable: false,
	}

	GetCurrentBlockFailed = &types.Error{
		Code:      GetCurrentBlockError,
		Message:   "failed to get current block height",
		Retriable: false,
	}

	BlockNotExist = &types.Error{
		Code:      BlockNotExistError,
		Message:   "failed to get block by rpc",
		Retriable: false,
	}

	GetNeighborsFailed = &types.Error{
		Code:      GetNeighborsError,
		Message:   "failed to get neighbors",
		Retriable: false,
	}

	GetNodeStateFailed = &types.Error{
		Code:      GetNodeStateError,
		Message:   "failed to get node state",
		Retriable: false,
	}

	GetMempoolFailed = &types.Error{
		Code:      GetMempoolError,
		Message:   "failed to get mempool",
		Retriable: false,
	}

	DecodeAddress = &types.Error{
		Code:      DecodeAddressError,
		Message:   "failed to decode address",
		Retriable: false,
	}

	EncodeToAddress = &types.Error{
		Code:      EncodeToAddressError,
		Message:   "failed to encode to address",
		Retriable: false,
	}

	UnsupportNetwork = &types.Error{
		Code:      UnsupportNetworkError,
		Message:   "unsupport network",
		Retriable: false,
	}
)

var APIErrorMap = []*types.Error{
	TransactionNotExist,
	ReferTransactionNotExist,
	GetCurrentBlockFailed,
	BlockNotExist,
	GetNeighborsFailed,
	GetNodeStateFailed,
	GetMempoolFailed,
	DecodeAddress,
	EncodeToAddress,
	UnsupportNetwork,
}
