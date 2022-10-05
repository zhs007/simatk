package battle5

type Skill struct {
	ID   SkillID
	Data *SkillData
}

func (skill *Skill) Clone() *Skill {
	return &Skill{
		ID:   skill.ID,
		Data: skill.Data,
	}
}

func NewSkill(data *SkillData) *Skill {
	skill := &Skill{
		ID:   data.ID,
		Data: data,
	}

	return skill
}
