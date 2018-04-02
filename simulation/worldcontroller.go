package simulation

import "github.com/soupstore/coda-world/simulation/model"

type WorldController interface {
	MakeRoom(name, description string) model.RoomID
	SetSpawnRoom(id model.RoomID)
}
