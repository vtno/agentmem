# agentmem

**Dead simple memory for coding agents.**

`agentmem` gives your agent a tiny, durable memory folder it can read and update with plain shell commands. No database. No Docker. No daemon. Just markdown files and one `memory` binary.

## Why

Agents forget useful things: project conventions, user preferences, gotchas, decisions, and feedback. `agentmem` makes remembering boring:

```bash
memory list
memory show project/conventions
memory add --file project/conventions --summary "[project] Go conventions" "Wrap errors. Keep comments rare."
```

## Dead Simple Integration

Install the binary, then wire it into your agent:

```bash
memory install claude-code
# or
memory install opencode
```

That's it. The installer drops in the right skill/plugin files so the agent knows to check memory at the start of real tasks and save new facts when work is done.

For agents that prefer MCP:

```bash
memory install --mcp claude-code
memory install --mcp opencode
```

## How It Works

Everything lives in `~/.agentmem/` by default:

```text
~/.agentmem/
  MEMORY.md
  notes.md
  project/
    conventions.md
  user/
    prefs.md
```

Each memory is a markdown file. `MEMORY.md` is a generated index from each file's summary line, so agents can quickly browse what exists before opening details.

## Core Commands

```bash
memory list [prefix]          # show the memory index, optionally scoped
memory show <name>            # read one memory
memory search <query>         # grep across memories
memory add <text>             # append a quick note
memory add --file <name> --summary "[tag] desc" <text>
memory edit <name>            # open in $EDITOR
memory rm <name>              # delete a memory
memory reindex                # rebuild MEMORY.md
memory serve                  # MCP server over stdio
memory targets                # supported harnesses + install status
```

Names can use `/` for folders, e.g. `project/api` or `user/prefs`. On Windows, `\` is accepted and normalized.

## Install

```bash
curl -fsSL https://agentmem.thamtech.co/install | sh
```

Then wire it into your agent:

```bash
memory install claude-code
# or
memory install opencode
```

Prefer manual install? Download the latest binary from GitHub Releases.

## Philosophy

- **Portable**: plain markdown, easy to inspect and edit
- **Small**: one Go binary, fast startup
- **Agent-friendly**: index first, files as source of truth
- **No magic**: no vectors, no sync, no LLM calls inside the CLI

## Status

Early but usable. Supports Claude Code and OpenCode today.
