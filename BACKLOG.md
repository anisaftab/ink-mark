# ink-mark — Feature Backlog

Ordered by implementation priority. Each item notes what's missing on the backend (Go) and frontend (Svelte) sides, plus hard dependencies.

---

## ✅ Done

- **New document** — save dialog → file on disk → DB record → editor opens
- **Save document** — writes content to disk, updates title in DB
- **Load document by ID** — reads file from disk
- **ProseKit rich-text editor** — renders, toolbar formatting, markdown serialisation
- **Shell UI** — sidebar, titlebar, tabs, doc-toolbar, status bar (visual)
- **Home page** — mica gradient card, New Document + Open Document CTAs
- **1. Open existing file** — `OpenDocument()` Go method + `GetByFilePath` repo method; "Open Document" button and "Open file" menu item wired correctly
- **2. List documents in sidebar** — `ListDocuments()` Go method; sidebar Documents section populated from DB on startup; clicking opens the doc in a tab
- **3. Active document state + tabs** — shared `openDocuments` / `activeDocumentId` / `activeDocument` stores; tab bar iterates real open docs; close (×) works; sidebar click activates or loads
- **4. Markdown preview pane** — `markdown-it` renders live HTML from the active document's markdown content; Edit / Split / Read mode switch fully functional; full prose stylesheet

---

## 🔲 Backlog

### 1. Open existing file (`Ctrl O`)
"Open Document" and "Open file" both incorrectly call `NewDocument`. The correct open-file dialog flow is missing entirely.

- **Backend:** `OpenDocument()` Wails method — open-file dialog → check DB for existing record or register new one → return DTO with file content
- **Frontend:** Wire the "Open Document" home page button, "Open file" create-dropdown item, and `Ctrl O` shortcut
- **Depends on:** —

### 2. List documents — sidebar Documents section
Sidebar shows hardcoded fake items. `repository.List()` already exists in Go.

- **Backend:** `ListDocuments()` Wails method — calls `repository.List()`, returns `[]DocumentDTO`
- **Frontend:** Load list on startup, replace hardcoded sidebar items with `{#each}`, clicking opens the document
- **Depends on:** —

### 3. Active document state + tab bar wiring
One `activeDocument` lives locally in `+page.svelte`. Tab bar is fully static HTML. No shared session model exists.

- **Backend:** —
- **Frontend:** Shared app state (open documents list + active index), tabs reflect open docs, clicking a sidebar item opens it in a tab, `+` triggers new document, active tab highlights correctly
- **Depends on:** 1, 2

### 4. Markdown preview pane
Preview pane shows a static placeholder. Mode switch works visually but preview is empty. `prosemirror-markdown` is already installed.

- **Backend:** —
- **Frontend:** Serialise ProseKit document to HTML, render live in preview pane, update on every editor change
- **Depends on:** —

### 5. Projects — backend model + repository
`Project` doesn't exist anywhere — no model, no DB table, no migration, no repository.

- **Backend:** DB migration adding `projects` table; `project` package with model + repository (Create, GetByID, List, AddDocument, ListDocuments); Wails methods: `NewProject`, `ListProjects`, `GetProjectDocuments`
- **Frontend:** —
- **Depends on:** —

### 6. Projects — frontend sidebar tree
Projects section shows hardcoded placeholder trees.

- **Backend:** —
- **Frontend:** Load projects + their documents on startup, render collapsible tree, clicking a document opens it in a tab, unsaved dot tied to dirty state
- **Depends on:** 5

### 7. Open folder as project
"Open folder as project" menu item is a stub.

- **Backend:** `OpenFolderAsProject()` — folder-picker dialog → scan for `.md` files → create project record → register all found files → return project DTO
- **Frontend:** Wire the create-dropdown menu item
- **Depends on:** 5

### 8. Auto-save + unsaved indicator
Save is entirely manual. Unsaved dot and save status in toolbar are static.

- **Backend:** —
- **Frontend:** Dirty flag on document, auto-save on 2s debounce, unsaved dot in sidebar, save status in doc-toolbar reflects real state, `Ctrl+S` keybinding
- **Depends on:** 3

### 9. Recent documents
"Recent" quick-access goes nowhere. No `last_opened_at` timestamp exists.

- **Backend:** DB migration adding `last_opened_at`; update it on `LoadDocument`; `ListRecent(limit int)` Wails method
- **Frontend:** Populate Recent section in sidebar on startup
- **Depends on:** 2

### 10. Status bar — live editor stats
Word count and cursor position are hardcoded static strings.

- **Backend:** —
- **Frontend:** Read cursor position from ProseKit editor state, count words from serialised markdown, update reactively
- **Depends on:** —

### 11. Tab bar — close + unsaved guard
Tabs are visual-only. No confirmation when closing a dirty document.

- **Backend:** —
- **Frontend:** Close tab (× button), confirm dialog if dirty, active tab title synced to document title
- **Depends on:** 3, 8

### 12. Window controls — Wails runtime
Minimize / maximize / close buttons do nothing.

- **Backend:** —
- **Frontend:** Wire to `runtime.WindowMinimise`, `runtime.WindowMaximise`, `runtime.WindowClose`
- **Depends on:** —

### 13. Search (`Ctrl K`)
Search bar is entirely unimplemented.

- **Backend:** `SearchDocuments(query string)` — SQLite `LIKE` or FTS5 full-text search across title and content
- **Frontend:** Command palette modal, result list, pressing Enter opens the document
- **Depends on:** 2

### 14. Pinned documents
"Pinned" section and pin icons are static.

- **Backend:** DB migration adding `pinned boolean`; `PinDocument` / `UnpinDocument` methods; `ListPinned()` query
- **Frontend:** Populate Pinned section, wire pin button in doc-toolbar, pin icon on sidebar items
- **Depends on:** 2

### 15. Export (PDF / HTML)
Export button is a stub. No rendering infrastructure exists.

- **Backend:** `ExportDocument(id, format string)` — markdown → HTML (Go library), optionally PDF; save-file dialog for output
- **Frontend:** Format picker in doc-toolbar, progress/error feedback
- **Depends on:** 4
