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
