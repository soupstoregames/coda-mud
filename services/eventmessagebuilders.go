package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/soupstore/coda-world/simulation/model"
	"github.com/soupstore/mud-experiment/world-server/services/character"
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
	default:
		// TODO: log warning
	}

	return eventMessage, nil
}

func buildEventRoomDescription(event model.EvtRoomDescription) (*EventMessage, error) {
	payload, err := proto.Marshal(&character.RoomDescriptionEvent{
		Name:        event.Room.Name,
		Description: event.Room.Description,
	})
	if err != nil {
		return nil, err
	}

	return &EventMessage{
		Type:    EventType_EvtRoomDescription,
		Payload: payload,
	}, nil
}
