package model

import (
	"github.com/google/uuid"
)

type CharacterID string

type Character struct {
	*Rig
	ID     CharacterID
	Name   string
	Awake  bool
	Room   *Room
	Events chan interface{}
}

func NewCharacter(name string, room *Room) *Character {
	return &Character{
		ID:   CharacterID(uuid.New().String()),
		Name: name,
		Room: room,
		Rig:  &Rig{},
	}
}

func (c *Character) WakeUp() {
	c.Events = make(chan interface{}, 10)
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

func (c *Character) TakeItem(item *Item) bool {
	if c.Rig.Backpack != nil {
		c.Rig.Backpack.Container.PutItem(item)
		return true
	}

	return false
}
