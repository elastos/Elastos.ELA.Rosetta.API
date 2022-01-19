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

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
}
func Test_account(t *testing.T) {
	skipShort(t)
	client := create_test_client()
	ctx := context()
	if primaryNetwork == nil || networkStatus == nil || networkOptions == nil {
		log.Fatal("primaryNetwork or networkStatus or networkOptions has not been initialized")
	}

	// Step 1: Fetch the account coins
	result, rosettaErr, err := client.AccountAPI.AccountCoins(
		ctx,
		&types.AccountCoinsRequest{
			NetworkIdentifier: primaryNetwork,
			AccountIdentifier: &types.AccountIdentifier{
				Address:    "EWEzHxxHewc5w2nRMsn9Q7PFCs2V9cuaZ9",
				SubAccount: nil,
				Metadata:   nil,
			},
			IncludeMempool: false,
			Currencies: []*types.Currency{
				&types.Currency{
					Symbol:   "ELA",
					Decimals: 8,
					Metadata: nil,
				},
			},
		},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Step 2: print the account coins result
	log.Printf("Account result: %s\n", types.PrettyPrintStruct(result.Coins))

	assert.Equal(t, "50000000000", result.Coins[0].Amount.Value)
}
