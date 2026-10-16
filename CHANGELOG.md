# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-16

- Added `weave`, a REST client generated with oapi-codegen from the OpenAPI document published with W&B's Weave Stainless SDK (W&B Weave trace server, spec version 0.1.0, OpenAPI 3.1.0). It covers all 74 paths: calls, objects, tables, files, feedback, costs, refs, threads, annotation queues, evaluations, OTel trace export, inference routing and the v2 entity/project resources (datasets, evaluations, evaluation runs, models, ops, predictions, scorers, scores).
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- The generated code is split into `types.gen.go` (models) and `client.gen.go` (client).
- `make generate` reads the current spec URL from the `.stats.yml` of wandb/weave-stainless, downloads the spec and regenerates the client. The spec is not committed.
- Added the `go-weave` agent skill: `claude plugin install go-weave@psyb0t` for Claude Code, `codex plugin add go-weave@psyb0t` for Codex, `openclaw skills install @psyb0t/go-weave` from ClawHub.
