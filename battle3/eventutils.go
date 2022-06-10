package battle3

func CountNumWithID(lst []*Event, id int) int {
	num := 0

	for _, v := range lst {
		if v.ID == id {
			num++
		}
	}

	return num
}

func CountWidth(lst []*Event) int {
	num := 0

	for _, v := range lst {
		if IsLeafNode(v) {
			num++
		}
	}

	return num
}

func IsLeafNode(e *Event) bool {
	return e.IsEnding || IsItem(e.ID) || IsEquipment(e.ID)
}

func CountAvgLastHPPer(lst []*Event) int {
	lasthp := 0

	for _, v := range lst {
		lasthp += v.EndHP * 100 / v.MaxHP
	}

	return lasthp / len(lst)
}

func CountEventNum(lst []*Event) int {
	num := 0

	for _, v := range lst {
		num++

		num += len(v.Awards)
	}

	return num
}
