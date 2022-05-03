package battle

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type Stats struct {
	Info    string          `json:"email"`
	Results []*BattleResult `json:"results"`
}

func (stats *Stats) Output(fn string) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(stats)
	if err != nil {
		goutils.Warn("Stats.Output:Marshal",
			zap.Error(err))

		return err
	}

	err = os.WriteFile(fn, b, 0644)
	if err != nil {
		goutils.Warn("Stats.Output:WriteFile",
			zap.Error(err))

		return err
	}

	return nil
}
