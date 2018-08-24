package telnet

import (
	"errors"
	"strings"

	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/coda/simulation/model"
)

// LoginCommand is a function alias for commands to be used in the login state.
type LoginCommand func(conn *connection, args []string) error

// CmdConnect is the command used to login to the MUD.
func CmdConnect(conn *connection, args []string) error {
	if len(args) != 2 {
		return errors.New("incorrect number of arguments")
	}

	username := args[0]
	password := args[1]

	characterID, ok := conn.usersManager.Login(username, password)
	if !ok {
		return errors.New("invalid login")
	}

	conn.ctx = WithCharacterID(conn.ctx, characterID)

	conn.loadState(&stateWorld{
		conn:   conn,
		config: conn.config,
	})

	return nil
}

// all of the commands available to be used in the world state.
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

// WorldCommand is a function alias for commands to be used in the world state.
type WorldCommand func(model.CharacterID, simulation.CharacterController, []string) error

// CmdDrop allows the character to drop an item from their inventory on to the floor.
func CmdDrop(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	alias := strings.Join(args, " ")
	return cc.DropItem(characterID, alias)
}

// CmdEquip allows the character to equip an item to his rig.
func CmdEquip(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	alias := strings.Join(args, " ")
	return cc.EquipItem(characterID, alias)
}

// CmdLook will trigger another description of the room the character is currently in.
func CmdLook(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Look(characterID)
}

// CmdQuit sends the character to sleep and disconnects the user.
func CmdQuit(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.SleepCharacter(characterID)
}

// CmdSay makes the character speak to all other charactes in the same room.
func CmdSay(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	content := strings.Join(args, " ")
	return cc.Say(characterID, content)
}

// CmdTake has the character pick up an item from the room and put it into their inventory.
func CmdTake(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	alias := strings.Join(args, " ")
	return cc.TakeItem(characterID, alias)
}

// CmdNorth attempts to move the character through the north exit.
func CmdNorth(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionNorth)
}

// CmdNorthEast attempts to move the character through the north east exit.
func CmdNorthEast(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionNorthEast)
}

// CmdEast attempts to move the character through the east exit.
func CmdEast(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionEast)
}

// CmdSouthEast attempts to move the character through the south east exit.
func CmdSouthEast(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionSouthEast)
}

// CmdSouth attempts to move the character through the south exit.
func CmdSouth(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionSouth)
}

// CmdSouthWest attempts to move the character through the south west exit.
func CmdSouthWest(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionSouthWest)
}

// CmdWest attempts to move the character through the west exit.
func CmdWest(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionWest)
}

// CmdNorthWest attempts to move the character through the north west exit.
func CmdNorthWest(characterID model.CharacterID, cc simulation.CharacterController, args []string) error {
	return cc.Move(characterID, model.DirectionNorthWest)
}
