# go-weave

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-weave.svg)](https://pkg.go.dev/github.com/psyb0t/go-weave)
[![CI](https://github.com/psyb0t/go-weave/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-weave/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-weave/badges/version.svg)](https://github.com/psyb0t/go-weave/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-weave/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-weave/badges/importers.svg)](https://github.com/psyb0t/go-weave/blob/badges/importers.md)

Unofficial Go client for the [W&B Weave](https://weave-docs.wandb.ai) trace server API, generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document that ships with W&B's [Weave Stainless SDK](https://github.com/wandb/weave-stainless). I am not affiliated with Weights & Biases.

The package is `weave`, in the repository root. Everything in `client.gen.go` and `types.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-weave
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-weave"
)

func main() {
	auth := weave.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.SetBasicAuth(os.Getenv("WANDB_USERNAME"), os.Getenv("WANDB_API_KEY"))

		return nil
	})

	client, err := weave.NewClientWithResponses("https://trace.wandb.ai", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ServerInfoServerInfoGetWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

Method names come from the spec's operation IDs, which FastAPI builds from the handler name, the path and the HTTP method. So they are long and repeat themselves: `CallReadCallReadPost` is `POST /call/read`, `ObjReadObjReadPost` is `POST /obj/read`, `ReadRootHealthGet` is `GET /health`. Most endpoints take a `project_id` in the request body, in the form `<entity>/<project>`.

## Auth

The spec declares HTTP Basic and HTTP Bearer. The official SDK uses Basic: the W&B username (`WANDB_USERNAME`) is the username and the W&B API key (`WANDB_API_KEY`) is the password. The API key comes from the user settings page on wandb.ai. For a dedicated or self-hosted W&B instance, swap `https://trace.wandb.ai` for that instance's trace server URL.

## Agent integrations

The [skill](.agents/skills/go-weave) teaches a coding agent how to use go-weave when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-weave@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-weave@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-weave
```

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` reads `openapi_spec_url` from the `.stats.yml` of [wandb/weave-stainless](https://github.com/wandb/weave-stainless) and downloads that URL to `api/openapi.json`. Stainless hosts the spec under a URL with a content hash, which changes on every spec update, so the script looks it up each time. The file is JSON even though the URL ends in `.yml`.
- oapi-codegen runs twice. `config.types.yaml` writes the models to `types.gen.go` and `config.yaml` writes the client to `client.gen.go`, because one file would be about 850 KB. Both map `number` to `float64`, because oapi-codegen would otherwise use `float32`.

The spec needs no patching. It is OpenAPI 3.1, and oapi-codegen turns its `anyOf` with `null` fields into pointers.

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

`make sec` skips one semgrep rule, the one that flags `json.Unmarshal` into `interface{}`. oapi-codegen emits that for the spec's untyped schemas, where any JSON value is valid.

## License

[MIT](LICENSE). The generated code comes from the OpenAPI specification W&B publishes with its Apache-2.0 licensed Weave Stainless SDK. The specification itself is not redistributed here.
