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

	lst := [][]*battle3.Event{}

	for i := 0; i < 10; i++ {
		unit, err := battle3.MgrStatic.MgrCharacter.NewUnit(1000)
		if err != nil {
			goutils.Error("NewUnit",
				zap.Error(err))

			return
		}

		lst0, err := battle3.GenEvent("./gamedata/mt/stage001.yaml", unit)
		if err != nil {
			goutils.Error("GenEvent",
				zap.Error(err))

			return
		}

		goutils.Info("event",
			goutils.JSON("event", lst0))

		lst = append(lst, lst0)
	}

	battle3.SaveEvents("genevent3.xlsx", lst)
}
