package telnet

import (
	"strconv"
	"strings"

	"github.com/soupstoregames/coda-mud/simulation"
	"github.com/soupstoregames/coda-mud/simulation/model"
)

//// LoginCommand is a function alias for commands to be used in the login state.
//type LoginCommand func(conn *connection, args []string) error
//
//var loginCommands = map[string]LoginCommand{
//	"connect":  CmdConnect,
//	"register": CmdRegister,
//}

// Command is a function alias for commands in the game state.
type Command func(characterID model.CharacterID, sim *simulation.Simulation, args []string) error

// all the commands available to be used in the world state.
var commands = map[string]Command{
	"@spawn":    CmdAdminSpawn,
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

// CmdAdminSpawn allows admins to @spawn in items into the world.
func CmdAdminSpawn(characterID model.CharacterID, sim *simulation.Simulation, args []string) error {
	switch args[0] {
	case "item":
		sItemDefinitionID := args[1]
		itemDefinitionID, err := strconv.ParseInt(sItemDefinitionID, 10, 64)
		if err != nil {
			return err
		}
		return sim.AdminSpawnItem(characterID, model.ItemDefinitionID(itemDefinitionID))
	}
	return nil
}

//// CmdConnect is the command used to login to the MUD.
//func CmdConnect(conn *connection, args []string) error {
//	if len(args) != 2 {
//		return errors.New("incorrect number of arguments")
//	}
//
//	username := args[0]
//	password := args[1]
//
//	characterID, ok := conn.usersManager.Login(username, password)
//	if !ok {
//		return errors.New("invalid login")
//	}
//
//	if characterID == ("") {
//		conn.loadState(&stateCharacterCreation{
//			conn:     conn,
//			config:   conn.config,
//			username: username,
//		}, true)
//
//		return nil
//	}
//
//	conn.ctx = WithCharacterID(conn.ctx, characterID)
//
//	conn.loadState(&stateWorld{
//		conn:   conn,
//		config: conn.config,
//	}, true)
//
//	return nil
//}

//// CmdRegister is used to register a new user.
//func CmdRegister(conn *connection, args []string) error {
//	if len(args) != 2 {
//		return errors.New("incorrect number of arguments")
//	}
//
//	username := args[0]
//	password := args[1]
//
//	if conn.usersManager.IsUsernameTaken(username) {
//		conn.writelnString("That username is taken.")
//		conn.writelnString("If this is your account type 'connect " + username + " <password>'.")
//		conn.writeln()
//		conn.writePrompt()
//		return nil
//	}
//
//	if err := conn.usersManager.Register(username, password); err != nil {
//		logging.Error(err.Error())
//	}
//
//	conn.loadState(&stateCharacterCreation{
//		conn:     conn,
//		config:   conn.config,
//		username: username,
//	}, true)
//
//	return nil
//}

// all of the commands available to be used in the world state.

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

// CmdSay makes the character speak to all other characters in the same room.
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
