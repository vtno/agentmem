# Integrations

Claude Code and OpenCode receive native integration in addition to the skill. Other shell-capable agents can use the shared skill against the same local memory folder.

[← Home](/) · [Install](/docs/install) · [Commands](/docs/commands) · [How it works](/docs/how-it-works)

## Harnesses

| Target | Notes |
|--------|--------|
| [Claude Code](#claude-code) | Official — skill + `CLAUDE.md` guidance · optional MCP |
| [OpenCode](#opencode) | Official — skill + plugin · optional MCP |
| [Other harnesses](#other-harnesses) | Skill only — no harness-specific integration |

## Claude Code · official {#claude-code}

Default install writes the memory skill and manages an agentmem section in your user-level `~/.claude/CLAUDE.md`, so Claude Code has the memory instructions in every session.

```bash
memory install claude-code
memory install --mcp claude-code
```

Skill suits shell-capable sessions. MCP suits sandboxed or MCP-first setups. Both modes keep the managed `CLAUDE.md` guidance; `--mcp` also registers `memory serve` in Claude settings.

## OpenCode · official {#opencode}

Default install writes the skill and an OpenCode plugin for session hooks. Use `--mcp` when you want the MCP server registered in OpenCode config.

```bash
memory install opencode
memory install --mcp opencode
```

## Other harnesses · skill only {#other}

Any agent that can load a skill and run shell commands can use `agentmem` against the same local memory folder. No Claude Code guidance, OpenCode plugin, or MCP configuration is written automatically for these harnesses.

Install the skill (and only the skill) automatically:

```bash
memory install skill
# override skill directory (default: ~/.agents/skills):
memory install skill --dir /path/to/skills
```

`skill` is not a harness name — it installs the portable skill only. Default path is under `~/.agents` for maximum compatibility: `~/.agents/skills/memory/SKILL.md`. Use `--dir` to point at another skills root. Point your harness at that skill, keep `memory` on `PATH`, and you’re done.

- No harness-specific integration is installed — the skill tells the agent how to use memory
- Not combinable with a harness or `--mcp` (e.g. `install skill claude-code` errors)
- Optional: register `memory serve` as MCP if the harness is sandboxed
- Remove later with `memory uninstall skill` (same `--dir` if you used one)

Prefer the official installer when you use Claude Code or OpenCode — it adds the relevant guidance or plugin. Default path matches OpenCode’s skill root; `uninstall skill` removes that shared skill file (an OpenCode plugin may remain).

## Skill vs MCP

- **Skill (default on harness)** — agent runs shell `memory …` commands. Best when the harness has a real shell.
- **MCP (`--mcp`)** — agent calls tools over the MCP server (`memory serve`). Better for sandboxed or shell-less harnesses.
- **Native integration** — Claude Code gets managed `CLAUDE.md` guidance; OpenCode gets a plugin. Portable `install skill` does not add either.

Check status anytime with `memory targets`. Remove with `memory uninstall [--mcp] <harness>` or `memory uninstall skill`.

---

Free for personal and commercial use. Source private.
