package state

import (
	"github.com/go-pg/pg"
	"github.com/soupstore/coda/simulation/model"
)

type Character struct {
	ID    model.CharacterID
	Name  string
	Room  model.RoomID
	World model.WorldID
}

func GetCharacters(db *pg.DB) ([]*Character, error) {
	var characters []*Character
	if err := db.Model(&characters).Select(); err != nil {
		return nil, err
	}
	return characters, nil
}

func UpdateCharacterLocation(db *pg.DB, id model.CharacterID, room model.RoomID, world model.WorldID) error {
	character := &Character{ID: id}
	if err := db.Select(character); err != nil {
		return err
	}

	character.Room = room
	character.World = world

	if err := db.Update(character); err != nil {
		return err
	}

	return nil
}
