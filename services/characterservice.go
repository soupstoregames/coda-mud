package services

import (
	"errors"
	"io"
	"strconv"

	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	// ErrMalformedCharacterID means that the client connected with a malformed or missing
	// character ID in the metadata.
	ErrMalformedCharacterID = errors.New("character ID in metadata was missing or malformed")

	// ErrUnknownEventType is returned when an unknown type is dispatched to a characters event queue
	ErrUnknownEventType = errors.New("unknown event type")

	errConnectionEnded = errors.New("connection ended")
)

// CharacterService is a GRPC service for controlling characters.
type CharacterService struct {
	controller simulation.CharacterController
	logger     *zap.Logger
}

// NewCharacterService returns a pointer to a character service and sets the character controller.
func NewCharacterService(controller simulation.CharacterController, logger *zap.Logger) *CharacterService {
	return &CharacterService{controller, logger}
}

// Subscribe is the handler for the bidrirectional GRPC stream of commands and events.
func (s *CharacterService) Subscribe(stream Character_SubscribeServer) error {
	var err error

	s.logger.Info("Player connected")

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
	quit := make(chan struct{})
	var g errgroup.Group
	g.Go(s.listenForCommands(stream, characterID, quit))
	g.Go(s.sendEvents(stream, events, quit))
	err = g.Wait()
	if err != nil && err != errConnectionEnded {
		s.logger.Error(err.Error())
	}

	s.logger.Info("Player disconnected")

	return err
}

// extractCharacterID gets the id of the character the stream client wants to control
// the id is stored in the metadata inside the stream context
func (s *CharacterService) extractCharacterID(stream Character_SubscribeServer) (model.CharacterID, error) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return 0, ErrMalformedCharacterID
	}

	if len(md["characterid"]) != 1 {
		return 0, ErrMalformedCharacterID
	}

	characterIDint, err := strconv.Atoi(md["characterid"][0])
	if err != nil {
		return 0, ErrMalformedCharacterID
	}

	return model.CharacterID(characterIDint), nil
}

func (s *CharacterService) listenForCommands(stream Character_SubscribeServer, id model.CharacterID, quit chan<- struct{}) func() error {
	return func() error {
		for {
			command, err := stream.Recv()

			// client disconnected
			if err == io.EOF {
				quit <- struct{}{}
				return errConnectionEnded
			}

			grpcStatus, ok := status.FromError(err)
			if ok {
				switch grpcStatus.Code() {
				case codes.Canceled:
					quit <- struct{}{}
					return errConnectionEnded
				}
			}

			// unknown error
			if err != nil {
				quit <- struct{}{}
				return err
			}

			switch command.GetType() {
			// send commands to the simulation here
			}
		}

	}
}

func (s *CharacterService) sendEvents(stream Character_SubscribeServer, events <-chan interface{}, quit <-chan struct{}) func() error {
	return func() error {
		for {
			select {
			case <-quit:
				return nil
			case event, more := <-events:
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
}
