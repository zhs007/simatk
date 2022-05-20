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
