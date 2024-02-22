package telnet

import (
	"fmt"
	"github.com/google/uuid"
	"net"

	"github.com/soupstoregames/coda-mud/config"
	"github.com/soupstoregames/coda-mud/services"
	"github.com/soupstoregames/coda-mud/simulation"
	"github.com/soupstoregames/go-core/logging"
)

const (
	charNULL     byte = 0
	charECHO          = 1
	charSGA           = 3
	charLF            = 10
	charCR            = 13
	charESC           = 27
	charNAWS          = 31
	charLINEMODE      = 34
	charSE            = 240
	charSB            = 250
	charWILL          = 251
	charWONT          = 252
	charDO            = 253
	charDONT          = 254
	charIAC           = 255
)

var byteToIAC = map[byte]string{
	charNAWS:     "NAWS",
	charLINEMODE: "LINEMODE",
	charSE:       "SE",
	charSB:       "SB",
	charWILL:     "WILL",
	charWONT:     "WONT",
	charDO:       "DO",
	charDONT:     "DONT",
	charIAC:      "IAC",
}

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
	connectionID := uuid.NewString()

	c := newTelnetConnection(tcpConn, server.Config, server.sim, server.usersManager)
	c.ctx = WithConnectionID(c.ctx, connectionID)

	logger := logging.BuildConnectionLogger(connectionID)

	logging.Info(fmt.Sprintf("New connection from %q.", tcpConn.RemoteAddr()))

	c.conn.Write([]byte{charIAC, charWILL, charECHO})
	c.conn.Write([]byte{charIAC, charWILL, charSGA})
	c.conn.Write([]byte{charIAC, charWONT, charLINEMODE})
	c.conn.Write([]byte{charIAC, charDO, charNAWS})

	c.listen()

	logger.Info(fmt.Sprintf("Connection closed from %q.", tcpConn.RemoteAddr()))
}
