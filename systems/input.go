package systems

import (
	"github.com/soupstore/coda/components"
	"github.com/soupstore/ecs"
)

type inputEntity struct {
	*ecs.BasicEntity
	*components.CommandQueue
}

type Input struct {
	entities []inputEntity
}

func (s *Input) Add(basic *ecs.BasicEntity, commandQueue *components.CommandQueue) {
	s.entities = append(s.entities, inputEntity{basic, commandQueue})
}

func (s *Input) Remove(basic ecs.BasicEntity) {
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

func (s *Input) Update(dt float32) {
	for i := range s.entities {
		_ = float64(s.entities[i].ID()+uint64(3)) * 0.5
	}
}
