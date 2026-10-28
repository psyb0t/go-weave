#!/bin/bash
# Fixes the response media type of POST /calls/stream_query in the downloaded
# W&B Weave spec. The server streams JSONL there, like the other stream
# endpoints, but the spec declares application/json. With that declaration the
# generated client runs json.Unmarshal on the whole JSONL body and fails as
# soon as the stream has more than one line. Declaring application/jsonl leaves
# the body untyped in Body, as for the other stream endpoints. The script fails
# when the spec no longer has the declaration it fixes, so a stale patch shows.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

readonly SPEC_IN="openapi.json"
readonly SPEC_OUT="openapi.patched.json"
SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

readonly STREAM_QUERY_200='.paths["/calls/stream_query"].post.responses["200"].content'

main() {
	if [[ ! -f "${SPEC_DIR}/${SPEC_IN}" ]]; then
		log ERROR "spec not found, run scripts/fetch-spec.sh first spec=${SPEC_IN}"
		exit 1
	fi

	if ! jq -e "${STREAM_QUERY_200} | keys == [\"application/json\"]" \
		"${SPEC_DIR}/${SPEC_IN}" >/dev/null; then
		log ERROR "POST /calls/stream_query no longer declares only application/json; review and drop this patch"
		exit 1
	fi

	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.${SPEC_OUT}.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}'" EXIT

	jq "${STREAM_QUERY_200} = {\"application/jsonl\": {\"schema\": {}}}" \
		"${SPEC_DIR}/${SPEC_IN}" >"$tmp"
	chmod 0644 "$tmp"
	mv -- "$tmp" "${SPEC_DIR}/${SPEC_OUT}"
	log INFO "declared application/jsonl for POST /calls/stream_query spec=${SPEC_OUT}"
}

main "$@"
