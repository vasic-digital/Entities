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
