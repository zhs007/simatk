package battle3

func GetHPAreaForUnitType(ut UnitType, total int) (int, int) {
	if ut == UnitTypeMoreDPS {
		return 0, total * 125 / 1000
	} else if ut == UnitTypeDPS {
		return total * 125 / 1000, total * 375 / 1000
	} else if ut == UnitTypeNormal {
		return total * 375 / 1000, total * 625 / 1000
	} else if ut == UnitTypeHP {
		return total * 625 / 1000, total * 875 / 1000
	}

	return total * 875 / 1000, total
}

func IsMonster(id int) bool {
	return id >= 2000 && id < 9000
}

func IsItem(id int) bool {
	return id >= 10000 && id < 20000
}

func IsEquipment(id int) bool {
	return id >= 20000 && id < 30000
}
