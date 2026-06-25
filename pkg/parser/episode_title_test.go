package parser_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"digital.vasic.entities/pkg/parser"
)

// redMode reports whether the RED-polarity baseline assertion is active.
//
// RED_MODE=1 reproduces the DEFECT-I Layer A defect on the pre-fix logic: it
// asserts EpisodeTitle is EMPTY for inputs that carry a human episode title
// (the old behavior, where the SxxExx regex discarded the trailing title and
// ParsedTitle had no EpisodeTitle field). On the fixed code this assertion
// FAILS — which is exactly the proof the test catches the defect (§11.4.115).
//
// RED_MODE=0 (default, GREEN) asserts the fixed behavior: the human episode
// title after the SxxExx / NxNN token is captured into EpisodeTitle.
func redMode() bool { return os.Getenv("RED_MODE") == "1" }

// episodeTitleCase is one representative TV input that carries a human episode
// title following the season/episode token.
type episodeTitleCase struct {
	name             string
	dirname          string
	wantEpisodeTitle string // expected EpisodeTitle on the FIXED code (RED_MODE=0)
	wantTitle        string
	wantSeason       int
	wantEpisode      int
}

// titledEpisodeCases all carry a real human episode title — the cases the
// defect broke. They drive both the RED baseline and the GREEN assertion.
var titledEpisodeCases = []episodeTitleCase{
	{"SxxExx dot Pilot with ext", "Breaking.Bad.S01E01.Pilot.mkv", "Pilot", "Breaking Bad", 1, 1},
	{"SxxExx dash Pilot", "Breaking Bad S01E01 - Pilot", "Pilot", "Breaking Bad", 1, 1},
	{"NxNN space Pilot", "Breaking Bad 1x01 Pilot", "Pilot", "Breaking Bad", 1, 1},
	{"SxxExx multiword title", "Breaking.Bad.S01E02.Cat's.in.the.Bag.mkv", "Cat's in the Bag", "Breaking Bad", 1, 2},
	{"SxxExx title before quality", "Show S01E04 Pilot 1080p WEB-DL", "Pilot", "Show", 1, 4},
}

// TestEpisodeTitle_RedPolarity is the §11.4.115 reproduce-first guard for
// DEFECT-I Layer A (generic episode titles discarded by the parser).
func TestEpisodeTitle_RedPolarity(t *testing.T) {
	for _, tc := range titledEpisodeCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parser.ParseTVShow(tc.dirname)

			// Show title + season/episode must be correct in BOTH polarities
			// (the fix must not disturb the previously-working fields).
			assert.Equal(t, tc.wantTitle, got.Title, "show title")
			if assert.NotNil(t, got.Season, "season parsed") {
				assert.Equal(t, tc.wantSeason, *got.Season, "season number")
			}
			if assert.NotNil(t, got.Episode, "episode parsed") {
				assert.Equal(t, tc.wantEpisode, *got.Episode, "episode number")
			}

			if redMode() {
				// RED baseline: reproduce the defect — pre-fix the trailing
				// human title was discarded, so EpisodeTitle was empty. This
				// FAILS on the fixed code, proving the guard is real.
				assert.Equal(t, "", got.EpisodeTitle,
					"RED_MODE=1 baseline: defect leaves EpisodeTitle empty (expected to FAIL on fixed code)")
			} else {
				// GREEN: the fix captures the human episode title.
				assert.Equal(t, tc.wantEpisodeTitle, got.EpisodeTitle,
					"fixed code captures the human episode title")
			}
		})
	}
}

// TestEpisodeTitle_Negative covers inputs that must leave EpisodeTitle empty so
// a caller can fall back to a generic "Episode N" label.
func TestEpisodeTitle_Negative(t *testing.T) {
	tests := []struct {
		name    string
		dirname string
	}{
		{"SxxExx no trailing title", "Breaking.Bad.S01E02"},
		{"NxNN no trailing title", "Breaking Bad 1x02"},
		{"SxxExx only quality tags", "Show.S01E03.1080p.BluRay.mkv"},
		{"SxxExx only file extension", "Show.S01E03.mkv"},
		{"Season-only no episode token", "Game of Thrones - Season 3"},
		{"Complete series", "Sopranos Complete"},
		{"plain name no episode", "Fargo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.ParseTVShow(tt.dirname)
			assert.Equal(t, "", got.EpisodeTitle,
				"no human episode title present -> EpisodeTitle must be empty")
		})
	}
}

// TestEpisodeTitle_MoviesUnaffected proves the new field is a no-op for the
// movie parser (movies have no episode title; the field stays empty).
func TestEpisodeTitle_MoviesUnaffected(t *testing.T) {
	tests := []string{
		"The Matrix (1999)",
		"The.Dark.Knight.2008.1080p.BluRay",
		"Inception [2010]",
	}
	for _, in := range tests {
		t.Run(in, func(t *testing.T) {
			got := parser.ParseMovieTitle(in)
			assert.Equal(t, "", got.EpisodeTitle, "movies carry no episode title")
		})
	}
}
