package weave_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-weave"
)

// This is the README usage example. It talks to the live API, so it has no
// Output comment and only has to compile.
//
//nolint:testableexamples // needs real credentials, compile check only
func Example() {
	auth := weave.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.SetBasicAuth(
				os.Getenv("WANDB_USERNAME"),
				os.Getenv("WANDB_API_KEY"),
			)

			return nil
		},
	)

	client, err := weave.NewClientWithResponses(
		"https://trace.wandb.ai",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ServerInfoServerInfoGetWithResponse(
		context.Background(),
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
