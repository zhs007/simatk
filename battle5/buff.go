package battle5

type Buff struct {
	InstanceID BuffInstanceID
	Data       *BuffData
	From       *Hero
	FromSkill  *Skill
	isRemoved  bool
}

func (buff *Buff) IsMe(b *Buff) bool {
	return buff.InstanceID == b.InstanceID
}
