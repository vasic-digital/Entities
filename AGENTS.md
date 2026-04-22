# AGENTS.md - Entities Module Multi-Agent Coordination

## Module Identity

- **Module**: `digital.vasic.entities`
- **Role**: Generic media entity system -- title parsing, entity models, hierarchy support
- **Packages**: `pkg/parser`, `pkg/models`
- **Go Version**: 1.24.0
- **External Dependencies**: None (stdlib only; `testify` for tests)

## Agent Responsibilities

### Entities Agent

The Entities agent owns both packages in this module:

1. **Title Parsing** (`pkg/parser`) -- Regex-based parsing functions for movies, TV shows, music albums, games, and software. Includes `CleanTitle`, `ExtractYear`, and `DetectMediaCategory` utilities that form the base template for all category-specific parsers.

2. **Entity Models** (`pkg/models`) -- Data structures for the media entity system: `MediaItem` (self-referential hierarchy via `ParentID`), `MediaFile` (physical file metadata), `MediaType` (enumeration entry), `DuplicateGroup`, `HierarchyNode`, `QualityInfo`, `CastCrew`.

## Cross-Agent Coordination

### Upstream Consumers

Any service importing `digital.vasic.entities` should coordinate with the Entities agent when:

- Adding a new media category or parser function
- Modifying `MediaItem`, `MediaFile`, or `MediaType` struct fields
- Changing hierarchy semantics (`ParentID`, `SeasonNumber`, `EpisodeNumber`)
- Altering `ParsedTitle` or any parsed output struct

### Integration Points

| Consumer | Package Used | Purpose |
|----------|-------------|---------|
| catalog-api `internal/services/title_parser.go` | `pkg/parser` | Post-scan title extraction |
| catalog-api `internal/media/models/media.go` | `pkg/models` | Entity type aliases for internal use |
| catalog-api `internal/services/aggregation_service.go` | Both | Entity creation from scanned files |

## Dependency Graph

```
pkg/models  (no internal dependencies -- leaf package)
pkg/parser  (no internal dependencies -- leaf package)
```

Both packages are fully independent with zero internal imports between them.

## Conventions

- Zero external dependencies in production code (standard library only)
- Table-driven tests with `testify`
- Test naming: `Test<Function>_<Scenario>`
- All exported identifiers have doc comments
- Pure functions: parsing never returns errors; invalid input yields zero values
- Value objects: `ParsedTitle`, `QualityInfo`, `CastCrew` are immutable data holders

## Testing Standards

```bash
go test ./... -count=1          # All tests
go test -v ./... -count=1       # Verbose
go test -bench=. ./...          # Benchmarks
go build ./...                  # Compile check
go vet ./...                    # Vet
```

## Constraints

- **No CI/CD pipelines**: GitHub Actions, GitLab CI/CD, and all automated pipeline configurations are permanently disabled. All testing is local.
- **No project-specific code**: This module must remain generic and reusable with no references to Catalogizer internals.
- **Struct field changes are breaking**: Any modification to exported model structs affects all consumers and requires coordination.


## ⚠️ MANDATORY: NO SUDO OR ROOT EXECUTION

**ALL operations MUST run at local user level ONLY.**

This is a PERMANENT and NON-NEGOTIABLE security constraint:

- **NEVER** use `sudo` in ANY command
- **NEVER** use `su` in ANY command
- **NEVER** execute operations as `root` user
- **NEVER** elevate privileges for file operations
- **ALL** infrastructure commands MUST use user-level container runtimes (rootless podman/docker)
- **ALL** file operations MUST be within user-accessible directories
- **ALL** service management MUST be done via user systemd or local process management
- **ALL** builds, tests, and deployments MUST run as the current user

### Container-Based Solutions
When a build or runtime environment requires system-level dependencies, use containers instead of elevation:

- **Use the `Containers` submodule** (`https://github.com/vasic-digital/Containers`) for containerized build and runtime environments
- **Add the `Containers` submodule as a Git dependency** and configure it for local use within the project
- **Build and run inside containers** to avoid any need for privilege escalation
- **Rootless Podman/Docker** is the preferred container runtime

### Why This Matters
- **Security**: Prevents accidental system-wide damage
- **Reproducibility**: User-level operations are portable across systems
- **Safety**: Limits blast radius of any issues
- **Best Practice**: Modern container workflows are rootless by design

### When You See SUDO
If any script or command suggests using `sudo` or `su`:
1. STOP immediately
2. Find a user-level alternative
3. Use rootless container runtimes
4. Use the `Containers` submodule for containerized builds
5. Modify commands to work within user permissions

**VIOLATION OF THIS CONSTRAINT IS STRICTLY PROHIBITED.**


### ⚠️⚠️⚠️ ABSOLUTELY MANDATORY: ZERO UNFINISHED WORK POLICY

NO unfinished work, TODOs, or known issues may remain in the codebase. EVER.

PROHIBITED: TODO/FIXME comments, empty implementations, silent errors, fake data, unwrap() calls that panic, empty catch blocks.

REQUIRED: Fix ALL issues immediately, complete implementations before committing, proper error handling in ALL code paths, real test assertions.

Quality Principle: If it is not finished, it does not ship. If it ships, it is finished.
