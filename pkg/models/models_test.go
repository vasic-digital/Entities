package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"digital.vasic.entities/pkg/models"
)

func TestMediaCategory_Constants(t *testing.T) {
	categories := []models.MediaCategory{
		models.CategoryMovie, models.CategoryTVShow, models.CategoryTVSeason,
		models.CategoryEpisode, models.CategoryArtist, models.CategoryAlbum,
		models.CategorySong, models.CategoryGame, models.CategorySoftware,
		models.CategoryBook, models.CategoryComic, models.CategoryUnknown,
	}
	for _, c := range categories {
		assert.NotEmpty(t, string(c))
	}
}

func TestMediaItem_ZeroValue(t *testing.T) {
	var item models.MediaItem
	assert.Zero(t, item.ID)
	assert.Empty(t, item.Title)
	assert.Nil(t, item.ParentID)
	assert.Nil(t, item.Year)
	assert.True(t, item.FirstDetected.IsZero())
}

func TestMediaItem_WithParent(t *testing.T) {
	parent := models.MediaItem{
		ID:          1,
		Title:       "Breaking Bad",
		MediaTypeID: 2,
		FirstDetected: time.Now(),
		LastUpdated:   time.Now(),
	}

	seasonID := parent.ID
	episode := models.MediaItem{
		ID:            2,
		ParentID:      &seasonID,
		Title:         "Pilot",
		EpisodeNumber: intPtr(1),
		SeasonNumber:  intPtr(1),
	}

	assert.NotNil(t, episode.ParentID)
	assert.Equal(t, parent.ID, *episode.ParentID)
	assert.Equal(t, 1, *episode.SeasonNumber)
	assert.Equal(t, 1, *episode.EpisodeNumber)
}

func TestDuplicateGroup(t *testing.T) {
	year := 1999
	g := models.DuplicateGroup{
		Title:       "The Matrix",
		MediaTypeID: 1,
		Year:        &year,
		Items: []models.MediaItem{
			{ID: 1, Title: "The Matrix"},
			{ID: 2, Title: "The Matrix"},
		},
	}
	assert.Len(t, g.Items, 2)
	assert.Equal(t, 1999, *g.Year)
}

func TestHierarchyNode(t *testing.T) {
	root := models.HierarchyNode{
		Item: models.MediaItem{ID: 1, Title: "Breaking Bad"},
		Children: []models.HierarchyNode{
			{Item: models.MediaItem{ID: 2, Title: "Season 1"}},
			{Item: models.MediaItem{ID: 3, Title: "Season 2"}},
		},
	}
	assert.Len(t, root.Children, 2)
	assert.Equal(t, "Season 1", root.Children[0].Item.Title)
}

func intPtr(v int) *int { return &v }
