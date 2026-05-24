# Agent Instructions

## Build Commands

```bash
# Install dependencies
go mod download && go mod verify

# Build binary to ./bin/mangathr
make build

# Cross-compile for all platforms (darwin/linux × arm64/amd64)
make cross-build VERSION=v1.0.0

# Clean build artifacts
make clean
```

There are no tests or linter configurations in this project.

## Regenerating ent ORM Code

Database schema is defined in `ent/schema/`. After modifying schema files, regenerate with:

```bash
go generate ./ent/...
```

## Architecture

### Data Flow

```
cmd/mangathr/  (Cobra CLI commands)
    └── internal/config/      (load YAML config, merge --override flags)
    └── internal/sources/     (Scraper interface → mangadex, cubari, mangaplus)
    └── internal/downloader/  (page worker pool, writer, templater)
    └── internal/database/    (ent ORM → SQLite or PostgreSQL)
    └── internal/metadata/    (Agent interface → ComicInfo)
    └── internal/hooks/       (Hook interface → discord, webhook, subcommand)
```

- `internal/manga/` holds the core data types (`Chapter`, `Metadata`, `Page`) shared across all packages.
- `internal/rester/` is a thin HTTP client with built-in retry logic used by scrapers.
- `ent/` contains entgo-generated ORM code. The only hand-edited files are `ent/schema/manga.go` and `ent/schema/chapter.go`.

### Adding a New Source

1. Create a new package under `internal/sources/<name>/` implementing all methods of the `Scraper` interface (defined in `internal/sources/scraperController.go`).
2. Register it in the `scrapers` and `scraperTitles` maps at the top of `scraperController.go`.
3. Add a `Config` struct and `SetConfig`/`Default` functions following the package-level config pattern (see below).
4. Wire the config into `internal/config/config.go`: add a field to `Config.Sources`, call `SetConfig` in `Propagate()`, and call `Default()` in `useDefaults()`.

### Adding a New Metadata Agent

1. Create a new file in `internal/metadata/` implementing the `Agent` interface (defined in `metadataController.go`).
2. Register it in the `NewAgent` factory map in `metadataController.go`.
3. Add the new agent name to the `validateMetadataAgent` allowlist in `internal/config/config.go`.

### Adding a New Hook Type

1. Create a new file in `internal/hooks/` with a `<Type>HookConfig` struct and a hook struct that **embeds `baseHook`** (from `hook.go`). Only `Fire(ctx HookContext) error` and `FireAggregate(ctx AggregateHookContext) error` need custom implementation — all other `Hook` interface methods are satisfied by `baseHook`.
2. Add the new config type as a field in `HooksConfig` (in `controller.go`) and instantiate the hook in `newController()`.
3. Add validation for the new hook type in `validateHooks()` in `internal/config/config.go`.

## Key Conventions

### Package-Level Config Pattern

Every configurable package (`downloader`, `database`, `mangadex`, `cubari`, `hooks`) owns its config in a `config.go` (or `controller.go` for hooks) file with:
- A `Config` struct with YAML tags
- A package-level `var config Config`
- `SetConfig(cfg Config)` — called by `internal/config.Config.Propagate()`
- `(c *Config) Default(...)` — sets defaults, called by `internal/config.Config.useDefaults()`

Config is never passed as function arguments; packages read from their own `config` variable.

### Error Handling

Scraper methods return `*logging.ScraperError` (not `error`). Use `logging.ExitIfError(err)` or `logging.ExitIfErrorWithFunc(err, cleanup)` to handle these at call sites in the CLI commands. A `nil` return means success.

### Interface-Based Extensibility

`Scraper` (sources), `Agent` (metadata), and `Hook` (hooks) are interface types. Implementations are registered in factory maps or slices in their respective controller files (`scraperController.go`, `metadataController.go`, `controller.go`). The scraper/agent factory maps use lowercase keys; source names are matched case-insensitively. Hook implementations embed `baseHook` rather than using a factory map — they are instantiated directly in `newController()`.

### Filename Templating

The custom template engine lives in `internal/downloader/templater/`. Templates use `{var: prefix <.> suffix}` syntax. The `<.>` is replaced by the variable value; the entire `{...}` block is omitted if the variable is empty. The `num` variable supports zero-padding via `{num:3}`.

## Keeping This File Up-to-Date

Update `AGENTS.md` whenever you make any of the following changes:

- **New source**: Add the source name to the data flow diagram and document its config shape.
- **New metadata agent**: Add the agent name to the `validateMetadataAgent` allowlist description and note any agent-specific behaviour.
- **New hook type**: Add it to the data flow diagram and document its config struct fields.
- **New configurable package**: Add it to the Package-Level Config Pattern list and describe any deviations from the standard pattern.
- **Build system changes**: Update the Build Commands section if `Makefile` targets or dependency commands change.
- **ent schema changes**: Note any new entities or fields that affect the hand-edited schema files.
- **New interface or extensibility pattern**: Document the registration mechanism so future agents can follow the same pattern.

When in doubt, err on the side of updating this file — stale instructions are worse than no instructions.
