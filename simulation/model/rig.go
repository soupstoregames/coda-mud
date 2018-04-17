package model

type RigSlot byte

const (
	RigSlotNone RigSlot = iota
	RigSlotBackpack
)

type Rig struct {
	Backpack *Item
}

func (r *Rig) Equip(item *Item) (*Item, error) {
	if item.RigSlot == RigSlotNone {
		return nil, ErrNotEquipable
	}
	oldBackpack := r.Backpack
	r.Backpack = item
	return oldBackpack, nil
}
