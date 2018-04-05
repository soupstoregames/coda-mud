package model

type CharacterID int64

type Character struct {
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
	}
}

func (c *Character) Dispatch(event interface{}) {
	if !c.Awake {
		return
	}
	c.Events <- event
}
