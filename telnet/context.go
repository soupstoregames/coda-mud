package telnet

import (
	"context"

	"github.com/soupstore/coda-world/simulation/model"
)

type key int

const (
	connectionIDkey = iota
	characterIDKey
)

func WithCharacterID(parent context.Context, characterID model.CharacterID) context.Context {
	return context.WithValue(parent, characterIDKey, characterID)
}

func CharacterIDFromContext(c context.Context) model.CharacterID {
	return c.Value(characterIDKey).(model.CharacterID)
}

func WithConnectionID(parent context.Context, connectionID string) context.Context {
	return context.WithValue(parent, characterIDKey, connectionID)
}

func ConnectionIDFromContext(c context.Context) string {
	return c.Value(connectionIDkey).(string)
}
