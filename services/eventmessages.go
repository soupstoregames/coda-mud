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
	case model.EvtCharacterArrives:
		if eventMessage, err = buildEventCharacterArrives(v); err != nil {
			return nil, err
		}
	case model.EvtCharacterLeaves:
		if eventMessage, err = buildEventCharacterLeaves(v); err != nil {
			return nil, err
		}
	case model.EvtNoExitInThatDirection:
		if eventMessage, err = buildEventNoExitInThatDirection(); err != nil {
			return nil, err
		}
	default:
		// TODO: log warning
	}

	return eventMessage, nil
}

func buildEventRoomDescription(event model.EvtRoomDescription) (*EventMessage, error) {
	// build present characters
	characters := []*CharacterDescription{}
	for _, ch := range event.Room.GetCharacters() {
		characters = append(characters, buildCharacterDesciption(ch))
	}

	// build items
	items := []*ItemDescription{}
	for _, v := range event.Room.Container.Items {
		items = append(items, &ItemDescription{Name: v.GetName()})
	}

	// build exits
	exits := []*ExitDescription{}
	keys := make([]model.Direction, 8)
	for k := range event.Room.Exits {
		keys = append(keys, k)
	}
	for _, k := range keys {
		exit := event.Room.Exits[k]
		if exit == nil {
			continue
		}
		exits = append(exits, &ExitDescription{
			Direction: mapDirection(k),
			Name:      exit.Name,
		})
	}

	payload, err := proto.Marshal(&RoomDescriptionEvent{
		Name:        event.Room.Name,
		Description: event.Room.Description,
		Characters:  characters,
		Items:       items,
		Exits:       exits,
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

func buildEventCharacterArrives(event model.EvtCharacterArrives) (*EventMessage, error) {
	payload, err := proto.Marshal(&CharacterArrivesEvent{
		Character: buildCharacterDesciption(event.Character),
		Direction: mapDirection(event.Direction),
	})

	if err != nil {
		return nil, err
	}

	return &EventMessage{
		Type:    EventType_EvtCharacterArrives,
		Payload: payload,
	}, nil
}

func buildEventCharacterLeaves(event model.EvtCharacterLeaves) (*EventMessage, error) {
	payload, err := proto.Marshal(&CharacterLeavesEvent{
		Character: buildCharacterDesciption(event.Character),
		Direction: mapDirection(event.Direction),
	})

	if err != nil {
		return nil, err
	}

	return &EventMessage{
		Type:    EventType_EvtCharacterLeaves,
		Payload: payload,
	}, nil
}

func buildEventNoExitInThatDirection() (*EventMessage, error) {
	return &EventMessage{
		Type: EventType_EvtNoExitInThatDirection,
	}, nil
}

func buildCharacterDesciption(character *model.Character) *CharacterDescription {
	return &CharacterDescription{
		Id:    int64(character.ID),
		Name:  character.Name,
		Awake: character.Awake,
	}
}

func mapDirection(d model.Direction) Direction {
	switch d {
	case model.North:
		return Direction_North
	case model.NorthEast:
		return Direction_NorthEast
	case model.East:
		return Direction_East
	case model.SouthEast:
		return Direction_SouthEast
	case model.South:
		return Direction_South
	case model.SouthWest:
		return Direction_SouthWest
	case model.West:
		return Direction_West
	case model.NorthWest:
		return Direction_NorthWest
	}

	return 0
}
