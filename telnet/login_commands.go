package telnet

import (
	"errors"
)

type LoginCommand func(conn *connection, args []string) error

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
