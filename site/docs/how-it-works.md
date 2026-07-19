# How it works

Everything is a folder of markdown files you own. The CLI is a thin tool over that folder — no database, no cloud, no embedding model. Harnesses plug in via skill, hooks, or MCP; the store stays the same when you switch tools.

[← Home](/) · [Integrations](/docs/integrations) · [Commands](/docs/commands)

## Storage

Default directory is `~/.agentmem/` (override with `AGENTMEM_DIR`):

```text
~/.agentmem/
  MEMORY.md              ← generated index
  notes.md
  project/
    conventions.md
  user/
    prefs.md
```

Names may contain `/` to group related memories. The on-disk tree is organization; the logical index always lists the full set of entries.

## Memory files

Each entry is plain markdown. An optional first-line summary blockquote is the source for the index. Everything after the following blank line is the body.

```markdown
> [user] prefers TypeScript, dark mode, short responses

Prefers TS over JS for all new work. Wants terse replies — no preamble.
Dark mode everywhere.
```

Files are the single source of truth. There is no YAML frontmatter. Tags like `[user]` or `[project]` ride in the summary line.

## Index

`MEMORY.md` is **derived** from each file’s summary line so agents can scan what exists before opening details. After editing files by hand, run `memory reindex`.

Writes through the CLI, MCP, or `memory ui` call `Store.Write`, which reindexes automatically.

## Local UI

`memory ui` opens a loopback-only browser UI (default `127.0.0.1:8787`) to list, search, create, edit, and delete entries. Same storage as the CLI — mutations go through the same store layer.

```bash
memory ui
memory ui --no-open
```

## Philosophy

- **Yours** — files on disk, not locked to a provider or single harness
- **Pluggable** — skill, MCP, and official hooks/plugins; same memory folder
- **Portable** — plain markdown, easy to inspect and edit
- **Small** — one binary, fast startup (important when hooks call the CLI often)
- **Vectorless** — no embeddings, no sync daemon, no LLM inside the CLI

---

Free for personal and commercial use. Source private.
