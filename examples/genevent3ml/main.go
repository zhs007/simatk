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

	// lst0, err := battle3.GenEvent("./gamedata/mt/stage001.yaml", unit.Clone())
	lst0, err := battle3.GenEventWithStage(unit, 1, 2)
	if err != nil {
		goutils.Error("GenEvent",
			zap.Error(err))

		return
	}

	goutils.Info("event",
		goutils.JSON("event", lst0))

	lst := []*battle3.Event{}

	for i := 0; i < 10; i++ {
		event0 := battle3.GenMultiLineEvent(lst0)

		num := event0.CountNodes()
		if num != battle3.CountEventNum(lst0) {
			goutils.Info("multi-line event",
				goutils.JSON("event", event0),
				goutils.JSON("lst", lst0))

			return
		}

		winnum := battle3.CalcWinTimesWithAI2(event0, 100, unit.Clone())
		event0.TotalNum = 100
		event0.WinNum = winnum

		goutils.Info("multi-line event",
			goutils.JSON("event", event0))

		lst = append(lst, event0)
	}

	battle3.SaveEvents2("genevent3ml.xlsx", lst)
}
