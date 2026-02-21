// Package models defines the core data structures for the media entity system.
//
// These types represent media items, files, types, and collections in a
// storage-agnostic way. Consumers (like Catalogizer) layer database-specific
// concerns on top.
//
// Design patterns applied:
//   - Value Object: ParsedTitle, QualityInfo, CastCrew are immutable data holders
//   - Repository interface (see Repository type): callers inject storage
package models

import "time"

// MediaCategory is a string-based category identifier for broad media classification.
// Examples: "movie", "tv_show", "music", "book", "game", "software".
type MediaCategory string

const (
	CategoryMovie    MediaCategory = "movie"
	CategoryTVShow   MediaCategory = "tv_show"
	CategoryTVSeason MediaCategory = "tv_season"
	CategoryEpisode  MediaCategory = "tv_episode"
	CategoryArtist   MediaCategory = "music_artist"
	CategoryAlbum    MediaCategory = "music_album"
	CategorySong     MediaCategory = "song"
	CategoryGame     MediaCategory = "game"
	CategorySoftware MediaCategory = "software"
	CategoryBook     MediaCategory = "book"
	CategoryComic    MediaCategory = "comic"
	CategoryUnknown  MediaCategory = "unknown"
)

// QualityInfo describes the detected quality attributes of a media file.
type QualityInfo struct {
	Resolution  string // "2160p", "1080p", "720p", "480p", ""
	VideoCodec  string // "H.264", "H.265", "VP9", ""
	AudioCodec  string // "AAC", "DTS", "AC3", ""
	Source      string // "BluRay", "WEB-DL", "WEBRip", "HDTV", ""
	HDR         bool
	BitDepth    int // 8, 10, 12; 0 = unknown
	FrameRate   float64
	FileSize    int64
	Bitrate     int64 // bits per second; 0 = unknown
}

// Actor represents a performer with a character name.
type Actor struct {
	Name      string
	Character string
	Order     int
}

// CastCrew groups director, writers, actors, and other credits.
type CastCrew struct {
	Director   string
	Writers    []string
	Actors     []Actor
	Producers  []string
	Musicians  []string
	Developers []string
}

// MediaType is the seeded enumeration of supported media categories
// (e.g. from a database media_types table).
type MediaType struct {
	ID          int64
	Name        string
	Description string
}

// MediaItem represents a detected media entity with aggregated metadata.
// It maps 1-to-1 with a row in a media_items-style table and forms the root
// of the hierarchy tree (parent_id → children).
type MediaItem struct {
	ID            int64
	MediaTypeID   int64
	Title         string
	OriginalTitle *string
	Year          *int
	Description   *string
	Genre         []string
	Director      *string
	CastCrew      *CastCrew
	Rating        *float64
	Runtime       *int    // minutes
	Language      *string
	Country       *string
	Status        string
	ParentID      *int64 // nil for root items; non-nil for seasons, episodes, songs
	SeasonNumber  *int
	EpisodeNumber *int
	TrackNumber   *int
	FirstDetected time.Time
	LastUpdated   time.Time
}

// MediaFile links a physical file to a MediaItem.
type MediaFile struct {
	ID          int64
	MediaItemID int64
	FileID      int64
	FilePath    string
	FileSize    int64
	FileHash    *string
	Quality     *QualityInfo
	IsPrimary   bool
	AddedAt     time.Time
}

// DuplicateGroup groups media items that share the same title, type, and year.
type DuplicateGroup struct {
	Title       string
	MediaTypeID int64
	Year        *int
	Items       []MediaItem
}

// HierarchyNode wraps a MediaItem with its children for tree traversal.
type HierarchyNode struct {
	Item     MediaItem
	Children []HierarchyNode
}
