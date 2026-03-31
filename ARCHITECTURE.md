# Architecture -- Entities

## Purpose

Generic media entity system for Go providing title parsing, entity models, and hierarchy support for media cataloging applications. Extracts structured metadata (title, year, season, episode, artist, album, platform, version) from filenames and directory names using regex-based parsers.

## Structure

```
pkg/
  parser/   Title parsing: ParseMovieTitle, ParseTVShow, ParseMusicAlbum, ParseGameTitle, ParseSoftwareTitle, CleanTitle, ExtractYear, DetectMediaCategory
  models/   Data structures: MediaItem, MediaFile, MediaType, DuplicateGroup, HierarchyNode, QualityInfo, CastCrew
```

## Key Components

- **`parser.ParseMovieTitle`** -- Extracts title and year from movie filenames
- **`parser.ParseTVShow`** -- Extracts title, season, and episode from TV show patterns (S01E02, 1x01)
- **`parser.ParseMusicAlbum`** -- Extracts artist, album, and year from music paths
- **`parser.ParseGameTitle`** -- Extracts title and platform from game filenames
- **`parser.ParseSoftwareTitle`** -- Extracts title and version from software filenames
- **`parser.DetectMediaCategory`** -- Infers media category from file extension and name patterns
- **`parser.CleanTitle`** / **`parser.ExtractYear`** -- Base utilities used by all parsers
- **`models.MediaItem`** -- Core entity with parent_id self-reference for hierarchy (TV Show -> Season -> Episode)
- **`models.MediaFile`** -- Physical file linked to a MediaItem
- **`models.HierarchyNode`** -- Tree node for entity hierarchy traversal

## Data Flow

```
filename -> DetectMediaCategory() -> category-specific ParseXxx()
    |
    ParseMovieTitle("The.Dark.Knight.2008.1080p.BluRay")
        -> CleanTitle("The.Dark.Knight") -> "The Dark Knight"
        -> ExtractYear("2008") -> 2008
        -> ParsedTitle{Title, Year}
    |
    models.MediaItem{Title, MediaTypeID, ParentID, SeasonNumber, EpisodeNumber}
```

## Dependencies

- `github.com/stretchr/testify` -- Test assertions (only dependency; zero production dependencies)

## Testing Strategy

Table-driven tests with `testify`. Tests cover all media types with comprehensive filename patterns including edge cases (missing years, unusual separators, multi-episode formats). All parsing functions are pure -- no errors returned, invalid input produces zero values.
