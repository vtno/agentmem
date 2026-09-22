# Install

Get the `memory` binary on your PATH, then connect a harness. Claude Code gets a skill and managed `CLAUDE.md` guidance; OpenCode gets a skill and plugin. Other agents use the portable skill.

[← Home](/) · [Commands](/docs/commands) · [Integrations](/docs/integrations)

## Recommended

One-liner (macOS / Linux):

```bash
curl -fsSL https://agentmem.thamtech.co/install | sh
```

The script detects your OS and architecture, downloads the matching release from [GitHub Releases](https://github.com/vtno/agentmem/releases), verifies its SHA-256 checksum, and installs `memory` to `~/.local/bin`. Set `INSTALL_DIR` to choose another location.

## Platforms

- macOS arm64 and amd64
- Linux amd64 and arm64
- Windows amd64 (zip from Releases; put `memory.exe` on PATH)

## Manual install

1. Open [GitHub Releases](https://github.com/vtno/agentmem/releases) and download the archive for your platform
2. Extract the `memory` binary
3. `chmod +x memory` and move it onto your PATH
4. Verify with `memory --help` or `memory targets`

Checksums ship as `SHA256SUMS` on each release. Source install via `go install` is not offered — the source repository is private.

## Wire your agent

Official harnesses:

```bash
memory install claude-code
memory install opencode
```

MCP (also registers the MCP server; the normal guidance remains installed):

```bash
memory install --mcp claude-code
memory install --mcp opencode
```

Other harnesses — portable skill only under `~/.agents/skills` (no Claude Code guidance or OpenCode plugin; `skill` is not a harness name):

```bash
memory install skill
memory install skill --dir /path/to/skills
```

Remove wiring later:

```bash
memory uninstall claude-code
memory uninstall --mcp opencode
memory uninstall skill
```

`memory targets` shows install status. See [Integrations](/docs/integrations) for what each mode writes.

Note: OpenCode’s skill path is the same default as `install skill` (`~/.agents/skills`) — intentional for compatibility.

## Configuration

- Default memory directory: `~/.agentmem/`
- Override with `AGENTMEM_DIR`

---

Free for personal and commercial use. Source private.
