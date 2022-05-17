package battle3

type Equipment struct {
	Data  *ItemData
	AddOn *AddOn
}

func (equip *Equipment) Clone() *Equipment {
	return &Equipment{
		Data:  equip.Data,
		AddOn: equip.AddOn.Clone(),
	}
}
