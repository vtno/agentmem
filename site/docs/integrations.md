# Integrations

Official harnesses get skill plus enforcement (hooks/plugins). Other shell-capable agents can still use the shared skill — without that enforcement layer.

[← Home](/) · [Install](/docs/install) · [Commands](/docs/commands) · [How it works](/docs/how-it-works)

## Harnesses

| Target | Notes |
|--------|--------|
| [Claude Code](#claude-code) | Official — skill + SessionStart hook · optional MCP |
| [OpenCode](#opencode) | Official — skill + plugin · optional MCP |
| [Other harnesses](#other-harnesses) | Skill only — no hook/plugin enforcement |

## Claude Code · official {#claude-code}

Default install writes the memory skill and a SessionStart hook so the agent is reminded to use memory across resume, clear, and compaction.

```bash
memory install claude-code
memory install --mcp claude-code
```

Skill path suits shell-capable sessions. MCP suits sandboxed or MCP-first setups. The SessionStart hook is installed in both modes.

## OpenCode · official {#opencode}

Default install writes the skill and an OpenCode plugin for session hooks. Use `--mcp` when you want the MCP server registered in OpenCode config.

```bash
memory install opencode
memory install --mcp opencode
```

## Other harnesses · skill only {#other}

Any agent that can load a skill and run shell commands can use `agentmem` against the same local memory folder. There is no dedicated install target for these harnesses — no SessionStart hook, no plugin, no automatic nudge to open or update memory.

Install the skill (and only the skill) automatically:

```bash
memory install skill
# override skill directory (default: ~/.agents/skills):
memory install skill --dir /path/to/skills
```

`skill` is not a harness name — it installs the portable skill only. Default path is under `~/.agents` for maximum compatibility: `~/.agents/skills/memory/SKILL.md`. Use `--dir` to point at another skills root. Point your harness at that skill, keep `memory` on `PATH`, and you’re done.

- No hook or plugin is installed — the agent is told *how* to use memory, not forced to
- Not combinable with a harness or `--mcp` (e.g. `install skill claude-code` errors)
- Optional: register `memory serve` as MCP if the harness is sandboxed
- Remove later with `memory uninstall skill` (same `--dir` if you used one)

Prefer full official installs when you can — hooks and plugins keep memory discipline consistent across resume, clear, and compaction. Default path matches OpenCode’s skill root; `uninstall skill` removes that shared skill file (an OpenCode plugin may remain).

## Skill vs MCP

- **Skill (default on harness)** — agent runs shell `memory …` commands. Best when the harness has a real shell.
- **MCP (`--mcp`)** — agent calls tools over the MCP server (`memory serve`). Better for sandboxed or shell-less harnesses.
- **Hooks / plugins** — official harnesses only. They enforce “check memory / save facts”; portable `install skill` does not.

Check status anytime with `memory targets`. Remove with `memory uninstall [--mcp] <harness>` or `memory uninstall skill`.

---

Free for personal and commercial use. Source private.
