package services

import (
	"errors"
	"io"
	"strconv"

	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

var (
	// ErrMalformedCharacterID means that the client connected with a malformed or missing
	// character ID in the metadata.
	ErrMalformedCharacterID = errors.New("character ID in metadata was missing or malformed")

	// ErrUnknownEventType is returned when an unknown type is dispatched to a characters event queue
	ErrUnknownEventType = errors.New("unknown event type")
)

// CharacterService is a GRPC service for controlling characters.
type CharacterService struct {
	controller simulation.CharacterController
}

// NewCharacterService returns a pointer to a character service and sets the character controller.
func NewCharacterService(controller simulation.CharacterController) *CharacterService {
	return &CharacterService{controller}
}

// Subscribe is the handler for the bidrirectional GRPC stream of commands and events.
func (s *CharacterService) Subscribe(stream Character_SubscribeServer) error {
	var err error

	// get characterID from metadata
	characterID, err := s.extractCharacterID(stream)
	if err != nil {
		return err
	}

	// wake up the character, and it put it back to sleep when the controller disconnects
	events, err := s.controller.WakeUpCharacter(characterID)
	if err != nil {
		return err
	}
	defer s.controller.SleepCharacter(characterID)

	// start listening for commands and events
	var g errgroup.Group
	g.Go(s.listenForCommands(stream, characterID))
	g.Go(s.sendEvents(stream, events))
	err = g.Wait()

	return err
}

// extractCharacterID gets the id of the character the stream client wants to control
// the id is stored in the metadata inside the stream context
func (s *CharacterService) extractCharacterID(stream Character_SubscribeServer) (model.CharacterID, error) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return 0, ErrMalformedCharacterID
	}

	if len(md["characterID"]) != 1 {
		return 0, ErrMalformedCharacterID
	}

	characterIDint, err := strconv.Atoi(md["characterID"][0])
	if err != nil {
		return 0, ErrMalformedCharacterID
	}

	return model.CharacterID(characterIDint), nil
}

func (s *CharacterService) listenForCommands(stream Character_SubscribeServer, id model.CharacterID) func() error {
	return func() error {
		for {
			command, err := stream.Recv()

			// client disconnected
			if err == io.EOF {
				return nil
			}

			// unknown error
			if err != nil {
				return err
			}

			switch command.GetType() {
			// send commands to the simulation here
			}
		}
	}
}

func (s *CharacterService) sendEvents(stream Character_SubscribeServer, events <-chan interface{}) func() error {
	return func() error {
		for {
			// read from the events channel
			event, more := <-events
			if !more {
				// TODO: log warning
				return nil
			}

			// parse event
			eventMessage, err := buildEventMessage(event)
			if err != nil {
				return err
			}
			if eventMessage == nil {
				continue
			}

			// post the event over the stream
			if err = stream.Send(eventMessage); err != nil {
				return err
			}
		}
	}
}
