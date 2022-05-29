package battle3

import (
	"math/rand"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

func GenEvent(fn string, unit *Unit) ([]*Event, error) {
	ep, err := LoadEventPool(fn)
	if err != nil {
		goutils.Error("GenEvent:LoadEventPool",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	lst := []*Event{}

	for {
		e := ep.GenEvent(lst, unit)
		if e == nil {
			goutils.Error("GenEvent:GenEvent",
				zap.Error(ErrNoEvent))

			return nil, ErrNoEvent
		}

		lst = append(lst, e)
		e.Index = len(lst)

		e.StartHP = unit.Props[PropTypeCurHP]
		unit.ProcEvent(e.ID)
		e.EndHP = unit.Props[PropTypeCurHP]
		e.MaxHP = unit.Props[PropTypeMaxHP]

		if e.IsEnding {
			break
		}
	}

	return lst, nil
}

func GenMultiLineEvent(lst []*Event) *Event {
	// 顺序 + 随机 构建 tree 即可。
	// 有一个限制，就是怪物节点（ending外），不可以是最终的叶子结点，所以需要事先决定最大宽度

	width := CountWidth(lst)
	root := &Event{}

	lastleafnum := width

	// 需要考虑不步满的情况，这里有一定概率减少若干叶子节点数
	if width/2 > 0 {
		tw := width / 2
		for i := 0; i < tw; i++ {
			if rand.Int()%2 == 0 {
				width--
			} else {
				break
			}
		}
	}

	for _, v := range lst {
		if IsLeafNode(v) { // 但如果是叶子节点，最好放没有被封住的叶子节点上
			lstLeaf2 := root.GetValidLeafNodes()
			if len(lstLeaf2) > 0 {
				cr := rand.Int() % len(lstLeaf2)
				lstLeaf2[cr].AddChild(v)

				lastleafnum--

				continue
			}

			lstLeaf := root.GetLeafNodes()
			cr := rand.Int() % len(lstLeaf)
			lstLeaf[cr].AddChild(v)

			lastleafnum--

			continue
		}

		lstLeaf := root.GetLeafNodes()
		if len(lstLeaf) < width { // 如果当前宽度富余
			lstNodes := root.GetNotValidLeafNodes()
			cr := rand.Int() % len(lstNodes)
			lstNodes[cr].AddChild(v)

			continue
		}

		// 如果宽度满了
		if lastleafnum < width { // 如果剩余叶子数不足
			// 这里需要考虑还没有被封住的叶子节点数量，尽可能的不要重新开启已经被封住的叶子节点
			lstLeaf2 := root.GetValidLeafNodes()
			if len(lstLeaf2) > 0 {
				cr := rand.Int() % len(lstLeaf2)
				lstLeaf2[cr].AddChild(v)

				continue
			}
		}

		cr := rand.Int() % len(lstLeaf)
		lstLeaf[cr].AddChild(v)
	}

	return root
}
