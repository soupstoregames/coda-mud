package model

import (
	"github.com/google/uuid"
)

// CharacterID is a type-aliased string, often set to a uuid.
type CharacterID string

// Character is a player character in the simulation.
type Character struct {
	*Rig
	ID       CharacterID
	Name     string
	Awake    bool
	Room     *Room
	Commands chan interface{}
	Events   chan interface{}
}

// NewCharacter is a helper function for creating a new character in the simulation.
// It requires the character's name and the room to spawn the character in.
func NewCharacter(name string, room *Room) *Character {
	return &Character{
		ID:   CharacterID(uuid.New().String()),
		Name: name,
		Room: room,
		Rig:  &Rig{},
	}
}

// WakeUp initializes a buffered channel of simulation events that happen to the character.
// It also wake it up, which allows the player to control it.
func (c *Character) WakeUp() {
	c.Commands = make(chan interface{}, 1)
	c.Events = make(chan interface{}, 10)
	c.Awake = true
}

// Sleep closes the channel of simulation events for this character and puts the character to sleep.
func (c *Character) Sleep() {
	c.Awake = false
	close(c.Events)
	close(c.Commands)
}

// Dispatch is used by the simulation to send events to the character's event stream.
func (c *Character) Dispatch(event interface{}) {
	if !c.Awake {
		return
	}
	c.Events <- event
}

// TakeItem attempts to find a free slot in the player's inventory to place the given item.
func (c *Character) TakeItem(item *Item) bool {
	// currently this is the only place we could put an item and it has no limits!
	if c.Rig.Backpack != nil {
		c.Rig.Backpack.Container.PutItem(item)
		return true
	}

	return false
}

func (c *Character) DropItem(item *Item) bool {
	return c.Rig.RemoveItemFromInventory(item)
}
