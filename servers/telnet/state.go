package telnet

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aybabtme/rgbterm"
	"github.com/soupstoregames/coda-mud/config"
	"github.com/soupstoregames/coda-mud/simulation/model"
	"github.com/soupstoregames/go-core/logging"
	"strings"
)

// state is the interface for all scenes in this package
type state interface {
	onEnter() error
	onExit() error
	handleInput(byte) error
	writePrompt()
}

type stateMainMenu struct {
	config *config.Config
	conn   *connection
}

// onEnter is called when the scene is first loaded
func (s *stateMainMenu) onEnter() error {
	s.conn.writelnString(`
                  .___
  _____  __ __  __| _/
 /     \|  |  \/ __ |
|  Y Y  \  |  / /_/ |
|__|_|  /____/\____ |
      \/           \/
`)
	s.conn.writeln()
	s.conn.writelnString("1: Login")
	s.conn.writelnString("2: Register")
	s.conn.writeln()
	s.conn.writelnString("0: Quit")
	s.conn.writeln()

	return nil
}

func (s *stateMainMenu) onExit() error {
	s.conn.close()
	return nil
}

func (s *stateMainMenu) handleInput(b byte) error {
	fmt.Println(b)
	switch b {
	case '0':
		s.conn.close()
	case '1':
		s.conn.state.Push(&stateLogin{
			config: s.config,
			conn:   s.conn,
		})
		s.conn.state.Peek().onEnter()
	case '2':
		s.conn.state.Push(&stateRegister{
			config: s.config,
			conn:   s.conn,
		})
		s.conn.state.Peek().onEnter()
	}
	return nil
}

func (s *stateMainMenu) writePrompt() {
	s.conn.write([]byte{charLF, charCR})
}

// stateLogin is the scene used for connecting to a character
type stateLogin struct {
	config *config.Config
	conn   *connection

	input bytes.Buffer

	subState byte
	attempts byte

	username string
	password string
}

// onEnter is called when the scene is first loaded
func (s *stateLogin) onEnter() error {
	s.writePrompt()
	return nil
}

func (s *stateLogin) onExit() error {
	return nil
}

func (s *stateLogin) handleInput(input byte) error {
	switch input {
	case charCR:
		// do nothing
	case charNULL:
		fallthrough
	case charLF:
		s.conn.writeln()
		switch s.subState {
		case 0:
			s.username = s.input.String()
			s.input.Reset()
			s.subState++
			s.writePrompt()
		case 1:
			s.password = s.input.String()
			s.input.Reset()

			characterID, ok := s.conn.usersManager.Login(s.username, s.password)
			if !ok {
				s.conn.writelnString("Invalid login.")
				s.attempts++

				if s.attempts >= 3 {
					s.conn.writelnString("Too many attempts.")
					s.conn.close()
					return nil
				}

				s.username = ""
				s.password = ""
				s.subState = 0
				s.writePrompt()
				return nil
			}

			if characterID == ("") {
				s.conn.state.Push(&stateCharacterCreation{
					conn:     s.conn,
					config:   s.config,
					username: s.username,
				})
				s.conn.state.Peek().onEnter()

				return nil
			}

			s.conn.ctx = WithCharacterID(s.conn.ctx, characterID)

			s.conn.state.Pop()
			s.conn.state.Push(&stateWorld{
				conn:   s.conn,
				config: s.config,
			})
			s.conn.state.Peek().onEnter()
		}
	case charDELETE:
		if s.input.Len() > 0 {
			s.input.Truncate(s.input.Len() - 1)
			s.conn.write([]byte{8, 32, 8})
		}
	default:
		s.input.WriteByte(input)
		switch s.subState {
		case 0:
			s.conn.write([]byte{input})
		case 1:
			s.conn.write([]byte{'*'})
		}
	}

	return nil
}

func (s *stateLogin) writePrompt() {
	switch s.subState {
	case 0:
		s.conn.writeString("username: ")
	case 1:
		s.conn.writeString("password: ")
	}
}

type stateRegister struct {
	config *config.Config
	conn   *connection

	input bytes.Buffer

	subState byte

	username string
	password string
}

// onEnter is called when the scene is first loaded
func (s *stateRegister) onEnter() error {
	s.writePrompt()
	return nil
}

func (s *stateRegister) onExit() error {
	return nil
}

func (s *stateRegister) handleInput(input byte) error {
	switch input {
	case charCR:
		// do nothing
	case charNULL:
		fallthrough
	case charLF:
		s.conn.writeln()
		switch s.subState {
		case 0:
			s.username = s.input.String()
			s.input.Reset()

			if s.conn.usersManager.IsUsernameTaken(s.username) {
				s.conn.writelnString("Username already taken.")
			} else {
				s.subState++
			}
			s.writePrompt()
		case 1:
			s.password = s.input.String()
			s.input.Reset()

			err := s.conn.usersManager.Register(s.username, s.password)
			if err != nil {
				logging.Error(err.Error())
				s.conn.writelnString("Failed to create account.")

				s.username = ""
				s.password = ""
				s.subState = 0
				s.writePrompt()
				return nil
			}

			s.conn.state.Push(&stateCharacterCreation{
				conn:     s.conn,
				config:   s.config,
				username: s.username,
			})
			s.conn.state.Peek().onEnter()

			return nil
		}
	case charDELETE:
		if s.input.Len() > 0 {
			s.input.Truncate(s.input.Len() - 1)
			s.conn.write([]byte{8, 32, 8})
		}
	default:
		s.input.WriteByte(input)
		switch s.subState {
		case 0:
			s.conn.write([]byte{input})
		case 1:
			s.conn.write([]byte{'*'})
		}
	}

	return nil
}

func (s *stateRegister) writePrompt() {
	switch s.subState {
	case 0:
		s.conn.writeString("username: ")
	case 1:
		s.conn.writeString("password: ")
	}
}

type stateCharacterCreation struct {
	config   *config.Config
	conn     *connection
	username string

	input bytes.Buffer
	phase byte

	name string
}

func (s *stateCharacterCreation) onEnter() error {
	s.phase = 0
	s.conn.writelnString("CHARACTER CREATION")
	s.conn.writelnString("What will you be known as?")
	s.writePrompt()

	return nil
}

func (s *stateCharacterCreation) onExit() error {
	return nil
}

func (s *stateCharacterCreation) writePrompt() {
	s.conn.write([]byte{charLF})
	s.conn.write([]byte{'>', ' '})
}

// handleInput parses input from the client and performs any appropriate command
func (s *stateCharacterCreation) handleInput(input byte) error {
	switch input {
	case charCR:
		// do nothing
	case charNULL:
		fallthrough
	case charLF:
		id := s.conn.sim.MakeCharacter(s.input.String())
		if err := s.conn.usersManager.AssociateCharacter(s.username, id); err != nil {
			return err
		}

		s.conn.ctx = WithCharacterID(s.conn.ctx, id)

		s.conn.state.Pop()
		s.conn.state.Push(&stateWorld{
			conn:   s.conn,
			config: s.config,
		})
		s.conn.state.Peek().onEnter()
	case charDELETE:
		if s.input.Len() > 0 {
			s.input.Truncate(s.input.Len() - 1)
			s.conn.write([]byte{8, 32, 8})
		}
	default:
		s.input.WriteByte(input)
		s.conn.write([]byte{input})
	}

	return nil
}

// stateWorld is the scene used interacting with the world
type stateWorld struct {
	config      *config.Config
	conn        *connection
	characterID model.CharacterID

	input bytes.Buffer
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
func (s *stateWorld) handleInput(input byte) error {
	switch input {
	case charCR:
		// do nothing
	case charNULL:
		fallthrough
	case charLF:
		s.conn.write([]byte{charLF})

		cleansed := strings.TrimSpace(s.input.String())
		s.input.Reset()

		if cleansed == "" {
			return nil
		}

		tokens := strings.Split(cleansed, " ")
		commandText := strings.ToLower(tokens[0])

		// TODO: I dont like this - need to fix it
		if commandText == "quit" {
			s.conn.close()
			return errors.New("closed")
		}

		if commandText[0] == '@' {
			// TODO: check admin
		}

		command, ok := commands[commandText]
		if !ok {
			echo := rgbterm.String("Huh?", 255, 100, 100, 0, 0, 0)
			s.conn.writelnString(echo)
			s.writePrompt()
			return nil
		}

		characterID := CharacterIDFromContext(s.conn.ctx)
		err := command(characterID, s.conn.sim, tokens[1:])
		if err != nil {
			echo := rgbterm.String(err.Error(), 255, 100, 100, 0, 0, 0)
			s.conn.writelnString(echo)
			s.writePrompt()
			return nil
		}

	case charDELETE:
		if s.input.Len() > 0 {
			s.input.Truncate(s.input.Len() - 1)
			s.conn.write([]byte{8, 32, 8})
		}
	default:
		s.input.WriteByte(input)
		s.conn.write([]byte{input})
	}

	return nil
}

func (s *stateWorld) writePrompt() {
	s.conn.write([]byte{charLF})
	s.conn.write([]byte{'>', ' '})
}
