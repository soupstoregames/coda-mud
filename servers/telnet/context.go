package telnet

import (
	"context"

	"github.com/soupstore/coda/simulation/model"
)

type key int

const (
	connectionIDkey key = iota
	characterIDKey
)

// WithCharacterID returns the given context with the character ID in it.
func WithCharacterID(parent context.Context, characterID model.CharacterID) context.Context {
	return context.WithValue(parent, characterIDKey, characterID)
}

// CharacterIDFromContext extracts the character ID embedded in the context.
func CharacterIDFromContext(c context.Context) model.CharacterID {
	return c.Value(characterIDKey).(model.CharacterID)
}

// WithConnectionID returns the given context with the connection ID in it.
func WithConnectionID(parent context.Context, connectionID string) context.Context {
	return context.WithValue(parent, characterIDKey, connectionID)
}

// ConnectionIDFromContext extracts the connection ID embedded in the context.
func ConnectionIDFromContext(c context.Context) string {
	return c.Value(connectionIDkey).(string)
}
