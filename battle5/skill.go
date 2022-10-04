package battle5

type Skill struct {
	ID   SkillID
	Data *SkillData
}

func NewSkill(data *SkillData) *Skill {
	skill := &Skill{
		ID:   data.ID,
		Data: data,
	}

	return skill
}
