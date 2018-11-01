package entities

import (
	"github.com/soupstore/coda/components"
	"github.com/soupstore/ecs"
)

type GameMap struct {
	ecs.BasicEntity
	components.Geography
}
