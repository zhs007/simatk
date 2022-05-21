package battle3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadEventPool(t *testing.T) {
	ep, err := LoadEventPool("../gamedata/mt/stage001.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, ep)

	assert.Equal(t, len(ep.Events), 7)

	t.Logf("Test_LoadEventPool OK")
}
