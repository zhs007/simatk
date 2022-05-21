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

	ret, err := battle3.GenDPSLevelUp(0, 100)
	if err != nil {
		goutils.Error("GenDPSLevelUp",
			zap.Error(err))

		return
	}

	ret.Output("./genlevelupdps3.json")
}
