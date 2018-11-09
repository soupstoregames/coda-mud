package telnet

import (
	"errors"
	"strings"

	"github.com/aybabtme/rgbterm"

	"github.com/soupstore/coda/config"
	"github.com/soupstore/coda/simulation/model"
	"github.com/soupstore/go-core/logging"
)

// state is the interface for all scenes in this package
type state interface {
	onEnter() error
	onExit() error
	handleInput(string) error
}

// stateLogin is the scene used for connecting to a character
type stateLogin struct {
	config   *config.Config
	conn     *connection
	username string
	password string
}

// onEnter is called when the scene is first loaded
func (s *stateLogin) onEnter() error {
	s.conn.writelnString(
		`                    .___ ` + "\r\n" +
			`    _____  __ __  __| _/ ` + "\r\n" +
			`   /     \|  |  \/ __ |  ` + "\r\n" +
			`  |  Y Y  \  |  / /_/ |  ` + "\r\n" +
			`  |__|_|  /____/\____ |  ` + "\r\n" +
			`        \/           \/  `)
	s.conn.writelnString()
	s.conn.writelnString("Type 'connect <username> <password>' to log in.")
	s.conn.writelnString("Type 'register <username> <password>' to create a new account.")

	s.conn.writePrompt()

	return nil
}

func (s *stateLogin) onExit() error {
	return nil
}

func (s *stateLogin) handleInput(input string) error {
	tokens := strings.Split(input, " ")
	commandText := strings.ToLower(tokens[0])

	if commandText == "quit" {
		s.conn.close()
		return errors.New("closed")
	}

	command, ok := loginCommands[commandText]
	if !ok {
		echo := rgbterm.String("Huh?", 255, 100, 100, 0, 0, 0)
		s.conn.writelnString(echo)
		s.conn.writePrompt()
		return nil
	}

	err := command(s.conn, tokens[1:])
	if err != nil {
		s.conn.writelnString(err.Error())
		s.conn.writePrompt()
	}

	return nil
}

type characterCreationPhase byte

const (
	characterCreationPhaseName characterCreationPhase = iota
)

type stateCharacterCreation struct {
	config   *config.Config
	conn     *connection
	phase    characterCreationPhase
	username string
	name     string
}

func (s *stateCharacterCreation) onEnter() error {
	s.phase = characterCreationPhaseName
	s.conn.writelnString("CHARACTER CREATION")
	s.conn.writelnString("What will you be known as?")
	s.conn.writePrompt()

	return nil
}

func (s *stateCharacterCreation) onExit() error {
	return nil
}

// handleInput parses input from the client and performs any appropriate command
func (s *stateCharacterCreation) handleInput(input string) error {
	switch s.phase {
	case characterCreationPhaseName:
		s.name = input
		id := s.conn.sim.MakeCharacter(s.name)
		if err := s.conn.usersManager.AssociateCharacter(s.username, id); err != nil {
			return err
		}

		s.conn.ctx = WithCharacterID(s.conn.ctx, id)

		s.conn.loadState(&stateWorld{
			conn:   s.conn,
			config: s.config,
		})
	}

	return nil
}

// stateWorld is the scene used interacting with the world
type stateWorld struct {
	config      *config.Config
	conn        *connection
	characterID model.CharacterID
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
	logging.Debug("Disconnecting from world server")
	// err := s.cc.Close()
	logging.Debug("Disconnected from world server")
	return nil
}

// handleInput parses input from the client and performs any appropriate command
func (s *stateWorld) handleInput(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	tokens := strings.Split(input, " ")
	commandText := strings.ToLower(tokens[0])

	if commandText[0] == '@' {
		command, ok := adminCommands[commandText]
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
	} else {
		command, ok := worldCommands[commandText]
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
	}

	// TODO: I dont like this - need to fix it
	if commandText == "quit" {
		s.conn.close()
		return errors.New("closed")
	}

	return nil
}
