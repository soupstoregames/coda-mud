package simulation

import (
	"github.com/soupstore/coda-world/simulation/model"
)

type Simulation struct {
	nextRoomID      model.RoomID
	nextCharacterID model.CharacterID
	spawnRoom       *model.Room
	rooms           map[model.RoomID]*model.Room
	items           map[model.ItemID]*model.Item
	characters      map[model.CharacterID]*model.Character
}

func NewSimulation() *Simulation {
	return &Simulation{
		nextRoomID:      0,
		nextCharacterID: 0,
		spawnRoom:       nil,
		rooms:           make(map[model.RoomID]*model.Room),
		items:           make(map[model.ItemID]*model.Item),
		characters:      make(map[model.CharacterID]*model.Character),
	}
}

func (s *Simulation) getNextRoomID() model.RoomID {
	roomID := s.nextRoomID
	s.nextRoomID = s.nextRoomID + 1
	return roomID
}

func (s *Simulation) getNextCharacterID() model.CharacterID {
	id := s.nextCharacterID
	s.nextCharacterID = s.nextCharacterID + 1
	return id
}

// MakeRoom creates a new room at the next available ID
func (s *Simulation) MakeRoom(name, description string) model.RoomID {
	roomID := s.getNextRoomID()
	room := model.NewRoom(roomID, name, description)
	s.rooms[roomID] = room
	return roomID
}

// SetSpawnRoom sets the room that all new characters will start in
func (s *Simulation) SetSpawnRoom(id model.RoomID) {
	spawnRoom, ok := s.rooms[id]
	if !ok {
		// error
	}
	s.spawnRoom = spawnRoom
}

// MakeCharacter creates a new character at the next available ID
// It returns the new character's ID
func (s *Simulation) MakeCharacter(name string) model.CharacterID {
	// get and increment character ID
	characterID := s.getNextCharacterID()

	// create new character and add to sim
	character := model.NewCharacter(characterID, name, s.spawnRoom)
	s.characters[characterID] = character

	// add character to room
	s.spawnRoom.AddCharacter(character)

	return characterID
}

// WakeUpCharacter make a character wake up.
// It sends a room description to the waking character.
// It sends a character waking event to the other characters in the room.
func (s *Simulation) WakeUpCharacter(id model.CharacterID) <-chan interface{} {
	actor, ok := s.characters[id]
	if !ok {
		// error
	}

	// wake character and send description
	actor.Awake = true
	actor.Dispatch(model.EvtRoomDescription{actor.Room})

	// send character wakes up
	wakeUpEvent := model.EvtCharacterWakesUp{Character: actor}
	for _, c := range actor.Room.GetCharacters() {
		// ignore the character that woke up
		if c == actor {
			continue
		}
		// dont both alerting sleeping players of events
		if !c.Awake {
			continue
		}

		c.Dispatch(wakeUpEvent)
	}

	return actor.Events
}

// SleepCharacter sets a character to sleeping.
// It sends a character sleeping event to all other characters in the room.
func (s *Simulation) SleepCharacter(id model.CharacterID) {
	actor, ok := s.characters[id]
	if !ok {
		// error
	}

	actor.Awake = false

	// send character sleeps
	sleepEvent := model.EvtCharacterFallsAsleep{Character: actor}
	for _, c := range actor.Room.GetCharacters() {
		// ignore sleeping characters
		if !c.Awake {
			continue
		}

		c.Dispatch(sleepEvent)
	}
}
