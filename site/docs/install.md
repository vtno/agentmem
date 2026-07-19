# Install

Get the `memory` binary on your PATH, then plug in a harness. Official targets add enforcement (hooks/plugins); other agents use the portable skill.

[← Home](/) · [Commands](/docs/commands) · [Integrations](/docs/integrations)

## Recommended

One-liner (macOS / Linux):

```bash
curl -fsSL https://agentmem.thamtech.co/install | sh
```

The script detects your OS and architecture, downloads the matching release from [GitHub Releases](https://github.com/vtno/agentmem/releases), and installs `memory` to `~/.local/bin` (or `/usr/local/bin` when available).

## Platforms

- macOS arm64 and amd64
- Linux amd64 and arm64
- Windows amd64 (zip from Releases; put `memory.exe` on PATH)

## Manual install

1. Open [GitHub Releases](https://github.com/vtno/agentmem/releases) and download the archive for your platform
2. Extract the `memory` binary
3. `chmod +x memory` and move it onto your PATH
4. Verify with `memory --help` or `memory targets`

Checksums ship as `checksums.txt` on each release. Source install via `go install` is not offered — the source repository is private.

## Wire your agent

Official harnesses (skill + enforcement):

```bash
memory install claude-code
memory install opencode
```

MCP (hooks/plugins still installed on official targets):

```bash
memory install --mcp claude-code
memory install --mcp opencode
```

Other harnesses — portable skill only under `~/.agents/skills` (no hooks/plugins; `skill` is not a harness name):

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
