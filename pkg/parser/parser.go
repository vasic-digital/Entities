// Package parser provides media title parsing for Catalogizer and other projects.
//
// It extracts structured metadata (title, year, season, episode, artist, album,
// platform, version, quality hints) from directory names and file names using
// pattern matching — with no external dependencies beyond the standard library.
//
// Supported media categories: movie, TV show, music album, game, software.
//
// Design patterns applied:
//   - Strategy: each ParseXxx function applies a category-specific parsing strategy
//   - Template Method: shared helpers (CleanTitle, ExtractYear) form the base template
package parser

import (
	"regexp"
	"strconv"
	"strings"
)

// ParsedTitle holds structured metadata extracted from a directory or file name.
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

var (
	// Year patterns
	yearParenRe   = regexp.MustCompile(`\((\d{4})\)`)
	yearBracketRe = regexp.MustCompile(`\[(\d{4})\]`)
	yearInlineRe  = regexp.MustCompile(`(?:^|[\s._-])(\d{4})(?:[\s._-]|$)`)

	// Movie: "The Matrix (1999)", "The.Matrix.1999.1080p.BluRay", "The Matrix [1999]"
	movieYearParenRe   = regexp.MustCompile(`^(.+?)[\s._-]*\((\d{4})\)`)
	movieYearBracketRe = regexp.MustCompile(`^(.+?)[\s._-]*\[(\d{4})\]`)
	movieYearDotRe     = regexp.MustCompile(`^(.+?)[\s._-]+(\d{4})[\s._-]`)

	// TV Show patterns
	tvSxxExxRe   = regexp.MustCompile(`(?i)^(.+?)[\s._-]+S(\d{1,2})E(\d{1,2})`)
	tvNxNNRe     = regexp.MustCompile(`(?i)^(.+?)[\s._-]+(\d{1,2})x(\d{2,3})`)
	tvSeasonRe   = regexp.MustCompile(`(?i)^(.+?)[\s._-]+(?:Season|S)[\s._-]*(\d{1,2})(?:[\s._-]+Episode[\s._-]*(\d{1,3}))?`)
	tvCompleteRe = regexp.MustCompile(`(?i)^(.+?)[\s._-]+(?:Complete|COMPLETE)`)

	// Music: "Artist - Album (Year)", path separator "Artist/Album"
	musicDashRe  = regexp.MustCompile(`^(.+?)\s*-\s*(.+?)(?:\s*\((\d{4})\)\s*)?$`)
	musicSlashRe = regexp.MustCompile(`^([^/]+)/([^/]+)$`)

	// Quality indicators
	qualityIndicators = []struct {
		label string
		re    *regexp.Regexp
	}{
		{"2160p", regexp.MustCompile(`(?i)2160p`)},
		{"4K", regexp.MustCompile(`(?i)\b4K\b`)},
		{"1080p", regexp.MustCompile(`(?i)1080p`)},
		{"720p", regexp.MustCompile(`(?i)720p`)},
		{"480p", regexp.MustCompile(`(?i)480p`)},
		{"BluRay", regexp.MustCompile(`(?i)(?:Blu[\s._-]?Ray|BDRip|BRRip)`)},
		{"WEB-DL", regexp.MustCompile(`(?i)(?:WEB[\s._-]*DL|WEBRip)`)},
		{"HDRip", regexp.MustCompile(`(?i)HDRip`)},
		{"DVDRip", regexp.MustCompile(`(?i)(?:DVDRip|DVD[\s._-]?Rip)`)},
		{"REMUX", regexp.MustCompile(`(?i)REMUX`)},
		{"HDR", regexp.MustCompile(`(?i)\bHDR(?:10)?\b`)},
		{"DTS", regexp.MustCompile(`(?i)\bDTS\b`)},
		{"Atmos", regexp.MustCompile(`(?i)\bAtmos\b`)},
	}

	// Game platform indicators
	gamePlatformRe = regexp.MustCompile(`(?i)\b(?:PC|Windows|Linux|Mac|macOS|PS[2-5]|PlayStation[\s._-]*[2-5]?|Xbox(?:[\s._-]*(?:One|360|Series[\s._-]*[XS]))?|Switch|Nintendo|GOG|Steam)\b`)

	// Software version: "v3.0.20", "3.0.20", "24.04"
	softwareVersionRe = regexp.MustCompile(`(?:^|[\s._-])v?(\d+(?:\.\d+)+)(?:[\s._-]|$)`)

	// Track number: leading "01 - ", "01.", "Track 01"
	trackNumberRe = regexp.MustCompile(`(?i)(?:^|[\s._-])(?:Track[\s._-]*)?(0?\d{1,2})(?:[\s._-]+|$)`)

	// Dots/underscores replacement
	dotUnderscoreRe = regexp.MustCompile(`[._]+`)

	// Collapse multiple spaces
	multiSpaceRe = regexp.MustCompile(`\s{2,}`)
)

// ParseMovieTitle extracts title and year from movie directory/file names.
// Handles patterns like "The Matrix (1999)", "The.Matrix.1999.1080p.BluRay",
// "The Matrix [1999]".
func ParseMovieTitle(dirname string) ParsedTitle {
	var result ParsedTitle

	if m := movieYearParenRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
		if y, err := strconv.Atoi(m[2]); err == nil && y >= 1900 && y <= 2099 {
			result.Year = &y
		}
	} else if m := movieYearBracketRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
		if y, err := strconv.Atoi(m[2]); err == nil && y >= 1900 && y <= 2099 {
			result.Year = &y
		}
	} else if m := movieYearDotRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
		if y, err := strconv.Atoi(m[2]); err == nil && y >= 1900 && y <= 2099 {
			result.Year = &y
		}
	} else {
		result.Title = CleanTitle(dirname)
		result.Year = ExtractYear(dirname)
	}

	result.QualityHints = extractQualityHints(dirname)
	return result
}

// ParseTVShow extracts show name, season, and episode from TV show directory/file names.
// Handles patterns like "Breaking.Bad.S01E02", "Breaking Bad - Season 1",
// "S01E02 - Pilot", "01x02".
func ParseTVShow(dirname string) ParsedTitle {
	var result ParsedTitle

	if m := tvSxxExxRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
		if s, err := strconv.Atoi(m[2]); err == nil {
			result.Season = &s
		}
		if e, err := strconv.Atoi(m[3]); err == nil {
			result.Episode = &e
		}
	} else if m := tvNxNNRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
		if s, err := strconv.Atoi(m[2]); err == nil {
			result.Season = &s
		}
		if e, err := strconv.Atoi(m[3]); err == nil {
			result.Episode = &e
		}
	} else if m := tvSeasonRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
		if s, err := strconv.Atoi(m[2]); err == nil {
			result.Season = &s
		}
		if len(m) > 3 && m[3] != "" {
			if e, err := strconv.Atoi(m[3]); err == nil {
				result.Episode = &e
			}
		}
	} else if m := tvCompleteRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
	} else {
		result.Title = CleanTitle(dirname)
	}

	result.QualityHints = extractQualityHints(dirname)
	return result
}

// ParseMusicAlbum extracts artist, album, and year from music directory names.
// Handles patterns like "Pink Floyd - The Wall (1979)", "Pink Floyd/The Wall".
func ParseMusicAlbum(dirname string) ParsedTitle {
	var result ParsedTitle

	if m := musicDashRe.FindStringSubmatch(dirname); m != nil {
		result.Artist = strings.TrimSpace(m[1])
		album := strings.TrimSpace(m[2])
		album = yearParenRe.ReplaceAllString(album, "")
		result.Album = strings.TrimSpace(album)
		result.Title = result.Album
		if m[3] != "" {
			if y, err := strconv.Atoi(m[3]); err == nil && y >= 1900 && y <= 2099 {
				result.Year = &y
			}
		}
	} else if m := musicSlashRe.FindStringSubmatch(dirname); m != nil {
		result.Artist = strings.TrimSpace(m[1])
		result.Album = strings.TrimSpace(m[2])
		result.Title = result.Album
		result.Year = ExtractYear(dirname)
	} else {
		result.Title = CleanTitle(dirname)
		result.Year = ExtractYear(dirname)
	}

	return result
}

// ParseGameTitle extracts title and platform from game directory names.
// Handles patterns like "Half-Life 2 (PC)", "The Legend of Zelda [Switch]".
func ParseGameTitle(dirname string) ParsedTitle {
	var result ParsedTitle

	if m := gamePlatformRe.FindString(dirname); m != "" {
		result.Platform = m
	}

	cleaned := dirname
	platformParenRe := regexp.MustCompile(`(?i)\s*[\(\[](?:PC|Windows|Linux|Mac|macOS|PS[2-5]|PlayStation[\s._-]*[2-5]?|Xbox(?:[\s._-]*(?:One|360|Series[\s._-]*[XS]))?|Switch|Nintendo|GOG|Steam)[\)\]]\s*`)
	cleaned = platformParenRe.ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	if m := movieYearParenRe.FindStringSubmatch(cleaned); m != nil {
		result.Title = CleanTitle(m[1])
		if y, err := strconv.Atoi(m[2]); err == nil && y >= 1900 && y <= 2099 {
			result.Year = &y
		}
	} else if m := movieYearBracketRe.FindStringSubmatch(cleaned); m != nil {
		result.Title = CleanTitle(m[1])
		if y, err := strconv.Atoi(m[2]); err == nil && y >= 1900 && y <= 2099 {
			result.Year = &y
		}
	} else if m := movieYearDotRe.FindStringSubmatch(cleaned); m != nil {
		result.Title = CleanTitle(m[1])
		if y, err := strconv.Atoi(m[2]); err == nil && y >= 1900 && y <= 2099 {
			result.Year = &y
		}
	} else {
		result.Title = CleanTitle(cleaned)
		result.Year = ExtractYear(cleaned)
	}

	return result
}

// ParseSoftwareTitle extracts name and version from software directory names.
// Handles patterns like "Ubuntu 24.04", "Microsoft Office 2021", "VLC 3.0.20".
func ParseSoftwareTitle(dirname string) ParsedTitle {
	var result ParsedTitle

	if m := softwareVersionRe.FindStringSubmatch(dirname); m != nil {
		result.Version = m[1]
	}

	if m := gamePlatformRe.FindString(dirname); m != "" {
		result.Platform = m
	}

	if m := movieYearParenRe.FindStringSubmatch(dirname); m != nil {
		result.Title = CleanTitle(m[1])
		if y, err := strconv.Atoi(m[2]); err == nil && y >= 1900 && y <= 2099 {
			result.Year = &y
		}
	} else {
		cleaned := dirname
		if result.Version != "" {
			versionEscaped := regexp.QuoteMeta(result.Version)
			vRe := regexp.MustCompile(`(?:^|[\s._-])v?` + versionEscaped + `(?:[\s._-]|$)`)
			cleaned = vRe.ReplaceAllString(cleaned, " ")
		}
		result.Title = CleanTitle(cleaned)
		result.Year = ExtractYear(dirname)
	}

	return result
}

// CleanTitle replaces dots and underscores with spaces, trims whitespace,
// and collapses multiple consecutive spaces into one.
func CleanTitle(raw string) string {
	s := dotUnderscoreRe.ReplaceAllString(raw, " ")
	s = strings.TrimSpace(s)
	for _, qi := range qualityIndicators {
		s = qi.re.ReplaceAllString(s, "")
	}
	s = strings.TrimRight(s, " -._[](){}|")
	s = multiSpaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// ExtractYear finds a 4-digit year (1900–2099) in the given string.
// Returns nil if no valid year is found.
func ExtractYear(s string) *int {
	if m := yearParenRe.FindStringSubmatch(s); m != nil {
		if y, err := strconv.Atoi(m[1]); err == nil && y >= 1900 && y <= 2099 {
			return &y
		}
	}
	if m := yearBracketRe.FindStringSubmatch(s); m != nil {
		if y, err := strconv.Atoi(m[1]); err == nil && y >= 1900 && y <= 2099 {
			return &y
		}
	}
	if m := yearInlineRe.FindStringSubmatch(s); m != nil {
		if y, err := strconv.Atoi(m[1]); err == nil && y >= 1900 && y <= 2099 {
			return &y
		}
	}
	return nil
}

// ExtractQualityHints finds quality indicator strings in a name.
func ExtractQualityHints(name string) []string {
	var hints []string
	seen := make(map[string]bool)
	for _, qi := range qualityIndicators {
		if qi.re.MatchString(name) && !seen[qi.label] {
			hints = append(hints, qi.label)
			seen[qi.label] = true
		}
	}
	return hints
}

// DetectMediaCategory infers a media category string from a file/directory path
// based on extension and TV show patterns.
//
// Returns one of: "movie", "tv_show", "music", "book", "game", "software", "unknown".
func DetectMediaCategory(path string) string {
	lower := strings.ToLower(path)

	// TV show patterns override extension-based detection
	if tvSxxExxRe.MatchString(path) || tvNxNNRe.MatchString(path) {
		return "tv_show"
	}

	switch {
	case strings.HasSuffix(lower, ".mp4"), strings.HasSuffix(lower, ".mkv"),
		strings.HasSuffix(lower, ".avi"), strings.HasSuffix(lower, ".mov"),
		strings.HasSuffix(lower, ".wmv"), strings.HasSuffix(lower, ".flv"):
		return "movie"
	case strings.HasSuffix(lower, ".mp3"), strings.HasSuffix(lower, ".flac"),
		strings.HasSuffix(lower, ".wav"), strings.HasSuffix(lower, ".aac"),
		strings.HasSuffix(lower, ".ogg"), strings.HasSuffix(lower, ".m4a"):
		return "music"
	case strings.HasSuffix(lower, ".pdf"), strings.HasSuffix(lower, ".epub"),
		strings.HasSuffix(lower, ".mobi"), strings.HasSuffix(lower, ".cbr"),
		strings.HasSuffix(lower, ".cbz"):
		return "book"
	case strings.HasSuffix(lower, ".exe"), strings.HasSuffix(lower, ".msi"),
		strings.HasSuffix(lower, ".dmg"), strings.HasSuffix(lower, ".deb"),
		strings.HasSuffix(lower, ".rpm"), strings.HasSuffix(lower, ".appimage"):
		return "software"
	}

	return "unknown"
}

// extractQualityHints is the unexported alias kept for backward compatibility.
func extractQualityHints(name string) []string {
	return ExtractQualityHints(name)
}

// unused but compiled to ensure trackNumberRe is referenced and not optimized away
var _ = trackNumberRe
