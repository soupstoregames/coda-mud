package telnet

import (
	"errors"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/soupstore/coda/config"
)

// stateLogin is the scene used for connecting to a character
type stateLogin struct {
	config   *config.Config
	conn     *connection
	username string
	password string
}

var loginCommands = map[string]LoginCommand{
	"connect": CmdConnect,
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
