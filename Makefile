# ── ink-mark build system ─────────────────────────────────────────────────────
#
# Usage:
#   make release          Bump patch (0.0.x), build, overwrite build/bin/ink-mark.exe
#   make release-minor    Bump minor (0.x.0), build, overwrite
#   make release-major    Bump major (x.0.0), build, overwrite
#   make build-ab NAME=b  Build without bumping; saves to build/bin/ab/ only
#
# All release targets:
#   • Bump VERSION file
#   • Cross-compile Windows amd64 .exe via wails build
#   • Overwrite build/bin/ink-mark.exe  (cp -f handles Windows bind-mount lock)
#   • Archive a versioned copy in build/bin/releases/
#
# build-ab never touches the version or the main ink-mark.exe.
# ─────────────────────────────────────────────────────────────────────────────

VERSION_FILE := VERSION
CURRENT_VERSION := $(shell cat $(VERSION_FILE))

# Split semver into parts.
MAJOR := $(word 1, $(subst ., ,$(CURRENT_VERSION)))
MINOR := $(word 2, $(subst ., ,$(CURRENT_VERSION)))
PATCH := $(word 3, $(subst ., ,$(CURRENT_VERSION)))

BIN_DIR      := build/bin
RELEASE_DIR  := $(BIN_DIR)/releases
AB_DIR       := $(BIN_DIR)/ab
# Wails always outputs to build/bin/<name>.exe — we use a stable temp name.
WAILS_OUT    := $(BIN_DIR)/ink-mark-build.exe
MAIN_EXE     := $(BIN_DIR)/ink-mark.exe

WAILS_BUILD  := wails build -platform windows/amd64

.PHONY: release release-patch release-minor release-major build-ab _build _overwrite _dirs

# ── Default: bump patch ───────────────────────────────────────────────────────
release: release-patch

# ── Bump patch (0.0.X) ───────────────────────────────────────────────────────
release-patch: _dirs
	@NEW_PATCH=$$(($(PATCH) + 1)); \
	NEW_VERSION="$(MAJOR).$(MINOR).$$NEW_PATCH"; \
	echo "  Bumping version: $(CURRENT_VERSION) → $$NEW_VERSION"; \
	echo "$$NEW_VERSION" > $(VERSION_FILE); \
	$(MAKE) _build _overwrite NEW_VERSION=$$NEW_VERSION

# ── Bump minor (0.X.0) ───────────────────────────────────────────────────────
release-minor: _dirs
	@NEW_MINOR=$$(($(MINOR) + 1)); \
	NEW_VERSION="$(MAJOR).$$NEW_MINOR.0"; \
	echo "  Bumping version: $(CURRENT_VERSION) → $$NEW_VERSION"; \
	echo "$$NEW_VERSION" > $(VERSION_FILE); \
	$(MAKE) _build _overwrite NEW_VERSION=$$NEW_VERSION

# ── Bump major (X.0.0) ───────────────────────────────────────────────────────
release-major: _dirs
	@NEW_MAJOR=$$(($(MAJOR) + 1)); \
	NEW_VERSION="$$NEW_MAJOR.0.0"; \
	echo "  Bumping version: $(CURRENT_VERSION) → $$NEW_VERSION"; \
	echo "$$NEW_VERSION" > $(VERSION_FILE); \
	$(MAKE) _build _overwrite NEW_VERSION=$$NEW_VERSION

# ── A/B build — no version bump, no overwrite of main exe ────────────────────
build-ab: _dirs
ifndef NAME
	$(error NAME is required. Usage: make build-ab NAME=b)
endif
	@echo "  Building A/B: v$(CURRENT_VERSION)-$(NAME)"
	@$(WAILS_BUILD) \
		-ldflags "-X main.Version=$(CURRENT_VERSION)-$(NAME)" \
		-o ink-mark-ab-$(NAME).exe
	@cp -f $(BIN_DIR)/ink-mark-ab-$(NAME).exe $(AB_DIR)/ink-mark-v$(CURRENT_VERSION)-$(NAME).exe
	@echo "  ✓ Saved: $(AB_DIR)/ink-mark-v$(CURRENT_VERSION)-$(NAME).exe"
	@echo "  (main ink-mark.exe was NOT overwritten)"

# ── Internal: build to stable temp filename via wails ─────────────────────────
_build:
	@echo "  Building v$(NEW_VERSION)..."
	@$(WAILS_BUILD) \
		-ldflags "-X main.Version=$(NEW_VERSION)" \
		-o ink-mark-build.exe

# ── Internal: overwrite main exe + archive release ───────────────────────────
_overwrite:
	@cp -f $(WAILS_OUT) $(MAIN_EXE)
	@cp -f $(WAILS_OUT) $(RELEASE_DIR)/ink-mark-v$(NEW_VERSION).exe
	@echo "  ✓ $(MAIN_EXE) updated to v$(NEW_VERSION)"
	@echo "  ✓ Archived: $(RELEASE_DIR)/ink-mark-v$(NEW_VERSION).exe"

# ── Internal: ensure output directories exist ─────────────────────────────────
_dirs:
	@mkdir -p $(RELEASE_DIR) $(AB_DIR)

