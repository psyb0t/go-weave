// Package weave is a Go client for the W&B Weave trace server API, generated
// with oapi-codegen from the OpenAPI document published with the W&B Weave
// Stainless SDK.
package weave

//go:generate bash scripts/fetch-spec.sh
//go:generate go tool oapi-codegen --config=config.types.yaml api/openapi.json
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.json
