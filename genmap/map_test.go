package genmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMap(t *testing.T) {
	m := NewMap(100, 100)
	assert.NotNil(t, m)

	assert.Equal(t, m.Tile[0][0], TileOutsideWall)
	assert.Equal(t, m.Tile[99][99], TileOutsideWall)
	assert.Equal(t, m.Tile[1][1], TileWall)

	t.Logf("Test_NewMap OK")
}
