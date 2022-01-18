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

package client

import (
	"log"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/types"
)

func init() {
	t := new(testing.T)
	Test_network(t)
}

func Test_block(t *testing.T) {
	client := create_test_client()
	ctx := context()
	if primaryNetwork == nil || networkStatus == nil || networkOptions == nil {
		log.Fatal("primaryNetwork or networkStatus or networkOptions has not been initialized")
	}
	// Step 10: Create an asserter using the retrieved NetworkStatus and
	// NetworkOptions.
	//
	// This will be used later to assert that a fetched block is
	// valid.
	asserter, err := asserter.NewClientWithResponses(
		primaryNetwork,
		networkStatus,
		networkOptions,
		"",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Step 11: Fetch the current block
	block, rosettaErr, err := client.BlockAPI.Block(
		ctx,
		&types.BlockRequest{
			NetworkIdentifier: primaryNetwork,
			BlockIdentifier: types.ConstructPartialBlockIdentifier(
				networkStatus.CurrentBlockIdentifier,
			),
		},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Step 12: Print the block
	log.Printf("Current Block: %s\n", types.PrettyPrintStruct(block.Block))

	// Step 13: Assert the block response is valid
	//
	// It is important to note that this only ensures
	// required fields are populated and that operations
	// in the block only use types and statuses that were
	// provided in the networkStatusResponse. To run more
	// intensive validation, use the Rosetta CLI. It
	// can be found at: https://github.com/coinbase/rosetta-cli
	err = asserter.Block(block.Block)
	if err != nil {
		log.Fatalf("Assertion Error: %s\n", err.Error())
	}

	// Step 14: Print remaining transactions to fetch
	//
	// If you want the client to automatically fetch these, consider
	// using the fetcher package.
	for _, txn := range block.OtherTransactions {
		log.Printf("Other Transaction: %+v\n", txn)
	}
}
