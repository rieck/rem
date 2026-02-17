# rem - CLI for macOS Reminders

## Non-Negotiables
- **Conventional Commits**: ALL commits MUST follow [Conventional Commits](https://www.conventionalcommits.org/). Format: `type(scope): description`. Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `build`, `ci`, `perf`. No exceptions.

## What is this?
Go CLI wrapping macOS Reminders. Uses `go-eventkit` (cgo + Objective-C EventKit) for fast reads AND writes (<200ms) as a single binary — including reminder CRUD and list CRUD. AppleScript only for flagged operations and default list name query. Provides CRUD for reminders/lists, natural language dates, and import/export. For programmatic Go access, use `go-eventkit` directly.

## Architecture
- `cmd/rem/commands/` - Cobra CLI commands (one file per command); `huh_helpers.go` has shared interactive utilities
- `internal/service/` - Service layer: `reminders.go` and `lists.go` wrap `go-eventkit` client. `executor.go` runs `osascript` for flagged operations and default list name query only.
- `internal/reminder/` - Domain models: `Reminder`, `List`, `Priority`
- `internal/parser/` - Custom NL date parser (no external deps)
- `internal/export/` - JSON/CSV import/export
- `internal/skills/` - Agent skill install/uninstall/status logic
- `internal/update/` - Background update check (GitHub releases, 24h cache)
- `internal/ui/` - Table (`olekukonko/tablewriter` v1.x), plain, JSON output
- `skills/rem-cli/` - Embedded agent skill files (go:embed'd into binary)

## Critical: Architecture Rules
- **ALL reads AND writes go through `go-eventkit`** (`github.com/BRO3886/go-eventkit/reminders`) — in-process EventKit via cgo, <200ms
- **Single binary** — go-eventkit's cgo code compiled into the binary
- **AppleScript only for**: flagged operations (EventKit doesn't expose flagged), default list name query
- **EventKit doesn't expose `flagged`** - JXA fallback only used when `--flagged` filter is active, AppleScript for flag/unflag writes
- **go-eventkit field names**: `Title` (not `Name`), `Notes` (not `Body`), `List` (not `ListName`), native `URL` field
- **List CRUD via go-eventkit**: `CreateList` (auto-discovers source), `UpdateList` (ID-based), `DeleteList` (ID-based). Immutable lists are rejected.
- Priority: 0=none, 1-4=high, 5=medium, 6-9=low
- `due date` and `remind me date` are independent

## Libraries
- `BRO3886/go-eventkit` - **EventKit bindings** (cgo + ObjC, reads AND writes)
- `spf13/cobra` - CLI framework
- `olekukonko/tablewriter` v1.x - **new API**: `NewTable()`, `.Header()`, `.Append()`, `.Render()` (NOT the old `SetHeader`/`SetBorder` API)
- `fatih/color` - terminal colors
- `olekukonko/tablewriter/tw` - alignment constants (`tw.AlignLeft`)
- `charmbracelet/huh` - interactive forms (multi-select, select, input, confirm) for all `-i` modes and skills

## Build & Test
```bash
make build        # -> bin/rem (single binary, includes EventKit via cgo)
make release      # -> bin/rem-darwin-arm64.tar.gz (for GitHub Releases upload)
go test ./...     # unit tests (date parser, export, models)
make completions  # bash/zsh/fish
```

## Releasing
- **CRITICAL: Tag BEFORE building.** The Makefile uses `git describe --tags` to embed the version. If you `make release` before `git tag`, the binary will have the wrong version (e.g., `v0.5.0-2-gae75da9` instead of `v0.6.0`). Always: tag → push tag → `make release` → verify with `./rem version` → upload.
- **Always use `make release`** to build release binaries — produces a `.tar.gz` that preserves execute permissions (HTTP downloads strip +x from raw binaries)
- Upload `bin/rem-darwin-arm64.tar.gz` to GitHub Releases
- No CI release workflow — macOS runners are too expensive for free tier

## Agent Skills
- `rem skills install/uninstall/status` manages embedded agent skill files
- Skills are `go:embed`'d from `skills/rem-cli/` into the binary
- `internal/skills/` handles install/uninstall logic, `internal/update/` handles background update checks
- Background update check runs in a goroutine during `PersistentPreRun`, prints notice in `PersistentPostRun`
- `REM_NO_UPDATE_CHECK=1` disables update checks; also skipped for dev builds, json output, non-TTY, meta commands

## Conventions
- Short IDs displayed as first 8 chars of full `x-apple-reminder://UUID` ID
- Prefix matching: users can pass partial IDs to any command
- All commands support `-o json|table|plain`
- `NO_COLOR` env var respected
- **Interactive mode (`-i`)**: All mutation commands (add, complete, delete, flag, unflag, update, list-mgmt create/rename/delete) support `-i` for huh-based interactive forms. Shared helpers in `huh_helpers.go`. Theme: `ThemeCatppuccin()`. All handle `ErrUserAborted` gracefully.

## Website & Hosting
- Documentation site: `rem.sidv.dev` (Hugo on Cloudflare Pages)
- Source: `website/` dir, deployed via `.github/workflows/deploy.yml`
- Install script served at `rem.sidv.dev/install` (from `website/static/install`)
- Domain: `sidv.dev` (owned by user, managed on Cloudflare)

## Journal
Engineering journals live in `journals/` dir. See `.claude/commands/journal.md` for the journaling command.
