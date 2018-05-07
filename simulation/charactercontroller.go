package simulation

import (
	"github.com/soupstore/coda-world/simulation/model"
	"go.uber.org/zap"
)

// CharacterController is an interface over the Simulation that exposes all actions a connected
// player will need to perform
type CharacterController interface {
	WakeUpCharacter(model.CharacterID) (<-chan interface{}, error)
	SleepCharacter(model.CharacterID) error
	Look(model.CharacterID) error
	Say(model.CharacterID, string) error
	Move(model.CharacterID, model.Direction) error
	TakeItem(id model.CharacterID, alias string) error
	DropItem(id model.CharacterID, alias string) error
	EquipItem(id model.CharacterID, itemID model.ItemID) error
}

// WakeUpCharacter make a character wake up.
// It sends a room description to the waking character.
// It sends a character waking event to the other characters in the room.
func (s *Simulation) WakeUpCharacter(id model.CharacterID) (<-chan interface{}, error) {
	actor, ok := s.characters[id]
	if !ok {
		return nil, ErrCharacterNotFound
	}

	if actor.Awake {
		return nil, ErrCharacterAwake
	}

	// wake character and send description
	actor.Events = make(chan interface{}, 1)
	actor.Awake = true
	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})

	// send character wakes up
	wakeUpEvent := model.EvtCharacterWakesUp{Character: actor}
	for _, c := range actor.Room.GetCharacters() {
		// ignore the character that woke up
		if c == actor {
			continue
		}

		c.Dispatch(wakeUpEvent)
	}

	zap.L().Debug("Character woke up")

	return actor.Events, nil
}

// SleepCharacter sets a character to sleeping.
// It sends a character sleeping event to all other characters in the room.
func (s *Simulation) SleepCharacter(id model.CharacterID) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	actor.Awake = false
	close(actor.Events)

	// send character sleeps
	sleepEvent := model.EvtCharacterFallsAsleep{Character: actor}
	for _, c := range actor.Room.GetCharacters() {
		c.Dispatch(sleepEvent)
	}

	zap.L().Debug("Character fell asleep")

	return nil
}

// Look gives the character a room description
func (s *Simulation) Look(id model.CharacterID) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})
	return nil
}

// Say allows a character to speak to everyone in the same room
func (s *Simulation) Say(id model.CharacterID, content string) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	speechEvent := model.EvtCharacterSpeaks{
		Character: actor,
		Content:   content,
	}

	for _, c := range actor.Room.GetCharacters() {
		c.Dispatch(speechEvent)
	}

	return nil
}

func (s *Simulation) Move(id model.CharacterID, direction model.Direction) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}
	originalRoom := actor.Room

	exit := actor.Room.Exits[direction]
	if exit == nil {
		actor.Dispatch(model.EvtNoExitInThatDirection{})
		return nil
	}

	newRoom, err := s.GetRoom(exit.WorldID, exit.RoomID)
	if err != nil {
		return err
	}

	// remove actor from current room
	originalRoom.RemoveCharacter(actor)

	// tell people in the room that the actor has left
	personLeftEvent := model.EvtCharacterLeaves{
		Character: actor,
		Direction: direction,
	}
	for _, c := range originalRoom.GetCharacters() {
		c.Dispatch(personLeftEvent)
	}

	// tell people in the target room that a character has arrived
	personArrivedEvent := model.EvtCharacterArrives{
		Character: actor,
		Direction: direction.Opposite(),
	}
	for _, c := range newRoom.GetCharacters() {
		c.Dispatch(personArrivedEvent)
	}

	// move actor to the new room
	actor.Room = newRoom
	newRoom.AddCharacter(actor)
	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})

	return nil
}

func (s *Simulation) TakeItem(id model.CharacterID, alias string) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	for itemID, item := range actor.Room.Container.Items {
		if item.KnownAs(alias) {
			actor.Backpack.Container.PutItem(item)
			actor.Room.Container.RemoveItem(itemID)
			event := model.EvtCharacterTakesItem{
				Character: actor,
				Item:      item,
			}
			for _, ch := range actor.Room.GetCharacters() {
				ch.Dispatch(event)
			}
			return nil
		}
	}

	actor.Dispatch(model.EvtItemNotHere{})

	return nil
}

func (s *Simulation) DropItem(id model.CharacterID, alias string) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	for itemID, item := range actor.Backpack.Container.Items {
		if item.KnownAs(alias) {
			actor.Room.Container.PutItem(item)
			actor.Backpack.Container.RemoveItem(itemID)
			event := model.EvtCharacterDropsItem{
				Character: actor,
				Item:      item,
			}
			for _, ch := range actor.Room.GetCharacters() {
				ch.Dispatch(event)
			}
			return nil
		}
	}

	actor.Dispatch(model.EvtItemNotHere{})

	return nil
}

func (s *Simulation) EquipItem(id model.CharacterID, itemID model.ItemID) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	item, ok := actor.Room.Container.Items[itemID]
	if !ok {
		return ErrItemNotFound
	}

	_, err = actor.Equip(item)
	if err == model.ErrNotEquipable {
		return ErrCannotEquipItem
	}
	actor.Room.Container.RemoveItem(item.ID)

	return nil
}

// this checks that the character exists in the simulation and that they are awake (connected to)
func (s *Simulation) findAwakeCharacter(id model.CharacterID) (*model.Character, error) {
	actor, ok := s.characters[id]
	if !ok {
		return nil, ErrCharacterNotFound
	}

	if !actor.Awake {
		return nil, ErrCharacterAsleep
	}

	return actor, nil
}
