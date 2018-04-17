package model_test

import (
	"testing"

	"github.com/soupstore/coda-world/simulation/model"
)

func Test_EquipBackpack_NothingEquipped(t *testing.T) {
	rig := model.Rig{}

	backpack := model.NewItem(0, 0, "Test Backpack", []string{}, model.RigSlotBackpack)
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
	rig.Backpack = model.NewItem(0, 0, "Old Backpack", []string{}, model.RigSlotBackpack)

	newBackpack := model.NewItem(1, 0, "New Backpack", []string{}, model.RigSlotBackpack)
	oldBackpack, err := rig.Equip(newBackpack)
	if err != nil {
		t.Error("Failed to equip backpack")
	}
	if rig.Backpack != newBackpack {
		t.Error("Backpack was not equipped")
	}
	if oldBackpack.ID != 0 {
		t.Error("Did not get pointer to old backpack")
	}
}
