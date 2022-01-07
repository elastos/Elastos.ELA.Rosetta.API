package client

import (
	ctx "context"
	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/elastos/Elastos.ELA.Rosetta.API/common/config"
	"net/http"
	"strconv"
	"time"
)

const (

	// agent is the user-agent on requests to the
	// Rosetta Server.
	agent = "rosetta-sdk-go"

	// defaultTimeout is the default timeout for
	// HTTP requests.
	defaultTimeout = 10 * time.Second
)

func init() {
	config.Initialize()
}

// Step 1: Create a client
func create_test_client() *client.APIClient {
	clientCfg := client.NewConfiguration(
		config.Parameters.MainNodeRPC.IpAddress+":"+strconv.Itoa(config.Parameters.ServerPort),
		agent,
		&http.Client{
			Timeout: defaultTimeout,
		},
	)
	return client.NewAPIClient(clientCfg)
}

func context() ctx.Context {
	return ctx.Background()
}
