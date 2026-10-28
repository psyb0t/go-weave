// Package weave is a Go client for the W&B Weave trace server API, generated
// with oapi-codegen from the OpenAPI document the trace server publishes.
package weave

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.types.yaml api/openapi.patched.json
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.patched.json
