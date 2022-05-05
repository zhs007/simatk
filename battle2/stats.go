package battle2

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type BattlesStatsNode struct {
	TotalNums int `json:"totalNums"`
	WinNums   int `json:"winNums"`
	LastHP    int `json:"lastHP"`
}

type StatsNode struct {
	MainUnit      *Unit            `json:"mainUnit"`
	Total         BattlesStatsNode `json:"total"`
	ForceFirst    BattlesStatsNode `json:"forceFirst"`
	ForceNotFirst BattlesStatsNode `json:"forceNotFirst"`
}

func NewStatsNode(unit *Unit) *StatsNode {
	return &StatsNode{
		MainUnit: unit,
	}
}

func (node *StatsNode) AddResult(ret *BattleResult) {
	node.Total.TotalNums++

	if ret.WinIndex == 0 {
		node.Total.WinNums++

		node.Total.LastHP += ret.Units[0].Props[PropTypeCurHP]
	}

	if ret.ForceFirstIndex == 0 {
		node.ForceFirst.TotalNums++

		if ret.WinIndex == 0 {
			node.ForceFirst.WinNums++

			node.ForceFirst.LastHP += ret.Units[0].Props[PropTypeCurHP]
		}
	} else if ret.ForceFirstIndex == 1 {
		node.ForceNotFirst.TotalNums++

		if ret.WinIndex == 0 {
			node.ForceNotFirst.WinNums++

			node.ForceNotFirst.LastHP += ret.Units[0].Props[PropTypeCurHP]
		}
	}
}

type Stats struct {
	Title string       `json:"title"`
	Nodes []*StatsNode `json:"nodes"`
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
