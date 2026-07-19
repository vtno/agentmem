# agentmem

**Vectorless memory for AI agents.**

Agents forget. Worse: when memory lives inside a model provider or a single harness, your conventions, decisions, and project knowledge leave with the product.

`agentmem` is a pluggable memory layer you own. Plain markdown on your machine, one `memory` binary, wired into any agent that can run shell or MCP — so the memory stays yours when you switch tools.

## Install

```bash
curl -fsSL https://agentmem.thamtech.co/install | sh
```

See [Install](/docs/install) · [Releases](https://github.com/vtno/agentmem/releases) · [Issues](https://github.com/vtno/agentmem/issues)

## Why

- **Yours** — markdown files on disk, not locked in a provider account
- **Pluggable** — skill or MCP install into the harness you already use
- **Portable** — switch agents; keep the same memory folder
- **Vectorless** — no embeddings, no daemon, no cloud required

## Quick start

1. Install the `memory` binary (curl script or Releases)
2. Plug in: `memory install claude-code`, `opencode`, or `memory install skill` for other harnesses
3. Agents read and write the same local store — your memory travels with you

## Supported harnesses

| Target | What you get |
|--------|----------------|
| [Claude Code](/docs/integrations#claude-code) (official) | Skill + SessionStart hook · optional MCP |
| [OpenCode](/docs/integrations#opencode) (official) | Skill + plugin · optional MCP |
| [Other harnesses](/docs/integrations#other) | Skill only — `memory install skill` |

## Docs

- [Install](/docs/install)
- [Commands](/docs/commands)
- [Integrations](/docs/integrations)
- [How it works](/docs/how-it-works)

Markdown variants: append `.md` to any path (e.g. [docs/install.md](/docs/install.md)).

---

Free for personal and commercial use. Source private.  
[github.com/vtno/agentmem](https://github.com/vtno/agentmem)
