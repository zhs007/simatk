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

	unit, err := battle3.MgrStatic.MgrCharacter.NewUnit(1000)
	if err != nil {
		goutils.Error("NewUnit",
			zap.Error(err))

		return
	}

	event, err := battle3.GenEvent("./gamedata/mt/stage001.yaml", unit)
	if err != nil {
		goutils.Error("GenEvent",
			zap.Error(err))

		return
	}

	goutils.Info("event",
		goutils.JSON("event", event))
}
