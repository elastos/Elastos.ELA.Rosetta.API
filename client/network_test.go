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

var (
	primaryNetwork *types.NetworkIdentifier
	networkStatus  *types.NetworkStatusResponse
	networkOptions *types.NetworkOptionsResponse
)

func Test_network(t *testing.T) {
	ctx := context()
	client := create_test_client()
	var rosettaErr *types.Error
	var err error
	// Step 2: Get all available networks
	var networkList *types.NetworkListResponse
	networkList, rosettaErr, err = client.NetworkAPI.NetworkList(
		ctx,
		&types.MetadataRequest{},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := asserter.NetworkListResponse(networkList); err != nil {
		log.Fatal(err)
	}

	if len(networkList.NetworkIdentifiers) == 0 {
		log.Fatal("no available networks")
	}

	primaryNetwork = networkList.NetworkIdentifiers[0]

	// Step 3: Print the primary network
	log.Printf("Primary Network: %s\n", types.PrettyPrintStruct(primaryNetwork))

	// Step 4: Fetch the network status
	networkStatus, rosettaErr, err = client.NetworkAPI.NetworkStatus(
		ctx,
		&types.NetworkRequest{
			NetworkIdentifier: primaryNetwork,
		},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Step 5: Print the response
	log.Printf("Network Status: %s\n", types.PrettyPrintStruct(networkStatus))

	// Step 6: Assert the response is valid
	err = asserter.NetworkStatusResponse(networkStatus)
	if err != nil {
		log.Fatalf("Assertion Error: %s\n", err.Error())
	}

	// Step 7: Fetch the network options
	networkOptions, rosettaErr, err = client.NetworkAPI.NetworkOptions(
		ctx,
		&types.NetworkRequest{
			NetworkIdentifier: primaryNetwork,
		},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Step 8: Print the response
	log.Printf("Network Options: %s\n", types.PrettyPrintStruct(networkOptions))

	// Step 9: Assert the response is valid
	err = asserter.NetworkOptionsResponse(networkOptions)
	if err != nil {
		log.Fatalf("Assertion Error: %s\n", err.Error())
	}
}
