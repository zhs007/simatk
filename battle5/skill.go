package battle5

type Skill struct {
	ID   SkillID
	Data *SkillData
}

func (skill *Skill) canUseSkill() bool {
	return true
}

func (skill *Skill) findTarget(src *Hero) *HeroList {
	lst := src.findTargetWithFuncData(skill.Data.Find)

	// src.targetSkills = lst

	return lst
}

func (skill *Skill) attack(parent *BattleLogNode, src *Hero, target *Hero) {
	// 伤害
	if skill.Data.Atk != nil {
		MgrStatic.MgrFunc.Run(skill.Data.Atk, NewLibFuncParams(src.battle, src, NewHeroListEx2(target), skill, parent))
	}
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
