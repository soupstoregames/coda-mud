package entities

import (
	"github.com/soupstore/coda/components"
	"github.com/soupstore/ecs"
)

// PlayerCharacter is the entity that represents the player in the game.
// It will be controlled by a user that is connected to the game.
type PlayerCharacter struct {
	ecs.BasicEntity
	components.MapPosition
	components.CommandQueue
}
