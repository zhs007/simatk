package main

import (
	"math/rand"
	"time"

	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/battle3"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	goutils.InitLogger("", "", "debug", true, "")

	err := battle3.InitSystem("./gamedata/mt")
	if err != nil {
		goutils.Error("InitSystem",
			zap.Error(err))

		return
	}

	retdps, err := battle3.GenDPSLevelUp(0, 200)
	if err != nil {
		goutils.Error("GenDPSLevelUp",
			zap.Error(err))

		return
	}

	retdps.Output("./genlevelupdps3.json")

	rethp, err := battle3.GenHPLevelUp(0, 200)
	if err != nil {
		goutils.Error("GenHPLevelUp",
			zap.Error(err))

		return
	}

	rethp.Output("./genleveluphp3.json")

	retdps2, err := battle3.GenDPSLevelUp2(0, 200)
	if err != nil {
		goutils.Error("GenDPSLevelUp2",
			zap.Error(err))

		return
	}

	retdps2.Output("./genlevelupdps32.json")
}
