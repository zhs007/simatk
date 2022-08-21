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

	params, err := battle3.LoadGenMapParams("./gamedata/mt/genmap001.yaml")
	if err != nil {
		goutils.Error("LoadGenMapParams",
			zap.Error(err))

		return
	}

	mapdata, err := battle3.GenMap(params)
	if err != nil {
		goutils.Error("GenMap",
			zap.Error(err))

		return
	}

	mapdata.Save("./map001.xlsx")
}
