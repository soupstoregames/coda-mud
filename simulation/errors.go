package simulation

import "errors"

var (
	// ErrCharacterNotFound means that an attempt was made to act on a character that does not exist
	ErrCharacterNotFound = errors.New("character not found")
	// ErrCharacterAwake is thrown when trying to wake a woke character
	ErrCharacterAwake = errors.New("character is awake")
	// ErrCharacterAsleep is thrown when trying do anything other than wake up a sleeping character
	ErrCharacterAsleep = errors.New("character is asleep")
	// ErrRoomNotFound means that an attempt was made to act on a room that does not exist
	ErrRoomNotFound = errors.New("room not found")
	// ErrContainerNotFound means that an attempt was made to act on a container that does not exist
	ErrContainerNotFound = errors.New("container not found")
	// ErrItemNotFound means that an attempt was made to act on an item that is not available to the character
	ErrItemNotFound = errors.New("item not found")
	// ErrCannotEquipItem means that a character attempted to equip an item that is not equipable
	ErrCannotEquipItem = errors.New("cannot equip item")
)
