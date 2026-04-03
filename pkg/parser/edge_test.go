package parser_test

import (
	"strings"
	"testing"

	"digital.vasic.entities/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func edgeIntPtr(v int) *int { return &v }

func TestParseMovieTitle_UnicodeCharacters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		dirname   string
		wantTitle string
		wantYear  *int
	}{
		{
			"chinese title with year",
			"\u5361\u8428\u5e03\u5170\u5361 (1942)",
			"\u5361\u8428\u5e03\u5170\u5361",
			edgeIntPtr(1942),
		},
		{
			"arabic title",
			"\u0641\u064a\u0644\u0645 \u0639\u0631\u0628\u064a (2020)",
			"\u0641\u064a\u0644\u0645 \u0639\u0631\u0628\u064a",
			edgeIntPtr(2020),
		},
		{
			"emoji in title",
			"\U0001f600 Movie \U0001f3ac (2023)",
			"\U0001f600 Movie \U0001f3ac",
			edgeIntPtr(2023),
		},
		{
			"japanese title",
			"\u5343\u3068\u5343\u5c0b\u306e\u795e\u96a0\u3057 (2001)",
			"\u5343\u3068\u5343\u5c0b\u306e\u795e\u96a0\u3057",
			edgeIntPtr(2001),
		},
		{
			"korean title",
			"\uae30\uc0dd\ucda9 (2019)",
			"\uae30\uc0dd\ucda9",
			edgeIntPtr(2019),
		},
		{
			"cyrillic title",
			"\u0411\u0440\u0430\u0442 (1997)",
			"\u0411\u0440\u0430\u0442",
			edgeIntPtr(1997),
		},
		{
			"mixed scripts",
			"The \u5343\u3068\u5343\u5c0b Movie (2001)",
			"The \u5343\u3068\u5343\u5c0b Movie",
			edgeIntPtr(2001),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ParseMovieTitle(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			if tt.wantYear == nil {
				assert.Nil(t, got.Year)
			} else {
				assert.NotNil(t, got.Year)
				assert.Equal(t, *tt.wantYear, *got.Year)
			}
		})
	}
}

func TestParseMovieTitle_EmptyString(t *testing.T) {
	t.Parallel()
	got := parser.ParseMovieTitle("")
	assert.Empty(t, got.Title)
	assert.Nil(t, got.Year)
}

func TestParseMovieTitle_ExtremelyLongTitle(t *testing.T) {
	t.Parallel()
	longTitle := strings.Repeat("A", 2000) + " (2020)"
	got := parser.ParseMovieTitle(longTitle)
	assert.NotEmpty(t, got.Title)
	assert.NotNil(t, got.Year)
	assert.Equal(t, 2020, *got.Year)
}

func TestParseMovieTitle_SpecialCharacters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		dirname string
	}{
		{"forward slash", "Movie/Title (2020)"},
		{"backslash", "Movie\\Title (2020)"},
		{"quotes", `"Movie" 'Title' (2020)`},
		{"ampersand", "Tom & Jerry (2021)"},
		{"colon", "Star Wars: A New Hope (1977)"},
		{"semicolon", "Movie; Part 2 (2020)"},
		{"equals", "Movie = Greatness (2020)"},
		{"at sign", "Movie @ Home (2020)"},
		{"hash", "Movie #1 (2020)"},
		{"percent", "100% Pure (2020)"},
		{"exclamation", "Wow! (2020)"},
		{"question mark", "What? (2020)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ParseMovieTitle(tt.dirname)
			// Should not panic; title should be non-empty
			assert.NotEmpty(t, got.Title)
		})
	}
}

func TestExtractYear_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		wantYear *int
	}{
		{"no year at all", "Just A Movie", nil},
		{"year in parentheses", "Movie (2020)", edgeIntPtr(2020)},
		{"year in brackets", "Movie [2020]", edgeIntPtr(2020)},
		{"multiple years takes first paren", "Movie (2020) (2021)", edgeIntPtr(2020)},
		{"year at start", "2020 Movie", edgeIntPtr(2020)},
		{"year at end", "Movie 2020", edgeIntPtr(2020)},
		{"year below 1900", "Movie (1899)", nil},
		{"year above 2099", "Movie (2100)", nil},
		{"year exactly 1900", "Movie (1900)", edgeIntPtr(1900)},
		{"year exactly 2099", "Movie (2099)", edgeIntPtr(2099)},
		{"four digits below 1900", "Movie 1234 something", nil},
		{"numbers that look like year", "123456789", nil},
		{"empty string", "", nil},
		{"only year", "(2020)", edgeIntPtr(2020)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ExtractYear(tt.input)
			if tt.wantYear == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, *tt.wantYear, *got)
			}
		})
	}
}

func TestCleanTitle_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"only dots", "...", ""},
		{"only underscores", "___", ""},
		{"only spaces", "   ", ""},
		{"dots become spaces", "The.Dark.Knight", "The Dark Knight"},
		{"underscores become spaces", "The_Dark_Knight", "The Dark Knight"},
		{"mixed separators", "The.Dark_Knight Rises", "The Dark Knight Rises"},
		{"trailing separators stripped", "Movie---", "Movie"},
		{"quality hints removed", "Movie 1080p BluRay", "Movie"},
		{"multiple spaces collapsed", "Movie    Title", "Movie Title"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.CleanTitle(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseTVShow_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		dirname     string
		wantTitle   string
		wantSeason  *int
		wantEpisode *int
	}{
		{"uppercase SxxExx", "SHOW.S01E01", "SHOW", edgeIntPtr(1), edgeIntPtr(1)},
		{"lowercase sxxexx", "show.s01e01", "show", edgeIntPtr(1), edgeIntPtr(1)},
		{"high season number", "Show S99E99", "Show", edgeIntPtr(99), edgeIntPtr(99)},
		{"season only", "Show Season 5", "Show", edgeIntPtr(5), nil},
		{"complete series", "Show Complete", "Show", nil, nil},
		{"empty string", "", "", nil, nil},
		{"just season episode", "Show-S01E01", "Show", edgeIntPtr(1), edgeIntPtr(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ParseTVShow(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			if tt.wantSeason == nil {
				assert.Nil(t, got.Season)
			} else {
				assert.NotNil(t, got.Season)
				assert.Equal(t, *tt.wantSeason, *got.Season)
			}
			if tt.wantEpisode == nil {
				assert.Nil(t, got.Episode)
			} else {
				assert.NotNil(t, got.Episode)
				assert.Equal(t, *tt.wantEpisode, *got.Episode)
			}
		})
	}
}

func TestParseMusicAlbum_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		dirname    string
		wantTitle  string
		wantArtist string
		wantYear   *int
	}{
		{"empty string", "", "", "", nil},
		{"no dash separator", "Pink Floyd The Wall", "Pink Floyd The Wall", "", nil},
		{"dash with year", "Pink Floyd - The Wall (1979)", "The Wall", "Pink Floyd", edgeIntPtr(1979)},
		{"slash separator", "Pink Floyd/The Wall", "The Wall", "Pink Floyd", nil},
		{"unicode artist", "\u4e2d\u6587\u827a\u4eba - \u4e13\u8f91 (2020)", "\u4e13\u8f91", "\u4e2d\u6587\u827a\u4eba", edgeIntPtr(2020)},
		{"multiple dashes in album", "Artist - Album - Part 2", "Album - Part 2", "Artist", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ParseMusicAlbum(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			assert.Equal(t, tt.wantArtist, got.Artist)
			if tt.wantYear == nil {
				assert.Nil(t, got.Year)
			} else {
				assert.NotNil(t, got.Year)
				assert.Equal(t, *tt.wantYear, *got.Year)
			}
		})
	}
}

func TestParseGameTitle_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		dirname      string
		wantTitle    string
		wantPlatform string
	}{
		{"empty string", "", "", ""},
		{"no platform", "Half-Life 2", "Half-Life 2", ""},
		{"PC platform", "Half-Life 2 (PC)", "Half-Life 2", "PC"},
		{"PlayStation", "Game (PS5)", "Game", "PS5"},
		{"Switch", "Zelda [Switch]", "Zelda", "Switch"},
		{"unicode game title", "\u30bc\u30eb\u30c0 (Switch)", "\u30bc\u30eb\u30c0", "Switch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ParseGameTitle(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			assert.Equal(t, tt.wantPlatform, got.Platform)
		})
	}
}

func TestParseSoftwareTitle_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		dirname     string
		wantTitle   string
		wantVersion string
	}{
		{"empty string", "", "", ""},
		{"version with v prefix", "VLC v3.0.20", "VLC", "3.0.20"},
		{"version without v prefix", "Ubuntu 24.04", "Ubuntu", "24.04"},
		{"no version", "Microsoft Office", "Microsoft Office", ""},
		{"complex version", "App 1.2.3.4.5", "App", "1.2.3.4.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ParseSoftwareTitle(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			assert.Equal(t, tt.wantVersion, got.Version)
		})
	}
}

func TestDetectMediaCategory_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty path", "", "unknown"},
		{"no extension", "just_a_name", "unknown"},
		{"uppercase extension", "Movie.MKV", "movie"},
		{"tv pattern overrides extension", "Show.S01E02.mp4", "tv_show"},
		{"music extension", "song.mp3", "music"},
		{"book extension", "book.epub", "book"},
		{"software extension", "app.exe", "software"},
		{"double extension", "file.tar.gz", "unknown"},
		{"path with dots", "/path/to/my.movie/file.mkv", "movie"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.DetectMediaCategory(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractQualityHints_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantLen int
	}{
		{"empty string", "", 0},
		{"no hints", "Just A Movie 2020", 0},
		{"multiple hints", "Movie 1080p BluRay HDR DTS", 4},
		{"case insensitive", "movie 1080P bluray hdr", 3},
		{"duplicate hints", "1080p 1080p 1080p", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parser.ExtractQualityHints(tt.input)
			assert.Len(t, got, tt.wantLen)
		})
	}
}
