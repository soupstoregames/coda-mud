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

func (r *Rig) FindItem(alias string) *Item {
	if r.Backpack != nil && r.Backpack.KnownAs(alias) {
		return r.Backpack
	}

	return nil
}

// FindItemInContents searches the inventory for an item with the matching alias.
func (r *Rig) FindItemInContents(alias string) *Item {
	if r.Backpack != nil {
		for _, item := range r.Backpack.Container.Items() {
			if item.KnownAs(alias) {
				return item
			}
		}
	}

	return nil
}

func (r *Rig) RemoveItemFromInventory(item *Item) bool {
	r.Backpack.Container.RemoveItem(item.ID)
	return true
}

// Equip attempts to place the item on the rig in the item's designated rig slot.
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

func (r *Rig) Unequip(Item *Item) bool {
	if r.Backpack == Item {
		r.Backpack = nil
		return true
	}

	return false
}
