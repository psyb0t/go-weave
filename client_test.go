package weave_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-weave"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUsername  = "test-user"
	testAPIKey    = "test-key"
	testProjectID = "entity-x/project-x"
	testCallID    = "call-x"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *weave.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := weave.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.SetBasicAuth(testUsername, testAPIKey)

			return nil
		},
	)

	client, err := weave.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertBasicAuth(t *testing.T, r *http.Request) {
	t.Helper()

	user, pass, ok := r.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, testUsername, user)
	assert.Equal(t, testAPIKey, pass)
}

func TestServerInfo(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/server_info", r.URL.Path)
		assertBasicAuth(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"min_required_weave_python_version":"1.0.0",` +
				`"trace_server_version":"2.0.0"}`,
		))
	})

	res, err := client.ServerInfoServerInfoGetWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, "1.0.0", res.JSON200.MinRequiredWeavePythonVersion)
	assert.Equal(t, "2.0.0", res.JSON200.TraceServerVersion)
}

func TestCallRead(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/call/read", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBasicAuth(t, r)

		var body weave.CallReadReq

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, testProjectID, body.ProjectId)
		assert.Equal(t, testCallID, body.Id)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"call":{"id":"` + testCallID + `",` +
				`"project_id":"` + testProjectID + `",` +
				`"op_name":"op-x","trace_id":"trace-x",` +
				`"started_at":"2026-01-02T03:04:05Z",` +
				`"attributes":{},"inputs":{"n":1}}}`,
		))
	})

	res, err := client.CallReadCallReadPostWithResponse(
		context.Background(),
		weave.CallReadReq{ProjectId: testProjectID, Id: testCallID},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.Call)

	call := res.JSON200.Call
	assert.Equal(t, testCallID, call.Id)
	assert.Equal(t, "op-x", call.OpName)
	assert.InDelta(t, 1.0, call.Inputs["n"], 0)
}
