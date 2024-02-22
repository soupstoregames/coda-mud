package simulation

import (
	"errors"
	"github.com/soupstoregames/coda-mud/simulation/model"
)

func (s *Simulation) say(actor *model.Character, c model.CommandSay) {
	actor.Room.Dispatch(model.EvtCharacterSpeaks{
		Character: actor,
		Content:   c.Content,
	})
}

func (s *Simulation) move(actor *model.Character, c model.CommandMove) {
	// save the actor's current room
	originalRoom := actor.Room

	// look for an exit in the direction the user specified
	exit := actor.Room.Exits[c.Direction]
	if exit == nil {
		actor.Dispatch(model.EvtNoExitInThatDirection{})
		return
	}

	// actor.Room.OnExit(actor)

	newRoom, err := s.GetRoom(exit.WorldID, exit.RoomID)
	if err != nil {
		// room doesn't exist, do something clever!
	}

	// remove actor from current room
	originalRoom.RemoveCharacter(actor)

	// tell people in the room that the actor has left
	if !originalRoom.Alone {
		originalRoom.Dispatch(model.EvtCharacterLeaves{
			Character: actor,
			Direction: c.Direction,
		})
	}

	// move actor to the new room
	actor.Room = newRoom
	newRoom.AddCharacter(actor)
	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})

	// tell people in the target room that a character has arrived
	if !newRoom.Alone {
		for _, ch := range newRoom.Characters {
			if ch == actor {
				continue
			}
			ch.Dispatch(model.EvtCharacterArrives{
				Character: actor,
				Direction: c.Direction.Opposite(),
			})
		}
	}

	actor.Room.OnEnter(actor)
}

func (s *Simulation) takeItem(actor *model.Character, c model.CommandTake) {
	actor.TakeItem(c.Item)
	actor.Room.Container.RemoveItem(c.Item.ID)

	if actor.Room.Alone {
		actor.Dispatch(model.EvtCharacterTakesItem{
			Character: actor,
			Item:      c.Item,
		})
	} else {
		actor.Room.Dispatch(model.EvtCharacterTakesItem{
			Character: actor,
			Item:      c.Item,
		})
	}
}

func (s *Simulation) dropItem(actor *model.Character, c model.CommandDrop) {
	actor.DropItem(c.Item)
	actor.Room.Container.PutItem(c.Item)

	if actor.Room.Alone {
		actor.Dispatch(model.EvtCharacterDropsItem{
			Character: actor,
			Item:      c.Item,
		})
	} else {
		actor.Room.Dispatch(model.EvtCharacterDropsItem{
			Character: actor,
			Item:      c.Item,
		})
	}
}

func (s *Simulation) equipItem(actor *model.Character, c model.CommandEquip) {
	actor.Container.RemoveItem(c.Item.ID)
	oldItem, err := actor.Equip(c.Item)
	if errors.Is(err, model.ErrNotEquipable) {
		// do something
		actor.TakeItem(c.Item)
	}

	// do we have an item that we replaced
	if oldItem != nil {
		// tell everyone that we took off an item
		if actor.Room.Alone {
			actor.Dispatch(model.EvtCharacterUnequipsItem{
				Character: actor,
				Item:      oldItem,
			})
		} else {
			actor.Room.Dispatch(model.EvtCharacterUnequipsItem{
				Character: actor,
				Item:      oldItem,
			})
		}
	}

	// tell everyone that we put on an item
	if actor.Room.Alone {
		actor.Dispatch(model.EvtCharacterEquipsItem{
			Character: actor,
			Item:      c.Item,
		})
	} else {
		actor.Room.Dispatch(model.EvtCharacterEquipsItem{
			Character: actor,
			Item:      c.Item,
		})
	}

	if oldItem != nil {
		actor.TakeItem(oldItem)
	}
}

func (s *Simulation) unequipItem(actor *model.Character, c model.CommandUnequip) {
	if ok := actor.Rig.Unequip(c.Item); !ok {
		// DO SOMETHING
	}

	if actor.Room.Alone {
		actor.Dispatch(model.EvtCharacterUnequipsItem{
			Character: actor,
			Item:      c.Item,
		})
	} else {
		actor.Room.Dispatch(model.EvtCharacterUnequipsItem{
			Character: actor,
			Item:      c.Item,
		})
	}

	actor.TakeItem(c.Item)
}
