package battle3

import "math/rand"

type RandWeights struct {
	Vals       []int `yaml:"vals"`
	ValWeights []int `yaml:"weights"`
	MaxWeight  int   `yaml:"maxweight"`
}

func (rw *RandWeights) onLoad() {
	rw.MaxWeight = 0

	for i := range rw.ValWeights {
		rw.MaxWeight += rw.ValWeights[i]
	}
}

func (rw *RandWeights) GenVal() int {
	w := rand.Int() % rw.MaxWeight
	for i := range rw.ValWeights {
		if w < rw.ValWeights[i] {
			return rw.Vals[i]
		}

		w -= rw.ValWeights[i]
	}

	return -1
}
