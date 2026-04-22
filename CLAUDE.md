# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`digital.vasic.entities` is a standalone Go module providing the media entity system:
- Title parsing for movies, TV shows, music, games, and software
- Generic entity models (MediaItem, MediaFile, MediaType, hierarchy, duplicates)
- Zero external dependencies (standard library only + testify for tests)

**Module**: `digital.vasic.entities` (Go 1.24.0)

## Commands

```bash
go test ./... -count=1          # All tests
go test -v ./... -count=1       # Verbose
go test -bench=. ./...          # Benchmarks
go build ./...                  # Compile check
go vet ./...                    # Vet
```

## Architecture

| Package | Purpose |
|---------|---------|
| `pkg/parser` | Title parsing: ParseMovieTitle, ParseTVShow, ParseMusicAlbum, ParseGameTitle, ParseSoftwareTitle, CleanTitle, ExtractYear, DetectMediaCategory |
| `pkg/models` | Data structures: MediaItem, MediaFile, MediaType, DuplicateGroup, HierarchyNode, QualityInfo, CastCrew |

## Design Patterns

- **Strategy** — each ParseXxx function applies a category-specific parsing strategy
- **Template Method** — CleanTitle and ExtractYear form the base template for all parsers
- **Value Object** — ParsedTitle, QualityInfo, CastCrew are immutable data holders
- **Repository** (interface, caller responsibility) — models are storage-agnostic

## Conventions

- Zero external dependencies in production code (stdlib only)
- Table-driven tests with `testify`
- All exported identifiers have doc comments
- Error handling: no errors from pure parsing functions; invalid input returns zero values

## Integration with Catalogizer

Catalogizer imports this module as a submodule at `Entities/`.
- `internal/services/title_parser.go` delegates to `pkg/parser`
- `internal/media/models/media.go` aliases types from `pkg/models`


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


