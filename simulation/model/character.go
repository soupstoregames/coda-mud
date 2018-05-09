package model

type CharacterID int64

type Character struct {
	*Rig
	ID     CharacterID
	Name   string
	Awake  bool
	Room   *Room
	Events chan interface{}
}

func NewCharacter(id CharacterID, name string, room *Room) *Character {
	return &Character{
		ID:   id,
		Name: name,
		Room: room,
		Rig:  &Rig{},
	}
}

func (c *Character) WakeUp() {
	c.Events = make(chan interface{}, 1)
	c.Awake = true
}

func (c *Character) Sleep() {
	c.Awake = false
	close(c.Events)
}

func (c *Character) Dispatch(event interface{}) {
	if !c.Awake {
		return
	}
	c.Events <- event
}
