package battle3

func findUnit(lst []*Unit, hp int, dps int) int {
	for i, v := range lst {
		if v.Props[PropTypeHP] == hp && v.Props[PropTypeDPS] == dps {
			return i
		}
	}

	return -1
}

func genUnits(totalval int, minhp int, mindps int, minspeed int) []*Unit {
	lst := []*Unit{}

	for hp := minhp; hp <= totalval-mindps-minspeed; hp++ {
		for dps := mindps; dps <= totalval-minhp-minspeed; dps++ {
			ci := findUnit(lst, hp, dps)
			if ci < 0 {
				cu := NewUnit(hp, dps)
				lst = append(lst, cu)
			}
		}
	}

	return lst
}

// 给定初始属性，模拟全战斗
func Sim(totalval int, minhp int, mindps int, minspeed int, title string, fn string) *Stats {
	stats := &Stats{
		Title: title,
	}

	units := genUnits(totalval, minhp, mindps, minspeed)
	for i := 0; i < len(units); i++ {
		node := NewStatsNode(units[i].ResetAndClone())

		for j := 0; j < len(units); j++ {
			arr := []*Unit{units[i], units[j]}

			ret0 := StartBattle(arr)
			node.AddResult(ret0)

			// if arr[0].Props[PropTypeCurSpeed] == arr[1].Props[PropTypeCurSpeed] {
			// ret1 := StartBattle(arr, 1)
			// node.AddResult(ret1)
			// stats.Results = append(stats.Results, ret0)
			// }

			// if i != j {
			// 	ret1 := StartBattle(arr, 1)
			// 	stats.Results = append(stats.Results, ret1)
			// }
		}

		stats.Nodes = append(stats.Nodes, node)
	}

	stats.Output(fn)

	return stats
}
