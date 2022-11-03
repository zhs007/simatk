package battle5buffs

// 简单状态
// 一个角色身上只有一个生效状态，且高级覆盖低级，同级间不可刷新
type SimpleState struct {
	AttachTurn int
}
