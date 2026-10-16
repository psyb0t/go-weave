# Go code

Refines `~/.claude/rule-details/working-with-generated-code.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `client.gen.go` and `types.gen.go` are oapi-codegen output. Never edit them by hand. Change `config.yaml`, `config.types.yaml` or a script under `scripts/`, then run `make generate` and commit the result.
- **GO2** The spec needs no patch today. If a future spec breaks oapi-codegen, add a jq patch script in `scripts/`, run it from a `//go:generate` line in `gen.go` between `fetch-spec.sh` and oapi-codegen, and point both configs at its output file.
- **GO3** The spec is never committed. `api/` is gitignored; `scripts/fetch-spec.sh` resolves the current URL from the `.stats.yml` of wandb/weave-stainless on every run. The file is JSON even though the URL ends in `.yml`.
- **GO4** `config.yaml` (client) and `config.types.yaml` (models) keep the same `output-options` and `compatibility`, or the two files disagree on type names. Both map `number` to `float64`.
- **GO5** Hand-written Go is limited to `gen.go` and the tests. The package never logs and never prints.
- **GO6** `make lint` passes with zero issues. Shell scripts pass shellcheck and shfmt. `make sec` excludes only the semgrep rule for `json.Unmarshal` into `interface{}`, which the spec's untyped schemas need.
