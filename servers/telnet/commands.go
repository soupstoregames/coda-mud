package telnet

import (
	"errors"
	"github.com/soupstore/go-core/logging"
	"strconv"
	"strings"

	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/coda/simulation/model"
)

// AdminCommand is a function alias for commands that only admins can use.
type AdminCommand func(model.CharacterID, simulation.AdminController, []string) error

// all of the commands available to be used in the world state.
var adminCommands = map[string]AdminCommand{
	"@spawn": CmdAdminSpawn,
}

// CmdAdminSpawn allows admins to spawn in items into the world.
func CmdAdminSpawn(characterID model.CharacterID, ac simulation.AdminController, args []string) error {
	switch args[0] {
	case "item":
		sItemDefinitionID := args[1]
		itemDefinitionID, err := strconv.ParseInt(sItemDefinitionID, 10, 64)
		if err != nil {
			return err
		}
		return ac.AdminSpawnItem(characterID, model.ItemDefinitionID(itemDefinitionID))
	}
	return nil
}

// LoginCommand is a function alias for commands to be used in the login state.
type LoginCommand func(conn *connection, args []string) error

var loginCommands = map[string]LoginCommand{
	"connect":  CmdConnect,
	"register": CmdRegister,
}

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

	if characterID == model.CharacterID("") {
		conn.loadState(&stateCharacterCreation{
			conn:     conn,
			config:   conn.config,
			username: username,
		})

		return nil
	}

	conn.ctx = WithCharacterID(conn.ctx, characterID)

	conn.loadState(&stateWorld{
		conn:   conn,
		config: conn.config,
	})

	return nil
}

// CmdRegister is used to register a new user.
func CmdRegister(conn *connection, args []string) error {
	if len(args) != 2 {
		return errors.New("incorrect number of arguments")
	}

	username := args[0]
	password := args[1]

	if conn.usersManager.IsUsernameTaken(username) {
		conn.writelnString("That username is taken.")
		conn.writelnString("If this is your account type 'connect " + username + " <password>'.")
		conn.writelnString()
		conn.writePrompt()
		return nil
	}

	if err := conn.usersManager.Register(username, password); err != nil {
		logging.Error(err.Error())
	}

	conn.loadState(&stateCharacterCreation{
		conn:     conn,
		config:   conn.config,
		username: username,
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
	"up":        CmdUp,
	"u":         CmdUp,
	"south":     CmdSouth,
	"s":         CmdSouth,
	"southwest": CmdSouthWest,
	"sw":        CmdSouthWest,
	"west":      CmdWest,
	"w":         CmdWest,
	"northwest": CmdNorthWest,
	"nw":        CmdNorthWest,
	"down":      CmdDown,
	"d":         CmdDown,
	"take":      CmdTake,
	"get":       CmdTake,
	"drop":      CmdDrop,
	"equip":     CmdEquip,
	"wear":      CmdEquip,
	"unequip":   CmdUnequip,
	"remove":    CmdUnequip,
	"inventory": CmdInventory,
	"i":         CmdInventory,
}

// WorldCommand is a function alias for commands to be used in the world state.
type WorldCommand func(model.CharacterID, *simulation.Simulation, []string) error

// CmdLook will trigger another description of the room the character is currently in.
func CmdLook(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.Look(characterID)
}

// CmdInventory lists the character's current equipment and items in containers.
func CmdInventory(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.Inventory(characterID)
}

// CmdQuit sends the character to sleep and disconnects the user.
func CmdQuit(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.SleepCharacter(characterID)
}

// CmdSay makes the character speak to all other charactes in the same room.
func CmdSay(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandSay{
		Content: strings.Join(args, " "),
	})
}

// CmdTake has the character pick up an item from the room and put it into their inventory.
func CmdTake(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	item, err := cc.FindItemInRoom(characterID, strings.Join(args, " "))
	if err != nil {
		return err
	}
	return cc.QueueCommand(characterID, model.CommandTake{
		Item: item,
	})
}

// CmdDrop allows the character to drop an item from their inventory on to the floor.
func CmdDrop(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	item, err := cc.FindItemInInventory(characterID, strings.Join(args, " "))
	if err != nil {
		return err
	}

	return cc.QueueCommand(characterID, model.CommandDrop{
		Item: item,
	})
}

// CmdEquip allows the character to equip an item to his rig.
func CmdEquip(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	item, err := cc.FindItemInInventory(characterID, strings.Join(args, " "))
	if err != nil {
		return err
	}

	return cc.QueueCommand(characterID, model.CommandEquip{
		Item: item,
	})
}

// CmdUnequip takes an item off the character's rig and stores it.
func CmdUnequip(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	item, err := cc.FindItemInRig(characterID, strings.Join(args, " "))
	if err != nil {
		return err
	}

	return cc.QueueCommand(characterID, model.CommandUnequip{
		Item: item,
	})
}

// CmdNorth attempts to move the character through the north exit.
func CmdNorth(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionNorth,
	})
}

// CmdNorthEast attempts to move the character through the north east exit.
func CmdNorthEast(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionNorthEast,
	})
}

// CmdEast attempts to move the character through the east exit.
func CmdEast(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionEast,
	})
}

// CmdSouthEast attempts to move the character through the south east exit.
func CmdSouthEast(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionSouthEast,
	})
}

// CmdUp attempts to move the character through the up exit.
func CmdUp(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionUp,
	})
}

// CmdSouth attempts to move the character through the south exit.
func CmdSouth(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionSouth,
	})
}

// CmdSouthWest attempts to move the character through the south west exit.
func CmdSouthWest(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionSouthWest,
	})
}

// CmdWest attempts to move the character through the west exit.
func CmdWest(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionWest,
	})
}

// CmdNorthWest attempts to move the character through the north west exit.
func CmdNorthWest(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionNorthWest,
	})
}

// CmdDown attempts to move the character through the down exit.
func CmdDown(characterID model.CharacterID, cc *simulation.Simulation, args []string) error {
	return cc.QueueCommand(characterID, model.CommandMove{
		Direction: model.DirectionDown,
	})
}
