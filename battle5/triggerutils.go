package battle5

import "strings"

func NewTriggerDataTurnStart(turn int) *TriggerData {
	return &TriggerData{
		Trigger: TriggerTypeTurnStart,
		Turn:    turn,
	}
}

func NewTriggerDataTurnEnd(turn int) *TriggerData {
	return &TriggerData{
		Trigger: TriggerTypeTurnEnd,
		Turn:    turn,
	}
}

func NewTriggerDataFind(turn int, src *Hero) *TriggerData {
	return &TriggerData{
		Trigger: TriggerTypeFind,
		Turn:    turn,
		Src:     src,
	}
}

func ParseTriggerList(str string) []TriggerType {
	lst := []TriggerType{}

	arr := strings.Split(str, "|")
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v != "" {
			ct := Str2TriggerType(v)
			if ct != TriggerTypeNone {
				lst = append(lst, ct)
			}
		}
	}

	return lst
}
