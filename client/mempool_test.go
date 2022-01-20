// Copyright (c) 2017-2022 The Elastos Foundation:
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
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
)

func init() {
	t := new(testing.T)
	Test_network(t)
}

//func Test_RunTestMempoolInloop(t *testing.T) {
//	for {
//		Test_Mempool(t)
//		time.Sleep(30 * time.Second)
//	}
//}

func Test_Mempool(t *testing.T) {
	client := create_test_client()
	ctx := context()
	if primaryNetwork == nil || networkStatus == nil || networkOptions == nil {
		log.Fatal("primaryNetwork or networkStatus or networkOptions has not been initialized")
	}
	request := &types.NetworkRequest{
		NetworkIdentifier: primaryNetwork,
		Metadata:          nil,
	}
	mempool, rosettaErr, err := client.MempoolAPI.Mempool(ctx, request)

	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Current Mempool transaction hashes : %s\n", types.PrettyPrintStruct(mempool))
	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserter, err := asserter.NewServer(
		[]string{base.MainnetNextworkType},
		false,
		[]*types.NetworkIdentifier{primaryNetwork},
		nil,
		false,
		"",
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, txIdent := range mempool.TransactionIdentifiers {
		memTxReq := &types.MempoolTransactionRequest{
			NetworkIdentifier:     primaryNetwork,
			TransactionIdentifier: txIdent,
		}
		err := asserter.MempoolTransactionRequest(memTxReq)
		if err != nil {
			log.Fatalf("Assertion Error: %s\n", err.Error())
		}
		txResp, rosettaErr, err := client.MempoolAPI.MempoolTransaction(ctx, memTxReq)
		if rosettaErr != nil {
			log.Printf("Rosetta Error: %+v\n", rosettaErr)
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Mempool transaction : %s\n", types.PrettyPrintStruct(txResp))
	}
}
