package simulation

import (
	"github.com/soupstore/coda/simulation/model"
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
	originalRoom.Dispatch(model.EvtCharacterLeaves{
		Character: actor,
		Direction: c.Direction,
	})

	// tell people in the target room that a character has arrived
	actor.Room.Dispatch(model.EvtCharacterArrives{
		Character: actor,
		Direction: c.Direction.Opposite(),
	})

	// move actor to the new room
	actor.Room = newRoom
	newRoom.AddCharacter(actor)
	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})

	// actor.Room.OnEnter(actor)
}

func (s *Simulation) takeItem(actor *model.Character, c model.CommandTake) {
	ok := actor.TakeItem(c.Item)
	if !ok {
		actor.Dispatch(model.EvtNoSpaceToTakeItem{})
		return
	}

	actor.Room.Container.RemoveItem(c.Item.ID)
	actor.Room.Dispatch(model.EvtCharacterTakesItem{
		Character: actor,
		Item:      c.Item,
	})
}

func (s *Simulation) dropItem(actor *model.Character, c model.CommandDrop) {
	ok := actor.DropItem(c.Item)
	if !ok {
		actor.Dispatch(model.EvtItemNotHere{})
		return
	}

	actor.Room.Container.PutItem(c.Item)
	actor.Room.Dispatch(model.EvtCharacterDropsItem{
		Character: actor,
		Item:      c.Item,
	})
}

func (s *Simulation) equipItem(actor *model.Character, c model.CommandEquip) {
	actor.Rig.RemoveItemFromInventory(c.Item)

	oldItem, err := actor.Equip(c.Item)
	if err == model.ErrNotEquipable {
		// do something
		actor.TakeItem(c.Item)
	}

	// do we have an item that we replaced
	if oldItem != nil {
		// tell everyone that we took off an item
		actor.Room.Dispatch(model.EvtCharacterUnequipsItem{
			Character: actor,
			Item:      oldItem,
		})
	}

	// tell everyone that we put on an item
	actor.Room.Dispatch(model.EvtCharacterEquipsItem{
		Character: actor,
		Item:      c.Item,
	})

	ok := actor.TakeItem(oldItem)
	if !ok {
		// do something
	}
}

func (s *Simulation) unequipItem(actor *model.Character, c model.CommandUnequip) {
	if ok := actor.Rig.Unequip(c.Item); !ok {
		// DO SOMETHING
	}

	actor.Room.Dispatch(model.EvtCharacterUnequipsItem{
		Character: actor,
		Item:      c.Item,
	})

	ok := actor.TakeItem(c.Item)
	if !ok {
		// DO SOMETHING
	}
}
