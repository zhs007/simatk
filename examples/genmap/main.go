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

	var md *battle3.MapData
	for i := 0; i < 100; i++ {
		mapdata0, err := battle3.GenMap(params)
		if err != nil {
			goutils.Error("GenMap",
				zap.Error(err))

			continue
		} else {
			md = mapdata0

			break
		}
	}

	if md != nil {
		md.ToXlsx("./map001.xlsx")
		md.ToYaml("./map001.yaml")
		md.ToJson("./map001.json")
	} else {
		goutils.Error("GenMap:100")
	}
}
