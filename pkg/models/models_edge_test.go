package models_test

import (
	"strings"
	"testing"
	"time"

	"digital.vasic.entities/pkg/models"
	"github.com/stretchr/testify/assert"
)

// --- MediaItem Edge Cases ---

func TestMediaItem_EmptyTitle(t *testing.T) {
	t.Parallel()
	item := models.MediaItem{
		ID:          1,
		MediaTypeID: 1,
		Title:       "",
	}
	assert.Empty(t, item.Title)
	assert.Equal(t, int64(1), item.ID)
}

func TestMediaItem_UnicodeTitle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		title string
	}{
		{"chinese", "\u5343\u4e0e\u5343\u5bfb\u7684\u795e\u96a0"},
		{"arabic", "\u0641\u064a\u0644\u0645 \u0639\u0631\u0628\u064a"},
		{"emoji", "\U0001f3ac Movie \U0001f37f"},
		{"japanese", "\u6771\u4eac\u7269\u8a9e"},
		{"korean", "\uae30\uc0dd\ucda9"},
		{"cyrillic", "\u0411\u0440\u0430\u0442"},
		{"thai", "\u0e20\u0e32\u0e1e\u0e22\u0e19\u0e15\u0e23\u0e4c\u0e44\u0e17\u0e22"},
		{"mixed_scripts", "The \u5343\u4e0e\u5343\u5bfb Movie"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			item := models.MediaItem{
				Title: tt.title,
			}
			assert.Equal(t, tt.title, item.Title)
		})
	}
}

func TestMediaItem_ExtremelyLongTitle(t *testing.T) {
	t.Parallel()
	longTitle := strings.Repeat("A", 10000)
	item := models.MediaItem{
		Title: longTitle,
	}
	assert.Len(t, item.Title, 10000)
}

func TestMediaItem_SpecialCharactersInTitle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		title string
	}{
		{"forward_slash", "Movie/Title"},
		{"backslash", "Movie\\Title"},
		{"double_quotes", `"Movie Title"`},
		{"single_quotes", "'Movie Title'"},
		{"null_bytes", "Movie\x00Title"},
		{"newlines", "Movie\nTitle"},
		{"tabs", "Movie\tTitle"},
		{"angle_brackets", "<Movie> & <Title>"},
		{"sql_injection", "'; DROP TABLE media_items; --"},
		{"xss_payload", "<script>alert('xss')</script>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			item := models.MediaItem{
				Title: tt.title,
			}
			assert.Equal(t, tt.title, item.Title)
		})
	}
}

func TestMediaItem_NilOptionalFields(t *testing.T) {
	t.Parallel()
	item := models.MediaItem{
		ID:            1,
		MediaTypeID:   1,
		Title:         "Test",
		FirstDetected: time.Now(),
		LastUpdated:   time.Now(),
	}
	assert.Nil(t, item.OriginalTitle)
	assert.Nil(t, item.Year)
	assert.Nil(t, item.Description)
	assert.Nil(t, item.Director)
	assert.Nil(t, item.CastCrew)
	assert.Nil(t, item.Rating)
	assert.Nil(t, item.Runtime)
	assert.Nil(t, item.Language)
	assert.Nil(t, item.Country)
	assert.Nil(t, item.ParentID)
	assert.Nil(t, item.SeasonNumber)
	assert.Nil(t, item.EpisodeNumber)
	assert.Nil(t, item.TrackNumber)
}

func TestMediaItem_ZeroValueYear(t *testing.T) {
	t.Parallel()
	year := 0
	item := models.MediaItem{
		Title: "Test",
		Year:  &year,
	}
	assert.NotNil(t, item.Year)
	assert.Equal(t, 0, *item.Year)
}

func TestMediaItem_NegativeYear(t *testing.T) {
	t.Parallel()
	year := -500
	item := models.MediaItem{
		Title: "Ancient Film",
		Year:  &year,
	}
	assert.Equal(t, -500, *item.Year)
}

func TestMediaItem_SelfReferentialParent(t *testing.T) {
	t.Parallel()
	id := int64(42)
	item := models.MediaItem{
		ID:       42,
		ParentID: &id,
	}
	assert.Equal(t, item.ID, *item.ParentID)
}

// --- QualityInfo Edge Cases ---

func TestQualityInfo_ZeroValues(t *testing.T) {
	t.Parallel()
	q := models.QualityInfo{}
	assert.Empty(t, q.Resolution)
	assert.Empty(t, q.VideoCodec)
	assert.Empty(t, q.AudioCodec)
	assert.Empty(t, q.Source)
	assert.False(t, q.HDR)
	assert.Zero(t, q.BitDepth)
	assert.Zero(t, q.FrameRate)
	assert.Zero(t, q.FileSize)
	assert.Zero(t, q.Bitrate)
}

func TestQualityInfo_NegativeFileSize(t *testing.T) {
	t.Parallel()
	q := models.QualityInfo{
		FileSize: -1,
	}
	assert.Equal(t, int64(-1), q.FileSize)
}

// --- MediaFile Edge Cases ---

func TestMediaFile_EmptyFilePath(t *testing.T) {
	t.Parallel()
	mf := models.MediaFile{
		ID:          1,
		MediaItemID: 1,
		FileID:      1,
		FilePath:    "",
	}
	assert.Empty(t, mf.FilePath)
}

func TestMediaFile_NilQuality(t *testing.T) {
	t.Parallel()
	mf := models.MediaFile{
		ID: 1,
	}
	assert.Nil(t, mf.Quality)
}

func TestMediaFile_NilFileHash(t *testing.T) {
	t.Parallel()
	mf := models.MediaFile{
		ID: 1,
	}
	assert.Nil(t, mf.FileHash)
}

// --- DuplicateGroup Edge Cases ---

func TestDuplicateGroup_EmptyItems(t *testing.T) {
	t.Parallel()
	dg := models.DuplicateGroup{
		Title:       "Test",
		MediaTypeID: 1,
	}
	assert.Empty(t, dg.Items)
	assert.Nil(t, dg.Year)
}

// --- HierarchyNode Edge Cases ---

func TestHierarchyNode_NoChildren(t *testing.T) {
	t.Parallel()
	node := models.HierarchyNode{
		Item: models.MediaItem{
			Title: "Leaf",
		},
	}
	assert.Empty(t, node.Children)
}

func TestHierarchyNode_DeeplyNested(t *testing.T) {
	t.Parallel()
	// Build a hierarchy 100 levels deep
	leaf := models.HierarchyNode{
		Item: models.MediaItem{Title: "leaf"},
	}
	current := leaf
	for i := 0; i < 100; i++ {
		current = models.HierarchyNode{
			Item:     models.MediaItem{Title: "level"},
			Children: []models.HierarchyNode{current},
		}
	}
	// Verify we can traverse all the way down
	node := current
	depth := 0
	for len(node.Children) > 0 {
		node = node.Children[0]
		depth++
	}
	assert.Equal(t, 100, depth)
	assert.Equal(t, "leaf", node.Item.Title)
}

// --- CastCrew Edge Cases ---

func TestCastCrew_EmptySlices(t *testing.T) {
	t.Parallel()
	cc := models.CastCrew{}
	assert.Empty(t, cc.Director)
	assert.Nil(t, cc.Writers)
	assert.Nil(t, cc.Actors)
	assert.Nil(t, cc.Producers)
	assert.Nil(t, cc.Musicians)
	assert.Nil(t, cc.Developers)
}

func TestCastCrew_UnicodeNames(t *testing.T) {
	t.Parallel()
	cc := models.CastCrew{
		Director: "\u5bae\u5d0e\u99ff",
		Actors: []models.Actor{
			{Name: "\u67cf\u539f\u5d07", Character: "\u5343\u5c0b", Order: 0},
		},
	}
	assert.Equal(t, "\u5bae\u5d0e\u99ff", cc.Director)
	assert.Len(t, cc.Actors, 1)
}

// --- MediaCategory Edge Cases ---

func TestMediaCategory_AllConstants(t *testing.T) {
	t.Parallel()
	categories := []models.MediaCategory{
		models.CategoryMovie,
		models.CategoryTVShow,
		models.CategoryTVSeason,
		models.CategoryEpisode,
		models.CategoryArtist,
		models.CategoryAlbum,
		models.CategorySong,
		models.CategoryGame,
		models.CategorySoftware,
		models.CategoryBook,
		models.CategoryComic,
		models.CategoryUnknown,
	}
	assert.Len(t, categories, 12)

	// Each should be non-empty and unique
	seen := make(map[models.MediaCategory]bool)
	for _, c := range categories {
		assert.NotEmpty(t, string(c))
		assert.False(t, seen[c], "duplicate category: %s", c)
		seen[c] = true
	}
}

// --- MediaType Edge Cases ---

func TestMediaType_ZeroID(t *testing.T) {
	t.Parallel()
	mt := models.MediaType{
		ID:   0,
		Name: "test",
	}
	assert.Zero(t, mt.ID)
}

func TestMediaType_EmptyName(t *testing.T) {
	t.Parallel()
	mt := models.MediaType{
		ID: 1,
	}
	assert.Empty(t, mt.Name)
}
