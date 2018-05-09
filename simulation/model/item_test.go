package model_test

import (
	"testing"

	"github.com/soupstore/coda-world/simulation/model"
)

func TestNewItemDefinition(t *testing.T) {
	itemDef := model.NewItemDefinition(model.ItemDefinitionID(1), "Test Item", []string{"test", "item"}, model.RigSlotBackpack, &model.ContainerDefinition{})
	item := itemDef.Spawn(model.ItemID(333))

	if item.ID != model.ItemID(333) {
		t.Error("Item has incorrect ID assigned")
	}

	if item.Name != "Test Item" {
		t.Error("Item has been given the wrong name")
	}

	if item.RigSlot != model.RigSlotBackpack {
		t.Error("Item does not fit the request rig slot ")
	}
}

func TestKnownAs(t *testing.T) {
	itemDef := model.NewItemDefinition(model.ItemDefinitionID(1), "Test Item", []string{"test", "item"}, model.RigSlotBackpack, &model.ContainerDefinition{})
	item := itemDef.Spawn(model.ItemID(333))

	if !item.KnownAs("Test Item") {
		t.Error("Item should be known as its name")
	}

	if !item.KnownAs("test ITEm") {
		t.Error("Aliases should be case insensitive")
	}

	if !item.KnownAs("test") {
		t.Error("Item not known by its alias")
	}

	if !item.KnownAs("ITEm") {
		t.Error("Aliases should be case insensitive")
	}

	if !item.KnownAs("item") {
		t.Error("Item not known by its alias")
	}

	if item.KnownAs("something else") {
		t.Error("Item was known as something random")
	}
}
