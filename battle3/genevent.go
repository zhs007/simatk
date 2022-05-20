package battle3

import (
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
		e, isending := ep.GenEvent(lst, unit)
		if e == nil {
			goutils.Error("GenEvent:GenEvent",
				zap.Error(ErrNoEvent))

			return nil, ErrNoEvent
		}

		lst = append(lst, e)

		unit.ProcEvent(e.ID)

		if isending {
			break
		}
	}

	return lst, nil
}
