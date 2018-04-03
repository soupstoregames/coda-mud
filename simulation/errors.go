package simulation

import "errors"

var (
	// ErrCharacterNotFound means that an attempt was made to act on a character that does not exist
	ErrCharacterNotFound = errors.New("character not found")
	// ErrCharacterAlreadyAwake is thrown when trying to sleep a woke character
	ErrCharacterAlreadyAwake = errors.New("character is already awake")
	// ErrCharacterAlreadyAsleep is thrown when trying to sleep a sleeping character
	ErrCharacterAlreadyAsleep = errors.New("character is already asleep")
	// ErrRoomNotFound means that an attempt was made to act on a room that does not exist
	ErrRoomNotFound = errors.New("room not found")
)
