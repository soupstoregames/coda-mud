package model

type Rig struct {
	Backpack *Backpack
}

func (r *Rig) EquipBackpack(backpack *Backpack) *Backpack {
	oldBackpack := r.Backpack
	r.Backpack = backpack
	return oldBackpack
}
