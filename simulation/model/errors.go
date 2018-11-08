package model

import "errors"

var (
	// ErrNotEquipable is returned when a character attempts to equip an item that has no rig slot.
	ErrNotEquipable = errors.New("item is not equipable")

	// ErrItemNotInRig is returned when a character attempts to remove an item from a rig slot but no item with that alias can be found.
	ErrItemNotInRig = errors.New("item is not in rig")
)
