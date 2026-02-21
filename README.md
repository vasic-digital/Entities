# digital.vasic.entities

Generic media entity system for Go — title parsing, entity models, and hierarchy support for media cataloging applications.

## Features

- **Title Parsing** — extract title, year, season, episode, artist, album, platform, version from filenames/directory names
- **Entity Models** — `MediaItem`, `MediaFile`, `MediaType`, `DuplicateGroup`, `HierarchyNode`
- **Media Detection** — infer media category from file extensions and patterns
- **Zero Dependencies** — pure Go standard library in production code
- **Fully Tested** — table-driven tests with `testify`

## Install

```bash
go get digital.vasic.entities@latest
```

## Usage

### Title Parsing

```go
import "digital.vasic.entities/pkg/parser"

// Movies
movie := parser.ParseMovieTitle("The.Dark.Knight.2008.1080p.BluRay")
// movie.Title == "The Dark Knight"
// *movie.Year == 2008

// TV Shows
ep := parser.ParseTVShow("Breaking.Bad.S01E02")
// ep.Title == "Breaking Bad"
// *ep.Season == 1, *ep.Episode == 2

// Music
album := parser.ParseMusicAlbum("Pink Floyd - The Wall (1979)")
// album.Artist == "Pink Floyd"
// album.Album == "The Wall"
// *album.Year == 1979

// Games
game := parser.ParseGameTitle("Half-Life 2 (PC)")
// game.Title == "Half-Life 2"
// game.Platform == "PC"

// Software
sw := parser.ParseSoftwareTitle("Ubuntu 24.04")
// sw.Title == "Ubuntu"
// sw.Version == "24.04"

// Utilities
category := parser.DetectMediaCategory("episode.S01E02.mkv") // "tv_show"
year := parser.ExtractYear("Title (2023)")                    // *int(2023)
clean := parser.CleanTitle("The.Dark.Knight")                 // "The Dark Knight"
```

### Entity Models

```go
import "digital.vasic.entities/pkg/models"

item := models.MediaItem{
    Title:       "Breaking Bad",
    MediaTypeID: 1,
    Status:      "active",
}

season := models.MediaItem{
    ParentID:     &item.ID,
    Title:        "Season 1",
    SeasonNumber: intPtr(1),
}
```

## Module Name

`digital.vasic.entities`

## License

MIT
