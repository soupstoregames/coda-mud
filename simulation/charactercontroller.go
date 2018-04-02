package simulation

import "github.com/soupstore/coda-world/simulation/model"

type CharacterController interface {
	MakeCharacter(name string) model.CharacterID
	WakeUpCharacter(model.CharacterID) <-chan interface{}
	SleepCharacter(model.CharacterID)
}
