package battle5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Attack(t *testing.T) {
	h0 := NewHero(10, 10, 39, 10, 10, false)
	assert.NotNil(t, h0)

	h1 := NewHero(10, 10, 10, 10, 39, false)
	assert.NotNil(t, h1)

	ret := SimBattle(h0, h1)
	assert.Equal(t, ret, -1)

	t.Logf("Test_Attack OK")
}

func Test_Attack1(t *testing.T) {
	h0 := NewHero(10, 11, 39, 10, 10, false)
	assert.NotNil(t, h0)

	h1 := NewHero(10, 11, 10, 10, 39, false)
	assert.NotNil(t, h1)

	ret := SimBattle(h0, h1)
	assert.Equal(t, ret, 1)

	t.Logf("Test_Attack1 OK")
}
