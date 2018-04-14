package model_test

import (
	"testing"

	"github.com/soupstore/coda-world/simulation/model"
)

func Test_EquipBackpack_NothingEquipped(t *testing.T) {
	rig := model.Rig{}

	backpack := model.NewBackpack(0, "Test Backpack")
	oldBackpack := rig.EquipBackpack(backpack)
	if rig.Backpack != backpack {
		t.Error("Backpack was not equipped")
	}
	if oldBackpack != nil {
		t.Error("Old backpack was not nil")
	}
}

func Test_EquipBackpack_ReplacesCurrent(t *testing.T) {
	rig := model.Rig{}
	rig.Backpack = model.NewBackpack(0, "Old Backpack")

	newBackpack := model.NewBackpack(1, "New Backpack")
	oldBackpack := rig.EquipBackpack(newBackpack)
	if rig.Backpack != newBackpack {
		t.Error("Backpack was not equipped")
	}
	if oldBackpack.ID != 0 {
		t.Error("Did not get pointer to old backpack")
	}
}
