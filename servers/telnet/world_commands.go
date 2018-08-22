package telnet

import (
	"strings"

	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/coda/simulation/model"
)

var worldCommands = map[string]WorldCommand{
	"look":      CmdLook,
	"l":         CmdLook,
	"say":       CmdSay,
	"quit":      CmdQuit,
	"north":     CmdNorth,
	"n":         CmdNorth,
	"northeast": CmdNorthEast,
	"ne":        CmdNorthEast,
	"east":      CmdEast,
	"e":         CmdEast,
	"southeast": CmdSouthEast,
	"se":        CmdSouthEast,
	"south":     CmdSouth,
	"s":         CmdSouth,
	"southwest": CmdSouthWest,
	"sw":        CmdSouthWest,
	"west":      CmdWest,
	"w":         CmdWest,
	"northwest": CmdNorthWest,
	"nw":        CmdNorthWest,
	"take":      CmdTake,
	"get":       CmdTake,
	"drop":      CmdDrop,
	"equip":     CmdEquip,
	"wear":      CmdEquip,
}

type WorldCommand func(model.CharacterID, simulation.CharacterController, []string) error

func CmdDrop(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	alias := strings.Join(args, " ")
	return cc.DropItem(characterID, alias)
}

func CmdEquip(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	alias := strings.Join(args, " ")
	return cc.EquipItem(characterID, alias)
}

func CmdLook(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Look(characterID)
}

func CmdQuit(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.SleepCharacter(characterID)
}

func CmdSay(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	content := strings.Join(args, " ")
	return cc.Say(characterID, content)
}

func CmdTake(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	alias := strings.Join(args, " ")
	return cc.TakeItem(characterID, alias)
}

func CmdNorth(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionNorth)
}

func CmdNorthEast(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionNorthEast)
}

func CmdEast(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionEast)
}

func CmdSouthEast(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionSouthEast)
}

func CmdSouth(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionSouth)
}

func CmdSouthWest(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionSouthWest)
}

func CmdWest(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionWest)
}

func CmdNorthWest(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionNorthWest)
}
