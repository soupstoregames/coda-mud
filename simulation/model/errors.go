package model

import "errors"

var (
	// ErrNotEquipable is returned when a character attemps to equip an item that has no rig slot.
	ErrNotEquipable = errors.New("item is not equipable")
)
