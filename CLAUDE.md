# Ink Mark — Agent Instructions

## About This Project

**Ink Mark** is a local-first, single-user markdown writing app built with Wails v2 (Go backend + SvelteKit + TypeScript frontend).

- **Runtime**: Wails v2 — Go + SvelteKit rendered in a native webview
- **Frontend**: SvelteKit + TypeScript (`frontend/`)
- **Backend**: Go (`main.go`, `app.go`)

No server, no database, no sync, no multi-user.

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

### Go — Test-Driven Development

All Go code is written test-first. No production Go code without a corresponding test. This applies to every layer — repository, service, handler, utility.

#### TDD cycle

1. Write the test first — it must fail (red).
2. Write the minimum production code to make it pass (green).
3. Refactor. Keep tests passing.

Never write production code speculatively. If there is no test demanding the behaviour, the code should not exist.

#### What must be tested

Every testable behaviour needs a test. When writing or reviewing Go code, ask: *what can go wrong here, and what must go right?* Cover all of it.

**Success paths:**
- The happy path — expected input, expected output.
- Every distinct valid variant (e.g. different valid inputs that produce different valid outputs).
- Boundary values at the edges of valid ranges.

**Failure paths — every single one:**
- Invalid input (empty strings, zero values, nil, malformed data).
- Out-of-range values.
- Not-found cases (record does not exist).
- Duplicate / constraint violation (unique key, FK).
- Database errors (connection failure, query error — use mocks or test doubles).
- Partial failure (some records succeed, one fails).
- Context cancellation and timeout.
- Unexpected data types or shapes in `jsonb` columns.

**State and side effects:**
- Verify the database state after a write (insert/update/delete) by reading it back.
- Verify `updated_at` is actually updated on mutations.
- Verify triggers fire as expected.
- Verify indexes are used where performance depends on them (use `EXPLAIN` in integration tests where relevant).

#### Test structure

Use **table-driven tests** for all cases with multiple inputs. Each case must have a descriptive name.

```go
func TestCreateDocument(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateDocumentInput
        wantErr bool
        errType error // nil if not checking specific error
    }{
        {
            name:  "success — valid title and content",
            input: CreateDocumentInput{Title: "My Note", Content: `{}`},
        },
        {
            name:    "failure — empty title",
            input:   CreateDocumentInput{Title: "", Content: `{}`},
            wantErr: true,
        },
        {
            name:    "failure — invalid json content",
            input:   CreateDocumentInput{Title: "Note", Content: `not-json`},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // arrange, act, assert
        })
    }
}
```

#### Testing with the database

Use **`testify`** (`github.com/stretchr/testify`) — `assert` for non-fatal, `require` for fatal assertions.

For database tests, use a **real test database** — not mocks for the repository layer. Spin up a transaction per test and roll it back at the end, so tests are isolated and leave no state:

```go
func withTestTx(t *testing.T, db *sql.DB, fn func(tx *sql.Tx)) {
    t.Helper()
    tx, err := db.Begin()
    require.NoError(t, err)
    t.Cleanup(func() { _ = tx.Rollback() })
    fn(tx)
}
```

Mock at the **service layer** when testing business logic that depends on a repository interface — never mock the database itself at the repository layer.

#### Test file conventions

- Test files live next to the code they test: `document_repo_test.go` alongside `document_repo.go`.
- Integration tests that require a database are in `_integration_test.go` files and guarded with a build tag or `testing.Short()` skip:

```go
func TestDocumentRepo_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    // ...
}
```

- Run unit tests only: `go test -short ./...`
- Run all tests including integration: `go test ./...`

#### Assertions checklist per test

Every test must assert:
- [ ] The returned value is correct (or error is nil on success).
- [ ] The returned error is non-nil on failure paths.
- [ ] The error message or type is what is expected (where meaningful).
- [ ] Any side effects (DB state, timestamps, derived fields) are verified.
- [ ] No unexpected rows were created/modified/deleted.

### Languages

- **Go**: Strongly typed. Use explicit types, avoid `interface{}` / `any` unless truly necessary.
- **TypeScript**: Strongly typed. `strict` mode on. No implicit `any`. Use proper interfaces and type aliases.

### Git

- Use **Conventional Commits** format only: `feat`, `fix`, `docs`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.
- Format: `<type>[optional scope]: <description>`
- Only commit when explicitly asked.

---

## Design

### Rendering — HTML/CSS First

This is a Wails webview app — not a browser with framework overhead constraints. Always prefer native HTML and CSS over framework abstractions when they are equivalent:

- Use `<dialog>` over a custom modal component.
- Use `<details>`/`<summary>` over a custom accordion.
- Use CSS transitions/animations over Svelte's `transition:` directive unless Svelte adds real value (e.g. list reordering, coordinated animations).
- Use CSS Grid and Flexbox directly — no layout helper components.
- Use Svelte constructs (`{#if}`, `{#each}`, runes, actions) where they genuinely simplify reactivity or templating. Svelte is a highly optimised compiler — trust it, but only reach for it when it earns its place.

### Color Tokens

All colors are CSS custom properties defined in a single `:root` block. **Never hardcode a color value** — always reference a token. This makes the entire palette swappable in one place.

```css
:root {
  /* ── Surfaces ────────────────────────────────────────────── */
  --color-base: #f7f6f3; /* warm paper-white — main app background  */
  --color-surface: #fdfcfa; /* slightly lighter — panels, cards        */

  /* ── Ink (text & strokes) ───────────────────────────────── */
  --color-ink: #1c1917; /* warm near-black — primary text          */
  --color-ink-muted: #57534e; /* stone-600 — secondary text, labels      */
  --color-ink-subtle: #a8a29e; /* stone-400 — placeholders, disabled      */

  /* ── Translucency ───────────────────────────────────────── */
  --glass-sidebar-bg: rgba(247, 246, 243, 0.99);
  --glass-sidebar-blur: blur(48px) saturate(160%);
  --glass-content-bg: rgba(253, 252, 250, 0.99);
  --glass-content-blur: blur(24px);

  /* ── Borders ─────────────────────────────────────────────── */
  --color-border: rgba(28, 25, 23, 0.08);
  --color-border-strong: rgba(28, 25, 23, 0.16);
}
```

### Translucency — Two Tiers

The app window background is transparent (Wails `BackgroundColour: A: 0`) so the desktop wallpaper bleeds through both layers.

| Layer       | Intent                                       | CSS                                                                                |
| ----------- | -------------------------------------------- | ---------------------------------------------------------------------------------- |
| **Sidebar** | Mica-like — wallpaper reflects strongly      | `background: var(--glass-sidebar-bg); backdrop-filter: var(--glass-sidebar-blur);` |
| **Content** | Readable glass — document text must be crisp | `background: var(--glass-content-bg); backdrop-filter: var(--glass-content-blur);` |

- Never use plain `background: white` or `background: var(--color-surface)` on any visible panel — always layer through the glass tokens.
- Text on both layers uses `--color-ink` (never pure `#000000`).

### Typography

**Font: [Plus Jakarta Sans](https://fonts.google.com/specimen/Plus+Jakarta+Sans)**

- Variable font — single file, full weight range (200–800)
- Geometric with humanist warmth — distinctive but not distracting
- Excellent legibility at all sizes; well-suited to a writing app

```css
:root {
  --font-sans: "Plus Jakarta Sans", system-ui, sans-serif;
  --font-weight-normal: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;
}

body {
  font-family: var(--font-sans);
  font-weight: var(--font-weight-normal);
  color: var(--color-ink);
  -webkit-font-smoothing: antialiased;
}
```

Do **not** use: Inter, Open Sans, Roboto, Nunito, Lato, or any other generic system default.

### Icons — Phosphor Duotone

Use **[Phosphor Icons](https://phosphoricons.com/)** in `duotone` weight.

- Duotone icons carry two tones — a base and a lighter accent — giving them natural color and depth without being garish.
- Icons should **not** be monochrome grey. Assign a meaningful accent color per icon category (e.g. document icons warm amber, settings icons cool slate).
- Icon size follows a consistent scale: `16px` (inline/label), `20px` (UI controls), `24px` (navigation/sidebar).
- Never use plain stroke-only icon sets (Lucide, Heroicons) without color treatment.

### Svelte Best Practices

- Use `$state` only for variables that are genuinely reactive (drive the template or a `$derived`/`$effect`). Use plain `let` for everything else.
- Prefer `$derived` over `$effect` for computed values — `$effect` is for side effects only.
- Use `$props()` with explicit TypeScript interfaces — no loose prop typing.
- Deeply reactive objects (`$state({...})`) have proxy overhead — use `$state.raw` for large objects that are only ever reassigned (e.g. document content blobs).
- Keep components small and focused. If a component exceeds ~150 lines, consider splitting.

### Active States & Accessibility

Any UI element that can be selected, active, or focused must be **visually unambiguous** — never rely on subtle opacity shifts alone.

- **Contrast requirement**: active/selected indicators must meet WCAG AA — minimum **4.5:1** for text, **3:1** for UI components against their background.
- Use a solid accent token for the active indicator — never just a background tint. A left border, filled pill, or underline with the accent color is preferred over a faint background alone.
- Define active state tokens alongside the palette:

```css
:root {
  --color-accent: #c2410c; /* warm amber-red — primary active/selected */
  --color-accent-subtle: rgba(194, 65, 12, 0.1); /* tinted background fill   */
  --color-focus-ring: rgba(194, 65, 12, 0.5); /* keyboard focus outline    */
}
```

- All interactive elements must have a visible `:focus-visible` ring — use `--color-focus-ring`.
- Active nav items: solid `--color-accent` left border (3px) + `--color-accent-subtle` background + `--color-ink` text.
- Selected tabs/segments: filled pill in `--color-accent` with white text — never rely on weight change alone.

### Animations & Micro-interactions

All motion follows a single consistent timing system. **Never use arbitrary durations or easing** — always pick from the scale below.

```css
:root {
  /* ── Timing scale ────────────────────────────────────────── */
  --duration-instant: 80ms; /* state toggles, checkboxes             */
  --duration-fast: 150ms; /* hover, focus, button press            */
  --duration-base: 250ms; /* panel reveals, dropdown open          */
  --duration-slow: 400ms; /* page transitions, modal enter         */
  --duration-lazy: 600ms; /* large layout shifts, onboarding       */

  /* ── Easing ──────────────────────────────────────────────── */
  --ease-out: cubic-bezier(0, 0, 0.2, 1); /* decelerate — entering elements  */
  --ease-in: cubic-bezier(0.4, 0, 1, 1); /* accelerate — leaving elements   */
  --ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1); /* slight overshoot — feel alive */
  --ease-smooth: cubic-bezier(0.4, 0, 0.2, 1); /* standard — state changes        */
}
```

**Rules:**

- Hover/focus transitions: `--duration-fast` + `--ease-smooth`.
- Elements entering the screen: `--duration-base` + `--ease-out`.
- Elements leaving the screen: `--duration-fast` + `--ease-in` (exits are always faster than entrances).
- Springy/playful interactions (icon bounce, button press): `--ease-spring`.
- Never animate `width`, `height`, or `top/left` directly — use `transform` and `opacity` for performance.
- `prefers-reduced-motion`: wrap all non-essential animations in the media query — keep state changes (active, focus) but remove decorative motion.

**Micro-interaction catalogue** — apply these consistently everywhere:

| Interaction          | Animation                                                                        |
| -------------------- | -------------------------------------------------------------------------------- |
| Button press         | `transform: scale(0.97)` on `:active`, `--duration-instant`, `--ease-spring`     |
| Nav item hover       | Background fade, icon nudge `translateX(2px)`, `--duration-fast`                 |
| Nav item active      | Left border slides in, background fills, `--duration-base` + `--ease-out`        |
| Sidebar item click   | Brief scale pulse on icon `scale(1.15 → 1)`, `--duration-fast` + `--ease-spring` |
| Panel/modal enter    | `opacity: 0 → 1` + `translateY(6px) → 0`, `--duration-base` + `--ease-out`       |
| Panel/modal exit     | `opacity: 1 → 0` + `translateY(4px)`, `--duration-fast` + `--ease-in`            |
| Toast / notification | Slide in from bottom `translateY(16px) → 0`, `--duration-base` + `--ease-spring` |
| Input focus          | Border colour transition + subtle shadow grow, `--duration-fast`                 |
| Checkbox / toggle    | Fill sweep or thumb slide, `--duration-base` + `--ease-spring`                   |
| Icon in active state | Duotone layer brightens, `--duration-fast`                                       |

---

## Skills

> Install with: `npx skills add <source> --skill <name>`

### Go — `samber/cc-skills-golang` ✅ All audits pass

| Skill                   | Purpose                                           |
| ----------------------- | ------------------------------------------------- |
| `golang-code-style`     | Style rules, line length, clarity over cleverness |
| `golang-error-handling` | Idiomatic Go error wrapping and propagation       |
| `golang-database`       | Database access patterns for Go                   |
| `golang-testing`        | Table-driven tests, testify, TDD patterns, coverage |

```sh
npx skills add samber/cc-skills-golang --skill golang-code-style
npx skills add samber/cc-skills-golang --skill golang-error-handling
npx skills add samber/cc-skills-golang --skill golang-database
npx skills add samber/cc-skills-golang --skill golang-testing
```

### Svelte / SvelteKit — `sveltejs/ai-tools` ✅ All audits pass (Official Svelte team)

| Skill                       | Purpose                                                 |
| --------------------------- | ------------------------------------------------------- |
| `svelte-code-writer`        | Writing correct Svelte 5 components with MCP doc access |
| `svelte-core-bestpractices` | `$state`, `$derived`, `$effect` runes best practices    |

```sh
npx skills add sveltejs/ai-tools --skill svelte-code-writer
npx skills add sveltejs/ai-tools --skill svelte-core-bestpractices
```

### Svelte 5 — `ejirocodes/agent-skills` ✅ All audits pass

| Skill                    | Purpose                                                     |
| ------------------------ | ----------------------------------------------------------- |
| `svelte5-best-practices` | Runes, snippets, TypeScript props, SvelteKit load functions |

```sh
npx skills add ejirocodes/agent-skills --skill svelte5-best-practices
```

### TypeScript — `wshobson/agents` ✅ All audits pass

| Skill                       | Purpose                                                  |
| --------------------------- | -------------------------------------------------------- |
| `typescript-advanced-types` | Generics, conditional types, mapped types, strict typing |

```sh
npx skills add wshobson/agents --skill typescript-advanced-types
```

### Git — `github/awesome-copilot` ✅ All audits pass

| Skill        | Purpose                                                      |
| ------------ | ------------------------------------------------------------ |
| `git-commit` | Conventional Commits: feat, fix, chore, docs, refactor, etc. |

```sh
npx skills add github/awesome-copilot --skill git-commit
```
