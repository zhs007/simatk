package battle5

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhs007/goutils"
)

func Test_LoadHeroData(t *testing.T) {
	goutils.InitLogger("", "", "debug", true, "")

	mgr, err := LoadHeroData("../gamedata/battle5/heros.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, mgr)

	assert.Equal(t, len(mgr.mapHeros), 6)

	t.Logf("Test_LoadHeroData OK")
}
