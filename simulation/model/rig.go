package model

// RigSlot is an enum of locations on a rig items can be equipped to.
type RigSlot byte

const (
	// RigSlotNone is used for items that cannot be equipped.
	RigSlotNone RigSlot = iota
	// RigSlotBackpack designates an item as wearable on the back.
	RigSlotBackpack
)

// Rig is a structure of various item mount points that represents where items can be equipped to.
type Rig struct {
	Backpack *Item
}

// Equip attemps to place the item on the rig in the item's designated rig slot.
// When an item is equipped, you get a reference to the item that was already there returned back. This will be nil if nothing was equipped there.
// If the item is not equippable for any reason you get an error.
func (r *Rig) Equip(item *Item) (*Item, error) {
	if item.Definition.RigSlot == RigSlotNone {
		return nil, ErrNotEquipable
	}
	oldBackpack := r.Backpack
	r.Backpack = item
	return oldBackpack, nil
}
