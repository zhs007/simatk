package battle

func findUnit(lst []*Unit, hp int, dps int) int {
	for i, v := range lst {
		if v.Props[PropTypeHP] == hp && v.Props[PropTypeDPS] == dps {
			return i
		}
	}

	return -1
}

func genUnits(totalval int, minhp int, mindps int) []*Unit {
	lst := []*Unit{}

	for hp := minhp; hp <= totalval-mindps; hp++ {
		ci := findUnit(lst, hp, totalval-hp)
		if ci < 0 {
			cu := NewUnit(hp, totalval-hp)
			lst = append(lst, cu)
		}
	}

	for dps := mindps; dps <= totalval-minhp; dps++ {
		ci := findUnit(lst, totalval-dps, dps)
		if ci < 0 {
			cu := NewUnit(totalval-dps, dps)
			lst = append(lst, cu)
		}
	}

	return lst
}

// 给定初始属性，模拟全战斗
func Sim(totalval int, minhp int, mindps int, title string, fn string) *Stats {
	stats := &Stats{
		Title: title,
	}

	units := genUnits(totalval, minhp, mindps)
	for i := 0; i < len(units)-1; i++ {
		for j := i + 1; j < len(units); j++ {
			arr := []*Unit{units[i], units[j]}

			ret0 := StartBattle(arr, 0)
			stats.Results = append(stats.Results, ret0)

			ret1 := StartBattle(arr, 1)
			stats.Results = append(stats.Results, ret1)
		}
	}

	stats.Output(fn)

	return stats
}
