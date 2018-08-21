package model_test

import (
	"testing"

	"github.com/soupstore/coda/simulation/model"
)

func Test_EquipBackpack_NothingEquipped(t *testing.T) {
	rig := model.Rig{}

	backpackDef := model.NewItemDefinition(0, "Test Backpack", []string{}, model.RigSlotBackpack, &model.ContainerDefinition{})
	backpack := backpackDef.Spawn()
	oldBackpack, err := rig.Equip(backpack)
	if err != nil {
		t.Error("Failed to equip backpack")
	}
	if rig.Backpack != backpack {
		t.Error("Backpack was not equipped")
	}
	if oldBackpack != nil {
		t.Error("Old backpack was not nil")
	}
}

func Test_EquipBackpack_ReplacesCurrent(t *testing.T) {
	rig := model.Rig{}
	backpackDef := model.NewItemDefinition(0, "Old Backpack", []string{}, model.RigSlotBackpack, &model.ContainerDefinition{})

	originalBackpack := backpackDef.Spawn()
	newBackpack := backpackDef.Spawn()

	rig.Backpack = originalBackpack

	replaced, err := rig.Equip(newBackpack)
	if err != nil {
		t.Error("Failed to equip backpack")
	}
	if rig.Backpack != newBackpack {
		t.Error("Backpack was not equipped")
	}
	if replaced != originalBackpack {
		t.Error("Did not get pointer to original backpack")
	}
}
