package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"digital.vasic.entities/pkg/parser"
)

func TestParseMovieTitle(t *testing.T) {
	tests := []struct {
		name      string
		dirname   string
		wantTitle string
		wantYear  *int
	}{
		{"paren year", "The Matrix (1999)", "The Matrix", intPtr(1999)},
		{"bracket year", "Inception [2010]", "Inception", intPtr(2010)},
		{"dot format", "The.Dark.Knight.2008.1080p.BluRay", "The Dark Knight", intPtr(2008)},
		{"no year", "Metropolis", "Metropolis", nil},
		{"quality hints removed", "Interstellar (2014) 1080p BluRay", "Interstellar", intPtr(2014)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestParseTVShow(t *testing.T) {
	tests := []struct {
		name        string
		dirname     string
		wantTitle   string
		wantSeason  *int
		wantEpisode *int
	}{
		{"SxxExx", "Breaking.Bad.S01E02", "Breaking Bad", intPtr(1), intPtr(2)},
		{"NxNN", "Breaking Bad 1x02", "Breaking Bad", intPtr(1), intPtr(2)},
		{"Season N", "Game of Thrones - Season 3", "Game of Thrones", intPtr(3), nil},
		{"Season N Ep N", "Lost Season 2 Episode 10", "Lost", intPtr(2), intPtr(10)},
		{"Complete", "Sopranos Complete", "Sopranos", nil, nil},
		{"plain name", "Fargo", "Fargo", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestParseMusicAlbum(t *testing.T) {
	tests := []struct {
		name       string
		dirname    string
		wantTitle  string
		wantArtist string
		wantAlbum  string
		wantYear   *int
	}{
		{
			"artist dash album year",
			"Pink Floyd - The Wall (1979)",
			"The Wall", "Pink Floyd", "The Wall", intPtr(1979),
		},
		{
			"artist dash album no year",
			"Led Zeppelin - IV",
			"IV", "Led Zeppelin", "IV", nil,
		},
		{
			"plain album",
			"Rumours",
			"Rumours", "", "", nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.ParseMusicAlbum(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			assert.Equal(t, tt.wantArtist, got.Artist)
			assert.Equal(t, tt.wantAlbum, got.Album)
			if tt.wantYear == nil {
				assert.Nil(t, got.Year)
			} else {
				assert.NotNil(t, got.Year)
				assert.Equal(t, *tt.wantYear, *got.Year)
			}
		})
	}
}

func TestParseGameTitle(t *testing.T) {
	tests := []struct {
		name         string
		dirname      string
		wantTitle    string
		wantPlatform string
	}{
		{"PC platform", "Half-Life 2 (PC)", "Half-Life 2", "PC"},
		{"Switch platform", "The Legend of Zelda [Switch]", "The Legend of Zelda", "Switch"},
		{"no platform", "Minecraft", "Minecraft", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.ParseGameTitle(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			assert.Equal(t, tt.wantPlatform, got.Platform)
		})
	}
}

func TestParseSoftwareTitle(t *testing.T) {
	tests := []struct {
		name        string
		dirname     string
		wantTitle   string
		wantVersion string
	}{
		{"ubuntu with version", "Ubuntu 24.04", "Ubuntu", "24.04"},
		{"vlc with version", "VLC 3.0.20", "VLC", "3.0.20"},
		{"no version", "Photoshop", "Photoshop", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.ParseSoftwareTitle(tt.dirname)
			assert.Equal(t, tt.wantTitle, got.Title)
			assert.Equal(t, tt.wantVersion, got.Version)
		})
	}
}

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"The.Matrix", "The Matrix"},
		{"Breaking_Bad", "Breaking Bad"},
		{"  multiple   spaces  ", "multiple spaces"},
		{"Inception  1080p  BluRay", "Inception"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, parser.CleanTitle(tt.input))
		})
	}
}

func TestExtractYear(t *testing.T) {
	tests := []struct {
		input string
		want  *int
	}{
		{"(1999)", intPtr(1999)},
		{"[2010]", intPtr(2010)},
		{"Title 2008 1080p", intPtr(2008)},
		{"no year here", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parser.ExtractYear(tt.input)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}

func TestDetectMediaCategory(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"movie.S01E02.mkv", "tv_show"},
		{"film.mp4", "movie"},
		{"song.mp3", "music"},
		{"book.pdf", "book"},
		{"setup.exe", "software"},
		{"random.xyz", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, parser.DetectMediaCategory(tt.path))
		})
	}
}

func TestExtractQualityHints(t *testing.T) {
	hints := parser.ExtractQualityHints("The Matrix (1999) 1080p BluRay DTS")
	assert.Contains(t, hints, "1080p")
	assert.Contains(t, hints, "BluRay")
	assert.Contains(t, hints, "DTS")
}

func intPtr(v int) *int { return &v }
