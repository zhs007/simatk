package battle5

type BuffList struct {
	Buffs []*Buff
}

func NewBuffList() *BuffList {
	return &BuffList{}
}

func (lst *BuffList) Add(buff *Buff) {
	lst.Buffs = append(lst.Buffs, buff)
}

func (lst *BuffList) RemoveAll() {
	nlst := []*Buff{}

	isremoved := false
	for _, v := range lst.Buffs {
		if !v.isRemoved {
			nlst = append(nlst, v)

			isremoved = true
		}
	}

	if isremoved {
		lst.Buffs = nlst
	}
}
