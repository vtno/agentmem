# Commands

The binary is named `memory`. Entry names may use `/` for folders (e.g. `project/api`). On Windows, `\` is accepted and normalized.

[← Home](/) · [Install](/docs/install) · [Integrations](/docs/integrations)

## Browse

### `memory list [prefix]`

Print the derived index (equivalent to MEMORY.md when up to date). Optional prefix scopes the list (e.g. `project`).

```bash
memory list
memory list project
```

### `memory show <name>`

Print one memory file.

```bash
memory show project/conventions
```

### `memory search <query>`

Full-text search across memories. Returns matching files and snippets.

```bash
memory search conventions
```

## Write

### `memory add`

Append a quick note, or write a named file. Named files require `--summary`. Reads stdin when no text argument is given.

```bash
memory add "remember to run tests before push"
memory add --file project/conventions --summary "[project] Go conventions" "Wrap errors. Keep comments rare."
memory add --file user/prefs --summary "[user] short answers" <<'EOF'
Prefer terse replies. No preamble.
EOF
```

### `memory edit <name>`

Open a memory file in `$EDITOR`.

```bash
memory edit user/prefs
```

### `memory rm <name>`

Delete a memory file and its index entry.

```bash
memory rm scratch/old-idea
```

### `memory reindex`

Rebuild MEMORY.md from memory files after hand edits.

```bash
memory reindex
```

## Serve & UI

### `memory serve`

Start the MCP server on stdio (spawned by harnesses that use MCP).

```bash
memory serve
```

### `memory ui`

Local browser UI to list, search, create, edit, and delete memories. Loopback only. Save reindexes automatically via `Store.Write`.

```bash
memory ui
memory ui --addr 127.0.0.1:8787 --no-open
```

## Harness

### `memory install` / `uninstall`

Official harnesses: skill plus hooks/plugins by default; pass `--mcp` for MCP config (hooks/plugins still applied). Portable skill target installs skill only under `~/.agents/skills` (or `--dir`) — no harness name, no hooks.

```bash
memory install claude-code
memory install --mcp opencode
memory install skill
memory install skill --dir /path/to/skills
memory uninstall claude-code
memory uninstall skill
```

### `memory targets`

List supported harnesses and install status (skill, mcp, hook/plugin), plus the portable `skill` path.

```bash
memory targets
```

---

Free for personal and commercial use. Source private.
