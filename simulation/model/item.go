package model

import (
	"strings"

	"github.com/google/uuid"
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

type ItemID string

type Item struct {
	ID         ItemID
	Definition *ItemDefinition
	Container  Container
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

func (b *ItemDefinition) Spawn() *Item {
	var container Container
	if b.Container != nil {
		container = NewItemContainer()
	}
	return &Item{
		ID:         ItemID(uuid.New().String()),
		Definition: b,
		Container:  container,
	}
}

func (b *Item) KnownAs(alias string) bool {
	if strings.ToLower(alias) == strings.ToLower(b.Definition.Name) {
		return true
	}

	for _, al := range b.Definition.Aliases {
		if strings.ToLower(alias) == strings.ToLower(al) {
			return true
		}
	}

	return false
}
