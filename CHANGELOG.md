# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.0, 2026-10-28

- The client is generated from the spec the W&B Weave trace server publishes at `https://trace.wandb.ai/openapi.json`. The Stainless-hosted copy of the spec it used before is no longer publicly readable, so scheduled regeneration had stopped.
- Breaking: the published spec adds 74 endpoint methods and removes the `/inference/v1/{path}` proxy operations (the `InferenceGet*` and `InferencePost*` methods). `Expr` is now `Query_Expr`, and `Operation` is split into inline union types. The v2 list endpoints for datasets, models, ops, scorers and evaluations, `/threads/stream_query` and `/annotation_queues/query` return newline-delimited JSON, so their responses no longer have a `JSON200` field.
- `scripts/patch-spec.sh` declares `application/jsonl` for `POST /calls/stream_query`, which the published spec lists as `application/json`. Without it the client would decode the stream as a single JSON value and fail on any stream longer than one line.
- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.1.3, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- The issue-pull workflow runs at hours derived from the repository name instead of the same four hours as every other repository.

## v0.1.2, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.1, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.0, 2026-10-16

- Added `weave`, a REST client generated with oapi-codegen from the OpenAPI document published with W&B's Weave Stainless SDK (W&B Weave trace server, spec version 0.1.0, OpenAPI 3.1.0). It covers all 74 paths: calls, objects, tables, files, feedback, costs, refs, threads, annotation queues, evaluations, OTel trace export, inference routing and the v2 entity/project resources (datasets, evaluations, evaluation runs, models, ops, predictions, scorers, scores).
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- The generated code is split into `types.gen.go` (models) and `client.gen.go` (client).
- `make generate` reads the current spec URL from the `.stats.yml` of wandb/weave-stainless, downloads the spec and regenerates the client. The spec is not committed.
- Added the `go-weave` agent skill: `claude plugin install go-weave@psyb0t` for Claude Code, `codex plugin add go-weave@psyb0t` for Codex, `openclaw skills install @psyb0t/go-weave` from ClawHub.
