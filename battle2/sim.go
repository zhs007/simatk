package battle2

func findUnit(lst []*Unit, hp int, dps int, speed int) int {
	for i, v := range lst {
		if v.Props[PropTypeHP] == hp && v.Props[PropTypeDPS] == dps && v.Props[PropTypeSpeed] == speed {
			return i
		}
	}

	return -1
}

func genUnits(totalval int, minhp int, mindps int, minspeed int) []*Unit {
	lst := []*Unit{}

	for hp := minhp; hp <= totalval-mindps-minspeed; hp++ {
		for dps := mindps; dps <= totalval-minhp-minspeed; dps++ {
			speed := totalval - hp - dps
			if speed >= minspeed && speed <= totalval-minhp-mindps {
				ci := findUnit(lst, hp, dps, speed)
				if ci < 0 {
					cu := NewUnit(hp, dps, speed)
					lst = append(lst, cu)
				}
			}
			// for speed := minspeed; dps <= totalval-mindps-minhp; speed++ {
			// }
		}
	}

	// for dps := mindps; dps <= totalval-minhp; dps++ {
	// 	ci := findUnit(lst, totalval-dps, dps)
	// 	if ci < 0 {
	// 		cu := NewUnit(totalval-dps, dps)
	// 		lst = append(lst, cu)
	// 	}
	// }

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

			ret0 := StartBattle(arr, 0)
			node.AddResult(ret0)

			if arr[0].Props[PropTypeCurSpeed] == arr[1].Props[PropTypeCurSpeed] {
				ret1 := StartBattle(arr, 1)
				node.AddResult(ret1)
				// stats.Results = append(stats.Results, ret0)
			}

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
