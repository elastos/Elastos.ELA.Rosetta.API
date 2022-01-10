package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elastos/Elastos.ELA.Rosetta.API/common/base"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"github.com/elastos/Elastos.ELA.Rosetta.API/server/services"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(
	network *types.NetworkIdentifier,
	asserter *asserter.Asserter,
) http.Handler {
	networkAPIService := services.NewNetworkAPIService(network)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserter,
	)

	blockAPIService := services.NewBlockAPIService(network)
	blockAPIController := server.NewBlockAPIController(
		blockAPIService,
		asserter,
	)

	mempoolAPIService := services.NewMempoolAPIService(network)
	mempoolAPIController := server.NewMempoolAPIController(
		mempoolAPIService,
		asserter,
	)

	accountAPIService := services.NewAccounAPIService(network)
	accountAPIController := server.NewAccountAPIController(
		accountAPIService,
		asserter,
	)

	constructionAPIService := services.NewConstructionAPIServicer(network)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		asserter,
	)

	return server.NewRouter(
		networkAPIController,
		blockAPIController,
		mempoolAPIController,
		accountAPIController,
		constructionAPIController)
}

func main() {
	config.Initialize()
	network := &types.NetworkIdentifier{
		Blockchain: base.BlockChainName,
		Network:    config.Parameters.ActiveNet,
	}

	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserter, err := asserter.NewServer(
		[]string{base.MainnetNextworkType},
		false,
		[]*types.NetworkIdentifier{network},
		nil,
		false,
		"",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create the main router handler then apply the logger and Cors
	// middlewares in sequence.
	router := NewBlockchainRouter(network, asserter)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	log.Printf("Listening on port %d\n", config.Parameters.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Parameters.ServerPort), corsRouter))
}
