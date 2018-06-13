package model_test

import (
	"testing"

	"github.com/soupstore/coda/simulation/model"
)

func Test_EquipBackpack_NothingEquipped(t *testing.T) {
	rig := model.Rig{}

	backpackDef := model.NewItemDefinition(0, "Test Backpack", []string{}, model.RigSlotBackpack, &model.ContainerDefinition{})
	backpack := backpackDef.Spawn(0)
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
	rig.Backpack = backpackDef.Spawn(1)
	newBackpack := backpackDef.Spawn(2)
	oldBackpack, err := rig.Equip(newBackpack)
	if err != nil {
		t.Error("Failed to equip backpack")
	}
	if rig.Backpack != newBackpack {
		t.Error("Backpack was not equipped")
	}
	if oldBackpack.ID != 1 {
		t.Error("Did not get pointer to old backpack")
	}
}
