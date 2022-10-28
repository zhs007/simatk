package battle5funcs

import "github.com/zhs007/simatk/battle5"

func isValidBool(val int) bool {
	return val == 0 || val == 1
}

func isValidTargetType(str string) bool {
	return str == "enemy" || str == "friend" || str == "all"
}

func toTargetType(str string) battle5.TargetType {
	switch str {
	case "enemy":
		return battle5.TargetTypeEnemy
	case "friend":
		return battle5.TargetTypeFriend
	case "all":
		return battle5.TargetTypeAll
	}

	return battle5.TargetTypeEnemy
}
