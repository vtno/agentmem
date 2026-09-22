#!/usr/bin/env bash
# Smoke scenario for a cross-compiled `memory` binary (any GOOS/GOARCH).
#
# Usage: scripts/smoke.sh /path/to/memory
#
# Covers: version, add (arg / stdin / quick note), list, show, search,
# reindex, serve (MCP initialize over stdio), ui (loopback HTTP API),
# rm, install skill + uninstall skill (portable, into a temp dir), targets.
#
# Uses temp AGENTMEM_DIR and temp skill dir; never touches host state.
set -euo pipefail

bin=${1:?usage: smoke.sh <path-to-memory-binary>}
[ -e "$bin" ] || { echo "binary not found: $bin" >&2; exit 1; }

work=$(mktemp -d)
export AGENTMEM_DIR=$(mktemp -d)
ui_pid=""
cleanup() {
	[ -n "$ui_pid" ] && kill "$ui_pid" 2>/dev/null
	rm -rf "$work" "$AGENTMEM_DIR"
}
trap cleanup EXIT

step() {
	printf '\n=== %s\n' "$*"
}

step "version"
"$bin" version

step "add (--file arg)"
"$bin" add --file smoke --summary "[smoke] ci round-trip" "hello from smoke"

step "add (stdin)"
printf 'piped body\n' | "$bin" add --file piped --summary "[smoke] from stdin"

step "add (quick note)"
"$bin" add "a quick note line"

step "list"
"$bin" list | grep -F "smoke"
"$bin" list | grep -F "piped"

step "show"
"$bin" show smoke | grep -F "hello from smoke"
"$bin" show piped | grep -F "piped body"

step "search"
"$bin" search "piped body" | grep -F "piped"

step "reindex"
"$bin" reindex
"$bin" list | grep -F "smoke"

step "serve (MCP initialize over stdio)"
mcp_out=$( {
	printf '%s\n' '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke","version":"0"}}}'
	sleep 5
} | "$bin" serve)
printf '%s' "$mcp_out" | grep -F "serverInfo"

step "ui (loopback API)"
addr=127.0.0.1:8791
"$bin" ui --addr "$addr" --no-open >/dev/null 2>&1 &
ui_pid=$!
up=
for _ in $(seq 1 100); do
	if curl -sf "http://$addr/api/meta" >/dev/null 2>&1; then
		up=1
		break
	fi
	sleep 0.2
done
[ -n "$up" ] || { echo "ui did not come up on $addr" >&2; exit 1; }
curl -sf "http://$addr/api/meta" | grep -F '"version"'
curl -sf "http://$addr/api/list" | grep -F '"smoke"'
curl -sf "http://$addr/api/show?name=smoke" | grep -F "hello from smoke"
curl -sf "http://$addr/api/search?q=piped" | grep -F '"piped"'
curl -sf "http://$addr/" | grep -F "agentmem"
kill "$ui_pid" 2>/dev/null
wait "$ui_pid" 2>/dev/null || true
ui_pid=""

step "rm"
"$bin" rm smoke
if "$bin" list | grep -E '^[[:space:]]+smoke[[:space:]]' >/dev/null; then
	echo "smoke still in index after rm" >&2
	exit 1
fi
"$bin" rm piped
"$bin" list | grep -E '^[[:space:]]+piped[[:space:]]' >/dev/null && {
	echo "piped still in index after rm" >&2
	exit 1
}

step "install skill (portable, temp dir)"
skilldir="$work/skills"
"$bin" install skill --dir "$skilldir"
[ -f "$skilldir/memory/SKILL.md" ] || {
	echo "SKILL.md missing after install" >&2
	exit 1
}

step "uninstall skill (portable)"
"$bin" uninstall skill --dir "$skilldir"
[ ! -e "$skilldir/memory" ] || {
	echo "skill dir still present after uninstall" >&2
	exit 1
}

step "targets"
"$bin" targets >/dev/null

printf '\nSMOKE OK: %s\n' "$bin"
