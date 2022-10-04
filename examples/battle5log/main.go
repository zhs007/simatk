package main

import (
	"math/rand"
	"time"

	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle5"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	goutils.InitLogger("", "", "debug", true, "")

	mgr, err := battle5.NewStaticMgr("./gamedata/battle5")
	if err != nil {
		goutils.Error("NewStaticMgr",
			zap.Error(err))

		return
	}

	battle5.MgrStatic = mgr

	battle := battle5.NewBattleEx(mgr,
		[]battle5.HeroID{10000, 10001, 10002, 10003, 10004},
		[]battle5.HeroID{10001, 10002, 10003, 10004, 10005},
		battle5.SceneWidth, battle5.SceneHeight)

	battle.StartBattle()

	battle.Log.SaveText("./output/battlo5log.txt")
}
