package model

type CharacterID int64

type Character struct {
	Rig
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
		Rig: Rig{
			Backpack: NewItem(0, 99, "CODA Recon Pack", []string{"backpack", "pack", "recon pack", "coda recon pack", "coda recon"}, RigSlotBackpack),
		},
	}
}

func (c *Character) Dispatch(event interface{}) {
	if !c.Awake {
		return
	}
	c.Events <- event
}
