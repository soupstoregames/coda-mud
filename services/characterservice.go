package services

import (
	"errors"
	"io"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/soupstore/coda-world/log"
	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
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
}

// NewCharacterService returns a pointer to a character service and sets the character controller.
func NewCharacterService(controller simulation.CharacterController) *CharacterService {
	return &CharacterService{controller}
}

// Subscribe is the handler for the bidrirectional GRPC stream of commands and events.
func (s *CharacterService) Subscribe(stream Character_SubscribeServer) error {
	var err error

	log.Logger().Info("Player connected")

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

	commands := make(chan *CommandMessage)
	go s.listenForCommands(stream, commands)
	err = s.loop(stream, events, commands, characterID)
	if err != nil && err != errConnectionEnded {
		log.Logger().Error(err.Error())
	}

	log.Logger().Info("Player disconnected")

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

func (s *CharacterService) loop(stream Character_SubscribeServer, events <-chan interface{}, commands <-chan *CommandMessage, characterID model.CharacterID) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case command, more := <-commands:
			if !more {
				return nil
			}
			err := s.handleCommand(characterID, command)
			if err != nil {
				return err
			}
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

func (s *CharacterService) listenForCommands(stream Character_SubscribeServer, commands chan *CommandMessage) error {
	for {
		command, err := stream.Recv()

		// client disconnected
		if err == io.EOF {
			close(commands)
			return errConnectionEnded
		}

		grpcStatus, ok := status.FromError(err)
		if ok {
			switch grpcStatus.Code() {
			case codes.Canceled:
				close(commands)
				return errConnectionEnded
			}
		}

		// unknown error
		if err != nil {
			close(commands)
			return err
		}

		commands <- command
	}
}

func (s *CharacterService) handleCommand(characterID model.CharacterID, cmd *CommandMessage) error {
	// TODO: address the concurrency issues with this approach
	switch cmd.Type {
	case CommandType_CmdLook:
		s.controller.Look(characterID)

	case CommandType_CmdSay:
		var msg SayCommand
		err := proto.Unmarshal(cmd.Payload, &msg)
		if err != nil {
			return err
		}
		s.controller.Say(characterID, msg.Content)

	case CommandType_CmdNorth:
		err := s.controller.Move(characterID, model.North)
		if err != nil {
			return err
		}

	case CommandType_CmdNorthEast:
		err := s.controller.Move(characterID, model.NorthEast)
		if err != nil {
			return err
		}

	case CommandType_CmdEast:
		err := s.controller.Move(characterID, model.East)
		if err != nil {
			return err
		}

	case CommandType_CmdSouthEast:
		err := s.controller.Move(characterID, model.SouthEast)
		if err != nil {
			return err
		}

	case CommandType_CmdSouth:
		err := s.controller.Move(characterID, model.South)
		if err != nil {
			return err
		}

	case CommandType_CmdSouthWest:
		err := s.controller.Move(characterID, model.SouthWest)
		if err != nil {
			return err
		}

	case CommandType_CmdWest:
		err := s.controller.Move(characterID, model.West)
		if err != nil {
			return err
		}

	case CommandType_CmdNorthWest:
		err := s.controller.Move(characterID, model.NorthWest)
		if err != nil {
			return err
		}

	case CommandType_CmdTake:
		var msg TakeCommand
		err := proto.Unmarshal(cmd.Payload, &msg)
		if err != nil {
			return err
		}
		err = s.controller.TakeItem(characterID, msg.Alias)
		if err != nil {
			return err
		}

	case CommandType_CmdDrop:
		var msg DropCommand
		err := proto.Unmarshal(cmd.Payload, &msg)
		if err != nil {
			return err
		}
		err = s.controller.DropItem(characterID, msg.Alias)
		if err != nil {
			return err
		}
	}

	return nil
}
