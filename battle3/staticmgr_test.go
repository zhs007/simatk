package battle3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadAllStatic(t *testing.T) {
	ep, err := LoadAllStatic("../gamedata/mt/")
	assert.NoError(t, err)
	assert.NotNil(t, ep)

	// assert.Equal(t, len(ep.Events), 6)

	t.Logf("Test_LoadAllStatic OK")
}
