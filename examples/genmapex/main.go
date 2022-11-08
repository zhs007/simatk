package main

import (
	"math/rand"
	"time"

	"github.com/zhs007/goutils"
	"github.com/zhs007/simatk/genmap"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	goutils.InitLogger("", "", "debug", true, "")

	m, err := genmap.GenMap(30, 30, 1, 1, 5, 5)
	if err != nil {
		goutils.Error("GenMap",
			zap.Error(err))

		return
	}

	err = m.Save("./output/map001.xlsx")
	if err != nil {
		goutils.Error("LoadGenMapParams",
			zap.Error(err))

		return
	}
}
