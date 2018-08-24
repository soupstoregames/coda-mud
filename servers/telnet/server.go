package telnet

import (
	"fmt"
	"net"

	"github.com/satori/go.uuid"
	"github.com/soupstore/coda/config"
	"github.com/soupstore/coda/services"
	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/go-core/logging"
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
	Config       *config.Config
	Addr         string
	sim          *simulation.Simulation
	usersManager *services.UsersManager
}

// NewServer is a helper constructor for building a server.
func NewServer(c *config.Config, sim *simulation.Simulation, usersManager *services.UsersManager) *Server {
	return &Server{
		Addr:         fmt.Sprintf("%s:%s", c.Address, c.Port),
		Config:       c,
		sim:          sim,
		usersManager: usersManager,
	}
}

// ListenAndServe tells the server to start listening for telnet connections.
func (server *Server) ListenAndServe() error {
	addr := server.Addr
	if "" == addr {
		addr = ":23"
	}

	listener, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	return server.serve(listener)
}

func (server *Server) serve(listener net.Listener) error {
	defer listener.Close()
	logging.Debug(fmt.Sprintf("Listening at %q.", listener.Addr()))

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

	c := newTelnetConnection(tcpConn, server.Config, server.sim, server.usersManager)
	c.ctx = WithConnectionID(c.ctx, connectionID)

	logger := logging.BuildConnectionLogger(connectionID)

	logging.Info(fmt.Sprintf("New connection from %q.", tcpConn.RemoteAddr()))

	if err := c.listen(); err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("Connection closed from %q.", tcpConn.RemoteAddr()))
}
