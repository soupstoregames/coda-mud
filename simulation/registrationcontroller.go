package simulation

import (
	"github.com/soupstore/coda/simulation/model"
)

// RegistrationController is an interface on Simulation that can be used to create
// new character for new accounts
type RegistrationController interface {
	MakeCharacter(name string) model.CharacterID
}

// MakeCharacter creates a new character at the next available ID
// It returns the new character's ID
func (s *Simulation) MakeCharacter(name string) model.CharacterID {
	// create new character and add to sim
	character := model.NewCharacter(name, s.spawnRoom)
	s.characters[character.ID] = character

	// add character to room
	s.spawnRoom.AddCharacter(character)

	return character.ID
}
