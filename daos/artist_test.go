package daos

import (
	"testing"

	"github.com/ederavilaprado/golang-web-architecture-template/app"
	"github.com/ederavilaprado/golang-web-architecture-template/models"
	"github.com/ederavilaprado/golang-web-architecture-template/testdata"
	"github.com/stretchr/testify/assert"
)

func TestArtistDAO(t *testing.T) {
	db := testdata.ResetDB()
	dao := NewArtistDAO(db)

	{
		// Get
		testDBCall(db, func(rs app.RequestContext) {
			artist, err := dao.Get(rs, 2)
			assert.Nil(t, err)
			if assert.NotNil(t, artist) {
				assert.Equal(t, 2, artist.Id)
			}
		})
	}

	{
		// Create
		testDBCall(db, func(rs app.RequestContext) {
			artist := &models.Artist{
				Id:   1000,
				Name: "tester",
			}
			err := dao.Create(rs, artist)
			assert.Nil(t, err)
			assert.NotEqual(t, 1000, artist.Id)
			assert.NotZero(t, artist.Id)
		})
	}

	{
		// Update
		testDBCall(db, func(rs app.RequestContext) {
			artist := &models.Artist{
				Id:   2,
				Name: "tester",
			}
			err := dao.Update(rs, artist.Id, artist)
			assert.Nil(t, err)
		})
	}

	{
		// Update with error
		testDBCall(db, func(rs app.RequestContext) {
			artist := &models.Artist{
				Id:   2,
				Name: "tester",
			}
			err := dao.Update(rs, 99999, artist)
			assert.NotNil(t, err)
		})
	}

	{
		// Delete
		testDBCall(db, func(rs app.RequestContext) {
			err := dao.Delete(rs, 2)
			assert.Nil(t, err)
		})
	}

	{
		// Delete with error
		testDBCall(db, func(rs app.RequestContext) {
			err := dao.Delete(rs, 99999)
			assert.NotNil(t, err)
		})
	}

	{
		// Query
		testDBCall(db, func(rs app.RequestContext) {
			artists, err := dao.Query(rs, 1, 3)
			assert.Nil(t, err)
			assert.Equal(t, 3, len(artists))
		})
	}

	{
		// Count
		testDBCall(db, func(rs app.RequestContext) {
			count, err := dao.Count(rs)
			assert.Nil(t, err)
			assert.NotZero(t, count)
		})
	}
}
