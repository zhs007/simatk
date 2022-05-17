package battle3

type AddOn struct {
	CacheVals []int
}

func (addon *AddOn) Clone() *AddOn {
	a := &AddOn{
		CacheVals: make([]int, len(addon.CacheVals)),
	}

	copy(a.CacheVals, addon.CacheVals)

	return a
}
