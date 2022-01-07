package client

import (
	ctx "context"
	"github.com/coinbase/rosetta-sdk-go/client"
	"net/http"
	"time"
)

const (
	// serverURL is the URL of a Rosetta Server.
	serverURL = "http://localhost:8080"

	// agent is the user-agent on requests to the
	// Rosetta Server.
	agent = "rosetta-sdk-go"

	// defaultTimeout is the default timeout for
	// HTTP requests.
	defaultTimeout = 10 * time.Second
)

// Step 1: Create a client
func create_test_client() *client.APIClient {
	clientCfg := client.NewConfiguration(
		serverURL,
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
