package telnet

import (
	"github.com/soupstore/coda-world/config"
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
	s.conn.writelnString("Contacting login servers...")

	// client := clients.NewLoginClient(s.config.LoginAddress)
	// characterID, err := client.Login(s.username, s.password)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	s.conn.writeString("Incorrect login")
	// 	s.conn.close()
	// 	return
	// }

	s.conn.ctx = WithCharacterID(s.conn.ctx, 1)

	s.conn.loadState(&stateWorld{
		conn:   s.conn,
		config: s.config,
	})
}
