# Go code

Refines `~/.claude/rule-details/working-with-generated-code.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `client.gen.go` and `types.gen.go` are oapi-codegen output. Never edit them by hand. Change `config.yaml`, `config.types.yaml` or a script under `scripts/`, then run `make generate` and commit the result.
- **GO2** Spec fixes live in `scripts/patch-spec.sh`, a jq script that `gen.go` runs between `fetch-spec.sh` and oapi-codegen; both oapi-codegen steps read its output `api/openapi.patched.json`. Every patch states why it exists and fails loudly when the spec no longer matches it. Drop a patch once upstream fixes the spec.
- **GO3** The spec is never committed. `api/` is gitignored; `scripts/fetch-spec.sh` downloads it from `https://trace.wandb.ai/openapi.json` on every run.
- **GO4** `config.yaml` (client) and `config.types.yaml` (models) keep the same `output-options` and `compatibility`, or the two files disagree on type names. Both map `number` to `float64`.
- **GO5** Hand-written Go is limited to `gen.go` and the tests. The package never logs and never prints.
- **GO6** `make lint` passes with zero issues. Shell scripts pass shellcheck and shfmt. `make sec` excludes only the semgrep rule for `json.Unmarshal` into `interface{}`, which the spec's untyped schemas need.
