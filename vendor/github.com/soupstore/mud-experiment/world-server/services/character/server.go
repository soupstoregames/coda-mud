package character

import (
	"io"

	"golang.org/x/sync/errgroup"
)

const bufferSize = 100

type Server struct {
	Commands chan *CommandMessage
	Events   chan *EventMessage
}

func NewServer() *Server {
	s := new(Server)
	s.Commands = make(chan *CommandMessage, bufferSize)
	s.Events = make(chan *EventMessage, bufferSize)
	return s
}

func (s *Server) Subscribe(stream Character_SubscribeServer) error {
	var g errgroup.Group

	g.Go(func() error {
		for {
			// listen for commands from the gateways
			command, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			s.Commands <- command
		}
	})

	g.Go(func() error {
		for {
			// read from the events channel
			event, more := <-s.Events
			if !more {
				return nil
			}

			// post the event over the stream
			err := stream.Send(event)
			if err != nil {
				return err
			}
		}
	})

	return g.Wait()
}
