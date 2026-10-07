# Ink Mark — Agent Instructions

## About This Project

**Ink Mark** is a local-first, single-user markdown writing app.

- **Backend**: Go REST API (`localhost:18080`)
- **Frontend**: SvelteKit + TypeScript (`localhost:13000`)
- **Database**: PostgreSQL only (`localhost:15432` — user: `inkmark` / pass: `inkmark_dev` / db: `inkmark`)

Full-text search via PostgreSQL `tsvector` + GIN index. No sync, no multi-user.

---

## Rules

### General

- **Never commit** unless explicitly asked to commit, even if the user says they are satisfied with changes.
- **Simpler is better.** Always go for the simpler design that makes things easier to understand and maintain.
- **Follow the correct design principle for the system.** Do not overcomplicate it. Introduce a new design pattern only when complexity demands it and that pattern actually simplifies things.
- **Good code practices always.** Readable, maintainable, consistent code.
- **Healthy separation of concerns.** Enough, but not too much — don't over-engineer boundaries.
- **Follow SOLID principles.** Apply them where they genuinely improve the design, not as dogma.
- **Advice mode.** When the user asks for advice, provide advice only — no code changes.

### Languages

- **Go**: Strongly typed. Use explicit types, avoid `interface{}` / `any` unless truly necessary.
- **TypeScript**: Strongly typed. `strict` mode on. No implicit `any`. Use proper interfaces and type aliases.

### Git

- Use **Conventional Commits** format only: `feat`, `fix`, `docs`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.
- Format: `<type>[optional scope]: <description>`
- Only commit when explicitly asked.

---

## Skills

> Install with: `npx skills add <source> --skill <name>`

### Go — `samber/cc-skills-golang` ✅ All audits pass

| Skill | Purpose |
|---|---|
| `golang-code-style` | Style rules, line length, clarity over cleverness |
| `golang-error-handling` | Idiomatic Go error wrapping and propagation |
| `golang-database` | Database access patterns for Go |
| `golang-testing` | Table-driven tests, testify, coverage |

```sh
npx skills add samber/cc-skills-golang --skill golang-code-style
npx skills add samber/cc-skills-golang --skill golang-error-handling
npx skills add samber/cc-skills-golang --skill golang-database
npx skills add samber/cc-skills-golang --skill golang-testing
```

### Svelte / SvelteKit — `sveltejs/ai-tools` ✅ All audits pass (Official Svelte team)

| Skill | Purpose |
|---|---|
| `svelte-code-writer` | Writing correct Svelte 5 components with MCP doc access |
| `svelte-core-bestpractices` | `$state`, `$derived`, `$effect` runes best practices |

```sh
npx skills add sveltejs/ai-tools --skill svelte-code-writer
npx skills add sveltejs/ai-tools --skill svelte-core-bestpractices
```

### Svelte 5 — `ejirocodes/agent-skills` ✅ All audits pass

| Skill | Purpose |
|---|---|
| `svelte5-best-practices` | Runes, snippets, TypeScript props, SvelteKit load functions |

```sh
npx skills add ejirocodes/agent-skills --skill svelte5-best-practices
```

### TypeScript — `wshobson/agents` ✅ All audits pass

| Skill | Purpose |
|---|---|
| `typescript-advanced-types` | Generics, conditional types, mapped types, strict typing |

```sh
npx skills add wshobson/agents --skill typescript-advanced-types
```

### Git — `github/awesome-copilot` ✅ All audits pass

| Skill | Purpose |
|---|---|
| `git-commit` | Conventional Commits: feat, fix, chore, docs, refactor, etc. |

```sh
npx skills add github/awesome-copilot --skill git-commit
```
