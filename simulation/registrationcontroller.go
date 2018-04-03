package simulation

import "github.com/soupstore/coda-world/simulation/model"

// RegistrationController is an interface on Simulation that can be used to create
// new character for new accounts
type RegistrationController interface {
	MakeCharacter(name string) model.CharacterID
}

// MakeCharacter creates a new character at the next available ID
// It returns the new character's ID
func (s *Simulation) MakeCharacter(name string) model.CharacterID {
	// get and increment character ID
	characterID := s.getNextCharacterID()

	// create new character and add to sim
	character := model.NewCharacter(characterID, name, s.spawnRoom)
	s.characters[characterID] = character

	// add character to room
	s.spawnRoom.AddCharacter(character)

	return characterID
}
