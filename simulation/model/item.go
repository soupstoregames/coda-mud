package model

import "strings"

type ItemID int64

type Item struct {
	ID        ItemID
	Name      string
	Aliases   []string
	RigSlot   RigSlot
	Container *Container
}

func NewItem(id ItemID, name string, aliases []string, RigSlot RigSlot) *Item {
	return &Item{
		ID:      id,
		Name:    name,
		Aliases: append(aliases, name),
		RigSlot: RigSlot,
	}
}

func (b *Item) KnownAs(alias string) bool {
	if strings.ToLower(alias) == strings.ToLower(b.Name) {
		return true
	}

	for _, al := range b.Aliases {
		if strings.ToLower(alias) == strings.ToLower(al) {
			return true
		}
	}

	return false
}
