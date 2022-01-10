package main

import (
	"context"
	"log"

	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
)

const (
	// serverURL is the URL of a Rosetta Server.
	serverURL = "http://localhost:8080"
)

func main() {
	config.Initialize()
	ctx := context.Background()

	// Step 1: Create a new fetcher
	newFetcher := fetcher.New(
		serverURL,
	)

	// Step 2: Initialize the fetcher's asserter
	//
	// Behind the scenes this makes a call to get the
	// network status and uses the response to inform
	// the asserter what are valid responses.
	primaryNetwork, networkStatus, err := newFetcher.InitializeAsserter(ctx, nil, "")
	if err != nil {
		log.Fatal(err)
	}

	// Step 3: Print the primary network and network status
	log.Printf("Primary Network: %s\n", types.PrettyPrintStruct(primaryNetwork))
	log.Printf("Network Status: %s\n", types.PrettyPrintStruct(networkStatus))

	// Step 4: Fetch the current block with retries (automatically
	// asserted for correctness)
	//
	// It is important to note that this assertion only ensures
	// required fields are populated and that operations
	// in the block only use types and statuses that were
	// provided in the networkStatusResponse. To run more
	// intensive validation, use the Rosetta Validator. It
	// can be found at: https://github.com/coinbase/rosetta-validator
	//
	// On another note, notice that fetcher.BlockRetry
	// automatically fetches all transactions that are
	// returned in BlockResponse.OtherTransactions. If you use
	// the client directly, you will need to implement a mechanism
	// to fully populate the block by fetching all these
	// transactions.
	block, err := newFetcher.BlockRetry(
		ctx,
		primaryNetwork,
		types.ConstructPartialBlockIdentifier(
			networkStatus.CurrentBlockIdentifier,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Step 5: Print the block
	log.Printf("Current Block: %s\n", types.PrettyPrintStruct(block))
}
