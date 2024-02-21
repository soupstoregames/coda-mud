package simulation

import (
	"fmt"
	"github.com/soupstoregames/coda-mud/simulation/model"
	"github.com/soupstoregames/go-core/logging"
)

func (s *Simulation) AdminSpawnItem(characterID model.CharacterID, id model.ItemDefinitionID) error {
	actor, err := s.findAwakeCharacter(characterID)
	if err != nil {
		return err
	}

	definition, ok := s.itemDefinitions[id]
	if !ok {
		logging.Warn(fmt.Sprintf("Tried to load item for non-existant definition %d in room %d in world %s", id, actor.Room.ID, actor.Room.WorldID))
		return nil // TODO RETURN ERROR
	}

	item := definition.Spawn()
	actor.Room.Container.PutItem(item)

	spawnEvent := model.EvtAdminSpawnsItem{Character: actor, Item: item}
	for _, c := range actor.Room.Characters {
		c.Dispatch(spawnEvent)
	}

	return nil
}
