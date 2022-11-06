package battle5

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhs007/goutils"
)

func Test_LoadSkillData(t *testing.T) {
	goutils.InitLogger("", "", "debug", true, "")

	NewStaticMgr("../gamedata/battle5", nil)

	mgr, err := LoadSkillData("../gamedata/battle5/skills.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, mgr)

	assert.Equal(t, len(mgr.mapSkills), 2)

	t.Logf("Test_LoadSkillData OK")
}
