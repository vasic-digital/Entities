# Architecture

## Purpose

`digital.vasic.entities` is a standalone Go module that defines the core data structures and title-parsing logic for the Catalogizer media entity system. It has zero production dependencies beyond the standard library, making it safe to import from any Go project.

## Package Overview

| Package | Path | Description |
|---------|------|-------------|
| `models` | `pkg/models/` | Storage-agnostic entity types: `MediaItem`, `MediaFile`, `MediaType`, `QualityInfo`, `CastCrew`, `DuplicateGroup`, `HierarchyNode`, and the `MediaCategory` enum (12 constants) |
| `parser` | `pkg/parser/` | Regex-based title parsing for 5 media categories (movie, TV, music, game, software), plus shared utilities `CleanTitle`, `ExtractYear`, `ExtractQualityHints`, and `DetectMediaCategory` |

## Design Patterns

| Pattern | Where | How |
|---------|-------|-----|
| Strategy | `parser.ParseMovieTitle`, `ParseTVShow`, `ParseMusicAlbum`, `ParseGameTitle`, `ParseSoftwareTitle` | Each function applies a category-specific set of regex patterns to extract structured metadata from a raw directory/file name |
| Template Method | `parser.CleanTitle`, `parser.ExtractYear` | Shared helpers that every `ParseXxx` function delegates to for dot/underscore normalization and year extraction |
| Value Object | `models.ParsedTitle`, `models.QualityInfo`, `models.CastCrew` | Immutable data holders with no behavior; constructed and returned by parsers |
| Self-Referential Hierarchy | `models.MediaItem.ParentID`, `models.HierarchyNode` | `ParentID *int64` enables TV Show -> Season -> Episode and Artist -> Album -> Song trees; `HierarchyNode` wraps items for recursive traversal |
| Repository (caller-provided) | `models` package doc | Models are storage-agnostic; the consuming project (Catalogizer) layers its own repository implementations on top |

## Key Interfaces and Types

### models

```go
type MediaCategory string   // "movie", "tv_show", "tv_season", "tv_episode", "music_artist", ...

type MediaItem struct {
    ID            int64
    MediaTypeID   int64
    Title         string
    Year          *int
    ParentID      *int64      // nil = root; non-nil = season, episode, song
    SeasonNumber  *int
    EpisodeNumber *int
    TrackNumber   *int
    Genre         []string
    CastCrew      *CastCrew
    // ... timestamps, ratings, language, country, status
}

type MediaFile struct {
    ID          int64
    MediaItemID int64
    FileID      int64
    FilePath    string
    Quality     *QualityInfo
    IsPrimary   bool
}

type HierarchyNode struct {
    Item     MediaItem
    Children []HierarchyNode
}

type DuplicateGroup struct {
    Title       string
    MediaTypeID int64
    Year        *int
    Items       []MediaItem
}
```

### parser

```go
type ParsedTitle struct {
    Title        string
    Year         *int
    QualityHints []string
    Season       *int
    Episode      *int
    Artist       string
    Album        string
    TrackNumber  *int
    Platform     string
    Version      string
}

func ParseMovieTitle(dirname string) ParsedTitle
func ParseTVShow(dirname string) ParsedTitle
func ParseMusicAlbum(dirname string) ParsedTitle
func ParseGameTitle(dirname string) ParsedTitle
func ParseSoftwareTitle(dirname string) ParsedTitle

func CleanTitle(raw string) string
func ExtractYear(s string) *int
func ExtractQualityHints(name string) []string
func DetectMediaCategory(path string) string   // "movie" | "tv_show" | "music" | "book" | "software" | "unknown"
```

## Usage Example

```go
import (
    "fmt"
    "digital.vasic.entities/pkg/parser"
    "digital.vasic.entities/pkg/models"
)

// Parse a directory name into structured metadata
parsed := parser.ParseTVShow("Breaking.Bad.S01E02.720p.BluRay")
fmt.Println(parsed.Title)   // "Breaking Bad"
fmt.Println(*parsed.Season) // 1
fmt.Println(*parsed.Episode) // 2

// Detect category from a file path
cat := parser.DetectMediaCategory("song.mp3") // "music"

// Build a hierarchy node
show := models.MediaItem{ID: 1, Title: "Breaking Bad", MediaTypeID: 2}
season := models.HierarchyNode{
    Item:     models.MediaItem{ID: 2, Title: "Season 1", ParentID: &show.ID},
    Children: nil,
}
root := models.HierarchyNode{Item: show, Children: []models.HierarchyNode{season}}
```
