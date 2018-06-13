package telnet

import (
	"fmt"
	"net"

	"github.com/satori/go.uuid"
	"github.com/soupstore/coda-world/config"
	"github.com/soupstore/coda-world/log"
	"github.com/soupstore/coda-world/simulation"
	"go.uber.org/zap"
)

const (
	charNULL byte = 0
	charLF        = 10
	charCR        = 13
	charWILL      = 251
	charDO        = 252
	charWONT      = 253
	charDONT      = 254
	charIAC       = 255
)

// Server listens for incoming telnet connections
type Server struct {
	Config *config.Config
	Addr   string
	sim    *simulation.Simulation
}

func NewServer(c *config.Config, addr string, sim *simulation.Simulation) *Server {
	return &Server{
		Addr:   addr,
		Config: c,
		sim:    sim,
	}
}

func (server *Server) ListenAndServe() error {
	addr := server.Addr
	if "" == addr {
		addr = ":23"
	}

	listener, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	return server.Serve(listener)
}

func (server *Server) Serve(listener net.Listener) error {
	defer listener.Close()
	log.Logger().Debug(fmt.Sprintf("Listening at %q.", listener.Addr()))

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go server.handle(conn)
	}
}

func (server *Server) handle(tcpConn net.Conn) {
	connectionID := uuid.NewV4().String()

	c := newTelnetConnection(tcpConn, server.Config, server.sim)
	c.ctx = WithConnectionID(c.ctx, connectionID)

	log.Logger().Info(
		fmt.Sprintf("New connection from %q.", tcpConn.RemoteAddr()),
		zap.String("connectionID", connectionID),
	)

	err := c.listen()
	if err != nil {
		log.Logger().Error(
			err.Error(),
			zap.String("connectionID", connectionID),
		)
	}

	log.Logger().Info(
		fmt.Sprintf("Connection closed from %q.", tcpConn.RemoteAddr()),
		zap.String("connectionID", connectionID),
	)
}
