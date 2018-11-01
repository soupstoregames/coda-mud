package systems

import (
	"github.com/soupstore/coda/components"
	"github.com/soupstore/ecs"
)

type weatherEntity struct {
	*ecs.BasicEntity
	*components.Geography
}

type Weather struct {
	entities []weatherEntity
}

func (s *Weather) Add(basic *ecs.BasicEntity, geography *components.Geography) {
	s.entities = append(s.entities, weatherEntity{basic, geography})
}

func (s *Weather) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *Weather) Update(dt float32) {
	for i := range s.entities {
		for x := 0; x < s.entities[i].Geography.Width; x++ {
			for y := 0; y < s.entities[i].Geography.Height; y++ {
				_ = s.entities[i].Geography.Tiles[x][y].Aquifer
			}
		}
	}
}
