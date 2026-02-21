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
