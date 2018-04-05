package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/soupstore/coda-world/simulation/model"
)

// buildEventMessage takes an event from a character's event stream and tries to convert it
// into a protobuf *EventMessage.
// if the event is an unknown type, it will not error but the *EventMessage will be nil.
func buildEventMessage(event interface{}) (*EventMessage, error) {
	var err error
	var eventMessage *EventMessage

	switch v := event.(type) {
	case model.EvtRoomDescription:
		if eventMessage, err = buildEventRoomDescription(v); err != nil {
			return nil, err
		}
	case model.EvtCharacterWakesUp:
		if eventMessage, err = buildEventCharacterWakesUp(v); err != nil {
			return nil, err
		}
	case model.EvtCharacterFallsAsleep:
		if eventMessage, err = buildEventCharacterSleeps(v); err != nil {
			return nil, err
		}
	case model.EvtCharacterSpeaks:
		if eventMessage, err = buildEventCharacterSpeaks(v); err != nil {
			return nil, err
		}
	default:
		// TODO: log warning
	}

	return eventMessage, nil
}

func buildEventRoomDescription(event model.EvtRoomDescription) (*EventMessage, error) {
	characters := []*CharacterDescription{}
	for _, ch := range event.Room.GetCharacters() {
		characters = append(characters, buildCharacterDesciption(ch))
	}

	payload, err := proto.Marshal(&RoomDescriptionEvent{
		Name:        event.Room.Name,
		Description: event.Room.Description,
		Characters:  characters,
	})
	if err != nil {
		return nil, err
	}

	return &EventMessage{
		Type:    EventType_EvtRoomDescription,
		Payload: payload,
	}, nil
}

func buildEventCharacterWakesUp(event model.EvtCharacterWakesUp) (*EventMessage, error) {
	payload, err := proto.Marshal(&CharacterWakesUpEvent{
		Character: buildCharacterDesciption(event.Character),
	})

	if err != nil {
		return nil, err
	}

	return &EventMessage{
		Type:    EventType_EvtCharacterWakesUp,
		Payload: payload,
	}, nil
}

func buildEventCharacterSleeps(event model.EvtCharacterFallsAsleep) (*EventMessage, error) {
	payload, err := proto.Marshal(&CharacterSleepsEvent{
		Character: buildCharacterDesciption(event.Character),
	})

	if err != nil {
		return nil, err
	}

	return &EventMessage{
		Type:    EventType_EvtCharacterSleeps,
		Payload: payload,
	}, nil
}

func buildEventCharacterSpeaks(event model.EvtCharacterSpeaks) (*EventMessage, error) {
	payload, err := proto.Marshal(&CharacterSpeaksEvent{
		Character: buildCharacterDesciption(event.Character),
		Content:   event.Content,
	})

	if err != nil {
		return nil, err
	}

	return &EventMessage{
		Type:    EventType_EvtCharacterSpeaks,
		Payload: payload,
	}, nil
}

func buildCharacterDesciption(character *model.Character) *CharacterDescription {
	return &CharacterDescription{
		Id:    int64(character.ID),
		Name:  character.Name,
		Awake: character.Awake,
	}
}
