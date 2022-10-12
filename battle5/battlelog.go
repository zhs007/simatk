package battle5

import (
	"fmt"
	"io/ioutil"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type BattleLogNodeID int
type BattleLogNodeType int

const (
	BLNTBattleStart     BattleLogNodeType = 1  // 战斗开始，根节点，不需要额外的数据
	BLNTHeroComeIn      BattleLogNodeType = 2  // 角色进场
	BLNTTurnStart       BattleLogNodeType = 3  // 回合开始
	BLNTTurnEnd         BattleLogNodeType = 4  // 回合结束
	BLNTHeroMove        BattleLogNodeType = 5  // 角色主动移动
	BLNTChgProp         BattleLogNodeType = 6  // 改变属性
	BLNTBattleReady     BattleLogNodeType = 7  // 战斗开始后的准备阶段，也不需要额外的数据，处理角色上场及角色特性技能、被动技能释放、装备等
	BLNTFindTarget      BattleLogNodeType = 8  // 寻找目标
	BLNTUseSkill        BattleLogNodeType = 9  // 开始放技能
	BLNTReadySkill      BattleLogNodeType = 10 // 准备放技能
	BLNTFindSkillTarget BattleLogNodeType = 11 // 寻找技能目标
	BLNTSkillAttack     BattleLogNodeType = 12 // 技能伤害
	BLNTKillHero        BattleLogNodeType = 13 // 击杀hero
	BLNTBattleEnd       BattleLogNodeType = 14 // 战斗结束
)

type BattleLogNode struct {
	NodeID           BattleLogNodeID   `json:"nodeid"`
	Type             BattleLogNodeType `json:"type"`
	Turn             int               `json:"turn,omitempty"` // turn从1开始，turnindex从0开始
	SrcTeam          int               `json:"srcteam,omitempty"`
	SrcHeroID        HeroID            `json:"srcheroid,omitempty"`
	SrcRealHeroID    int               `json:"srcrealheroid,omitempty"`
	SrcPos           *Pos              `json:"srcpos,omitempty"`
	srcHero          *Hero             `json:"-"`
	TargetTeam       int               `json:"targetteam,omitempty"`
	TargetHeroID     HeroID            `json:"targetheroid,omitempty"`
	TargetRealHeroID int               `json:"targetrealheroid,omitempty"`
	TargetPos        *Pos              `json:"targetpos,omitempty"`
	targetHero       *Hero             `json:"-"`
	Children         []*BattleLogNode  `json:"children,omitempty"`
	Props            []PropType        `json:"props,omitempty"`
	OldPropVals      []int             `json:"oldpropvals,omitempty"`
	PropVals         []int             `json:"propvals,omitempty"`
	TargetSkillID    SkillID           `json:"targetskillid,omitempty"`
	FromSkillID      SkillID           `json:"fromskillid,omitempty"`
}

// func (bln *BattleLogNode) SetSrcPos(hero *Hero) {
// 	bln.SrcPos = &Pos{
// 		X: hero.X,
// 		Y: hero.Y,
// 	}
// }

func (bln *BattleLogNode) SetSrc(hero *Hero) {
	bln.SrcPos = hero.Pos.Clone()

	bln.SrcTeam = hero.TeamIndex + 1
	bln.srcHero = hero
	bln.SrcRealHeroID = hero.RealBattleHeroID
	bln.SrcHeroID = hero.ID
}

func (bln *BattleLogNode) SetTarget(hero *Hero) {
	bln.TargetPos = hero.Pos.Clone()

	bln.TargetTeam = hero.TeamIndex + 1
	bln.targetHero = hero
	bln.TargetRealHeroID = hero.RealBattleHeroID
	bln.TargetHeroID = hero.ID
}

func (bln *BattleLogNode) SetTargetSkill(skill *Skill) {
	bln.TargetSkillID = skill.Data.ID
}

func (bln *BattleLogNode) SetFromSkill(skill *Skill) {
	bln.FromSkillID = skill.Data.ID
}

func (bln *BattleLogNode) SetTargetPos(pos *Pos) {
	bln.TargetPos = pos
}

func (bln *BattleLogNode) genTABs(tab string, tabnum int) string {
	str := ""

	for i := 0; i < tabnum; i++ {
		str += tab
	}

	return str
}

func (bln *BattleLogNode) GenString(tab string, tabnum int, ontext FuncOnText) {
	str := bln.genTABs(tab, tabnum)

	switch bln.Type {
	case BLNTBattleStart:
		str += "战斗开始\n"
	case BLNTBattleReady:
		str += "准备阶段\n"
	case BLNTHeroComeIn:
		str += fmt.Sprintf("队%v %v(%v.%v) 入场，坐标 (%v, %v) \n",
			bln.SrcTeam,
			bln.srcHero.Data.Name,
			bln.SrcHeroID,
			bln.SrcRealHeroID,
			bln.TargetPos.X,
			bln.TargetPos.Y)

		for i, v := range bln.Props {
			str += fmt.Sprintf("%v%v (%v): %v \n",
				bln.genTABs(tab, tabnum+1),
				MapPropTypeStr[v],
				v,
				bln.PropVals[i])
		}
	case BLNTTurnStart:
		str += fmt.Sprintf("回合%v开始\n",
			bln.Turn)
	case BLNTTurnEnd:
		str += fmt.Sprintf("回合%v结束\n",
			bln.Turn)
	case BLNTFindTarget:
		if bln.targetHero != nil {
			str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 选择目标 队%v %v(%v.%v)坐标(%v, %v) \n",
				bln.SrcTeam,
				bln.srcHero.Data.Name,
				bln.SrcHeroID,
				bln.SrcRealHeroID,
				bln.SrcPos.X,
				bln.SrcPos.Y,
				bln.TargetTeam,
				bln.targetHero.Data.Name,
				bln.TargetHeroID,
				bln.TargetRealHeroID,
				bln.TargetPos.X,
				bln.TargetPos.Y)
		} else {
			str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 没有找到目标 \n",
				bln.SrcTeam,
				bln.srcHero.Data.Name,
				bln.SrcHeroID,
				bln.SrcRealHeroID,
				bln.SrcPos.X,
				bln.SrcPos.Y)
		}
	case BLNTHeroMove:
		str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 移动到 坐标(%v, %v) \n",
			bln.SrcTeam,
			bln.srcHero.Data.Name,
			bln.SrcHeroID,
			bln.SrcRealHeroID,
			bln.SrcPos.X,
			bln.SrcPos.Y,
			bln.TargetPos.X,
			bln.TargetPos.Y)
	case BLNTUseSkill:
		str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 释放了技能 %v(%v) \n",
			bln.SrcTeam,
			bln.srcHero.Data.Name,
			bln.SrcHeroID,
			bln.SrcRealHeroID,
			bln.SrcPos.X,
			bln.SrcPos.Y,
			MgrStatic.MgrSkillData.GetSkillData(bln.TargetSkillID).Name,
			bln.TargetSkillID)
	case BLNTFindSkillTarget:
		if bln.targetHero != nil {
			str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 选择目标 队%v %v(%v.%v)坐标(%v, %v) \n",
				bln.SrcTeam,
				bln.srcHero.Data.Name,
				bln.SrcHeroID,
				bln.SrcRealHeroID,
				bln.SrcPos.X,
				bln.SrcPos.Y,
				bln.TargetTeam,
				bln.targetHero.Data.Name,
				bln.TargetHeroID,
				bln.TargetRealHeroID,
				bln.TargetPos.X,
				bln.TargetPos.Y)
		} else {
			str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 没有找到目标 \n",
				bln.SrcTeam,
				bln.srcHero.Data.Name,
				bln.SrcHeroID,
				bln.SrcRealHeroID,
				bln.SrcPos.X,
				bln.SrcPos.Y)
		}
	case BLNTSkillAttack:
		str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 对 队%v %v(%v.%v)坐标(%v, %v) 造成了 %v 的伤害(HP %v->%v) from 技能 %v(%v)。 \n",
			bln.SrcTeam,
			bln.srcHero.Data.Name,
			bln.SrcHeroID,
			bln.SrcRealHeroID,
			bln.SrcPos.X,
			bln.SrcPos.Y,
			bln.TargetTeam,
			bln.targetHero.Data.Name,
			bln.TargetHeroID,
			bln.TargetRealHeroID,
			bln.TargetPos.X,
			bln.TargetPos.Y,
			bln.PropVals[0]-bln.OldPropVals[0],
			bln.OldPropVals[0],
			bln.PropVals[0],
			MgrStatic.MgrSkillData.GetSkillData(bln.FromSkillID).Name,
			bln.FromSkillID)
	case BLNTKillHero:
		str += fmt.Sprintf("队%v %v(%v.%v)坐标(%v, %v) 击杀了 队%v %v(%v.%v)坐标(%v, %v) from 技能 %v(%v)。 \n",
			bln.SrcTeam,
			bln.srcHero.Data.Name,
			bln.SrcHeroID,
			bln.SrcRealHeroID,
			bln.SrcPos.X,
			bln.SrcPos.Y,
			bln.TargetTeam,
			bln.targetHero.Data.Name,
			bln.TargetHeroID,
			bln.TargetRealHeroID,
			bln.TargetPos.X,
			bln.TargetPos.Y,
			MgrStatic.MgrSkillData.GetSkillData(bln.FromSkillID).Name,
			bln.FromSkillID)
	case BLNTBattleEnd:
		str += fmt.Sprintf("战斗结束\n")
	}

	if ontext != nil {
		ontext(str)
	}

	for _, v := range bln.Children {
		v.GenString(tab, tabnum+1, ontext)
	}
}

type BattleLog struct {
	Root      *BattleLogNode  `json:"root"`
	HashCode  string          `json:"hashcode"`
	curNodeID BattleLogNodeID `json:"-"`
	tab       string          `json:"-"`
}

func (bl *BattleLog) GenNodeID() BattleLogNodeID {
	nid := bl.curNodeID

	bl.curNodeID++

	return nid
}

func (bl *BattleLog) StartBattle() *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTBattleStart,
	}

	bl.Root = node

	return node
}

func (bl *BattleLog) BattleReady(parent *BattleLogNode) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTBattleReady,
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) HeroComeIn(parent *BattleLogNode, hero *Hero) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTHeroComeIn,
		// srcHero: hero,
	}

	node.SetSrc(hero)

	for proptype, v := range hero.Props {
		node.Props = append(node.Props, proptype)
		node.PropVals = append(node.PropVals, v)
	}

	node.SetTargetPos(hero.Pos.Clone())

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) StartTurn(parent *BattleLogNode, turn int) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTTurnStart,
		Turn:   turn,
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) EndTurn(parent *BattleLogNode, turn int) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTTurnEnd,
		Turn:   turn,
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) FindTarget(parent *BattleLogNode, src *Hero, target *Hero) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTFindTarget,
		// srcHero: src,
	}

	node.SetSrc(src)

	if target != nil {
		node.SetTarget(target)
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) HeroMove(parent *BattleLogNode, src *Hero, target *Pos) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTHeroMove,
		// srcHero: src,
	}

	node.SetSrc(src)

	node.SetTargetPos(target)

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) UseSkill(parent *BattleLogNode, src *Hero, skill *Skill) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTUseSkill,
		// srcHero: src,
	}

	node.SetSrc(src)

	node.SetTargetSkill(skill)

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) FindSkillTarget(parent *BattleLogNode, src *Hero, target *Hero) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTFindSkillTarget,
		// srcHero: src,
	}

	node.SetSrc(src)

	if target != nil {
		node.SetTarget(target)
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) SkillAttack(parent *BattleLogNode, src *Hero, target *Hero, skill *Skill, starthp int, endhp int) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTSkillAttack,
		// srcHero: src,
	}

	node.SetSrc(src)
	node.SetFromSkill(skill)

	if target != nil {
		node.SetTarget(target)
	}

	node.Props = append(node.Props, PropTypeCurHP)
	node.OldPropVals = append(node.OldPropVals, starthp)
	node.PropVals = append(node.PropVals, endhp)

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) KillHero(parent *BattleLogNode, src *Hero, target *Hero, skill *Skill) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTKillHero,
		// srcHero: src,
	}

	node.SetSrc(src)
	node.SetFromSkill(skill)

	if target != nil {
		node.SetTarget(target)
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) BattleEnd(parent *BattleLogNode) *BattleLogNode {
	node := &BattleLogNode{
		NodeID: bl.GenNodeID(),
		Type:   BLNTBattleEnd,
	}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	return node
}

func (bl *BattleLog) SaveText(fn string) error {
	if bl.Root == nil {
		return nil
	}

	f, err := os.Create(fn)
	if err != nil {
		goutils.Error("BattleLog.SaveText",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}
	defer f.Close()

	bl.Root.GenString(bl.tab, 0, func(str string) {
		_, err := f.WriteString(str)
		if err != nil {
			goutils.Error("BattleLog.SaveText:WriteString",
				zap.String("fn", fn),
				zap.Error(err))

			return
		}
	})

	return nil
}

func (bl *BattleLog) SaveJson(fn string) error {
	if bl.Root == nil {
		return nil
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	buf, err := json.Marshal(bl)
	if err != nil {
		goutils.Error("BattleLog.SaveJson:Marshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	err = ioutil.WriteFile(fn, buf, 0666)
	if err != nil {
		goutils.Error("BattleLog.SaveJson:WriteFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return nil
}

func NewBattleLog() *BattleLog {
	return &BattleLog{
		curNodeID: 1,
		tab:       "  ",
	}
}
