package telnet

import (
	"github.com/soupstore/coda-world/common/config"
	"github.com/soupstore/coda-world/login"
)

// stateLogin is the scene used for connecting to a character
type stateLogin struct {
	config   *config.Config
	conn     *connection
	substate byte
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
	s.promptUsername()

	return nil
}

func (s *stateLogin) onExit() error {
	return nil
}

func (s *stateLogin) handleInput(input string) error {
	switch s.substate {
	case 0: // waiting for username
		s.username = input
		s.promptPassword()
	case 1: // waiting for password
		s.password = input
		s.attemptLogin()
	}

	return nil
}

func (s *stateLogin) promptUsername() {
	s.conn.writeString("What is your name? ")
	s.substate = 0
}

func (s *stateLogin) promptPassword() {
	s.conn.writeString("What is your password? ")
	s.substate = 1
}

func (s *stateLogin) attemptLogin() {
	characterID, ok := login.GetCharacter(s.username, s.password)
	if !ok {
		s.conn.writelnString("Incorrect login")
		s.conn.close()
		return
	}

	s.conn.ctx = WithCharacterID(s.conn.ctx, characterID)

	s.conn.loadState(&stateWorld{
		conn:   s.conn,
		config: s.config,
	})
}
