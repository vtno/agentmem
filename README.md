# agentmem

**Vectorless memory for AI agents.**

Agents forget. Worse: when memory lives inside a model provider or a single harness, your conventions, decisions, and project knowledge leave with the product.

`agentmem` is a pluggable memory layer you own. Plain markdown on your machine, one `memory` binary, wired into any agent that can run shell or MCP — so the memory stays yours when you switch tools.

**Free for personal and commercial use. Source is private.** This repository holds the documentation site, issue tracker, and release binaries.

> 🌐 อ่านภาษาไทย [ที่นี่](README.th.md)

## Install

```bash
curl -fsSL https://agentmem.thamtech.co/install | sh
```

Or download the latest binary from [GitHub Releases](https://github.com/vtno/agentmem/releases).

### Wire a harness

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

Other / skill-capable harnesses (portable skill only — no harness-specific integration):

```bash
memory install skill
# optional custom skills root:
memory install skill --dir /path/to/skills
```

Default skill path: `~/.agents/skills/memory/SKILL.md` (shared with OpenCode’s skill root).

Check status: `memory targets`.

## Why

- **Yours** — markdown files on disk, not locked in a provider account
- **Pluggable** — skill or MCP into the harness you already use
- **Portable** — switch agents; keep the same memory folder
- **Vectorless** — no embeddings, no daemon, no cloud required

## Docs site

A small Go server embeds the entire [`site/`](site/) tree (`//go:embed`) into one portable binary — templates, CSS, JS, assets, language packs, and Markdown.

### i18n

All UI copy lives in YAML under `site/lang/`. Root keys are **page names** (`common`, `index`, `install`, `commands`, `integrations`, `how_it_works`).

```yaml
# site/lang/en.yaml
common:
  nav:
    install: Install
index:
  title: "agentmem — vectorless memory for AI agents"
  tagline: Vectorless memory for AI agents
```

Templates call `{{.T "key"}}` (page-local, then `common`) or absolute paths `{{.T "common.nav.install"}}`.

Language selection (first match wins):

1. `?lang=en`
2. Cookie `lang=en`
3. `Accept-Language`
4. Default `en`

Add a locale by copying `site/lang/en.yaml` → `site/lang/<code>.yaml` and translating values.

| Route | |
|-------|--|
| `/` | Landing |
| `/docs/install` | Binary + harness wiring |
| `/docs/commands` | CLI reference |
| `/docs/integrations` | Claude Code, OpenCode, other harnesses |
| `/docs/how-it-works` | Storage, index, local UI |
| `/*.md` / `/docs/*.md` | Same pages as Markdown (`Content-Type: text/markdown`) |

### Makefile

```bash
make build                 # → ./agentmem-site
make install               # → ~/.local/bin/agentmem-site
make run                   # go run . -addr :5555  (ADDR=:8080 to override)
make test
make fmt
VERSION=0.1.0 make build   # stamp version into binary
./agentmem-site -version
```

### Manual build / run

```bash
# preview (rebuilds embed from site/)
go run . -addr :5555

# portable binary (no site/ dir needed at runtime)
make build
./agentmem-site -addr :5555
# ship only agentmem-site to the server
```

Layout: `site/templates/` (`layout`, `header`, `footer` + page bodies).  
Static: `site/css`, `site/js`, `site/assets`. Markdown: `site/index.md`, `site/docs/*.md`.  
Content changes require a rebuild so they are re-embedded.

## Core commands

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
memory ui                     # local loopback web UI
memory install <harness>      # Claude Code: skill + CLAUDE.md; OpenCode: skill + plugin
memory install --mcp <harness>
memory install skill [--dir <dir>]
memory uninstall …
memory targets                # harness + skill install status
```

## How it works

Everything lives in `~/.agentmem/` by default (override with `AGENTMEM_DIR`). Each memory is a markdown file. `MEMORY.md` is a generated index from each file’s summary line. Files are the source of truth.

## Supported harnesses

| Target | What you get |
|--------|----------------|
| **Claude Code** (official) | Skill + managed `CLAUDE.md` guidance · optional MCP |
| **OpenCode** (official) | Skill + plugin · optional MCP |
| **skill** (portable) | Skill only under `~/.agents/skills` — no harness-specific integration |

## Issues & releases

- Bugs and features: [GitHub Issues](https://github.com/vtno/agentmem/issues)
- Binaries: [GitHub Releases](https://github.com/vtno/agentmem/releases)

## License

See [LICENSE](LICENSE). Free to use; source is not open.
