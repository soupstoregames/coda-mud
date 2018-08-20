package model

import (
	"strings"
)

type ItemDefinitionID int64

type ItemDefinition struct {
	ID        ItemDefinitionID
	Name      string
	Aliases   []string
	RigSlot   RigSlot
	Container *ContainerDefinition
}

type ContainerDefinition struct {
}

type ItemID int64

type Item struct {
	ID        ItemID
	Name      string
	Aliases   []string
	RigSlot   RigSlot
	Container Container
}

func NewItemDefinition(id ItemDefinitionID, name string, aliases []string, RigSlot RigSlot, container *ContainerDefinition) *ItemDefinition {
	return &ItemDefinition{
		ID:        id,
		Name:      name,
		Aliases:   append(aliases, name),
		RigSlot:   RigSlot,
		Container: container,
	}
}

func (b *ItemDefinition) Spawn(itemID ItemID) *Item {
	var container Container
	if b.Container != nil {
		container = NewItemContainer(0)
	}
	return &Item{
		ID:        itemID,
		Name:      b.Name,
		Aliases:   b.Aliases,
		RigSlot:   b.RigSlot,
		Container: container,
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
