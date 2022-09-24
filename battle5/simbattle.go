package battle5

// 模拟战斗，就是基础的1打1; 返回1，hero0赢，-1，hero1赢，0平
func SimBattle(hero0 *Hero, hero1 *Hero) int {
	lst := []*Hero{hero0, hero1}
	first := 0
	second := 1

	if hero0.Props[PropTypeCurSpeed] < hero1.Props[PropTypeCurSpeed] {
		first = 1
		second = 0
	}

	turn := 0
	for {
		if lst[first].Attack(lst[second]) {
			return 1
		}

		if lst[second].Attack(lst[first]) {
			return -1
		}

		turn++
		if turn >= 10 {
			break
		}
	}

	return 0
}
