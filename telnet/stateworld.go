package telnet

import (
	"errors"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/soupstore/coda/common/config"
	"github.com/soupstore/coda/common/log"
	"github.com/soupstore/coda/simulation/model"
)

// stateWorld is the scene used interacting with the world
type stateWorld struct {
	config      *config.Config
	conn        *connection
	characterID model.CharacterID
}

var validCommands = map[string]Command{
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

// onEnter is called when the scene is first loaded
func (s *stateWorld) onEnter() error {
	s.conn.writelnString("You are in the world!\n\r")

	s.characterID = CharacterIDFromContext(s.conn.ctx)
	events, err := s.conn.sim.WakeUpCharacter(s.characterID)
	if err != nil {
		return err
	}

	go renderEvents(s.conn, events)

	return nil
}

func (s *stateWorld) onExit() error {
	log.Logger().Debug("Disconnecting from world server")
	// err := s.cc.Close()
	log.Logger().Debug("Disconnected from world server")
	return nil
}

// handleInput parses input from the client and performs any appropriate command
func (s *stateWorld) handleInput(input string) error {
	tokens := strings.Split(input, " ")
	commandText := strings.ToLower(tokens[0])

	command, ok := validCommands[commandText]
	if !ok {
		echo := rgbterm.String("Huh?", 255, 100, 100, 0, 0, 0)
		s.conn.writelnString(echo)
		s.conn.writePrompt()
		return nil
	}

	characterID := CharacterIDFromContext(s.conn.ctx)
	err := command(characterID, s.conn.sim, tokens[1:])
	if err != nil {
		return err
	}

	// TODO: I dont like this - need to fix it
	if commandText == "quit" {
		s.conn.close()
		return errors.New("closed")
	}

	return nil
}
