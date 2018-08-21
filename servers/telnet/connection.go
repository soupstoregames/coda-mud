package telnet

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/soupstore/coda/config"
	"github.com/soupstore/coda/services"
	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/go-core/logging"
)

// connection is a telnet connection to the MUD
type connection struct {
	config       *config.Config
	sim          *simulation.Simulation
	usersManager *services.UsersManager

	input         bytes.Buffer
	conn          net.Conn
	ctx           context.Context
	state         state
	closed        bool
	stopHeartbeat chan struct{}
}

func newTelnetConnection(c net.Conn, conf *config.Config, sim *simulation.Simulation, usersManager *services.UsersManager) *connection {
	conn := &connection{
		config:       conf,
		conn:         c,
		sim:          sim,
		usersManager: usersManager,
		ctx:          context.Background(),
	}

	conn.createHeartbeat(time.Minute)

	// create a login state and initiate it with the connection
	conn.state = &stateLogin{
		conn:   conn,
		config: conf,
	}

	// run the tasks at the start of the login state
	if err := conn.state.onEnter(); err != nil {
		logging.Error(err.Error())
	}

	return conn
}

func (c *connection) close() {
	if c.closed {
		return
	}

	logging.Debug("Closing telnet connection")

	// set flag to avoid close being called twice
	// without this you can try to close a closed channel
	c.closed = true

	// stop the heartbeat
	c.stopHeartbeat <- struct{}{}

	// close the client connection
	if err := c.conn.Close(); err != nil {
		logging.Error(err.Error())
	}

	// // run the clean up tasks on the current state and then clear it
	if err := c.state.onExit(); err != nil {
		logging.Error(err.Error())
	}
	c.state = nil
}

func (c *connection) handleInput() error {
	// empty input so do nothing
	if c.input.Len() == 0 {
		return nil
	}

	input := []byte{}
	for {
		b, err := c.input.ReadByte()
		if err == io.EOF {
			break
		}

		// the IAC character is always followed by one or two more characters
		if b == charIAC {
			d, err := c.input.ReadByte()
			if err == io.EOF {
				break
			}

			// these characters all mean that there'll be one more before the next content to control code
			if d == charWILL || d == charWONT || d == charDO || d == charDONT {
				_, err := c.input.ReadByte()
				if err == io.EOF {
					break
				}
			}
		} else {
			// add this character to the cleansed buffer
			input = append(input, b)
		}
	}

	// let the state handle the input
	if err := c.state.handleInput(string(input)); err != nil {
		return err
	}

	// clear the input buffer
	c.input.Reset()
	return nil
}

// listen is a for loop that reads bytes from the telnet connection
// and acts accordingly
func (c *connection) listen() error {
	// parse input
	var buffer [1]byte
	p := buffer[:]

	for {
		// read one byte from the connection into the buffer
		n, err := c.conn.Read(p)

		// not sure if this happens, but I feel like EOF could be -1?
		if n < 0 {
			fmt.Println("disconnected?")
		}

		// handle a couple of known errors but return any unknown ones
		if nil != err {
			switch err.Error() {
			case "EOF":
				// graceful disconnection
				c.close()
				return nil
			case "Corrupted":
				// forced disconnection
				c.close()
				return nil
			default:
				neterr, ok := err.(net.Error)
				if ok && neterr.Timeout() {
					logging.Debug("Connection timed out")
					c.close()
					return nil
				}

				if ok && neterr.Temporary() {
					logging.Debug("Temporary Net Error ???")
					c.close()
					return nil
				}

				c.close()
			}
		}

		// NUL CR and LF are all valid terminators according to telnet spec, I think
		switch buffer[0] {
		case charNULL:
			fallthrough
		case charCR:
			fallthrough
		case charLF:
			err := c.handleInput()
			if err != nil {
				// logging.Debug("error handlnginput")
				c.close()
				return nil
			}
		default:
			c.input.Write(p)
		}
	}
}

// used to load a new state for the connection
// takes care of calling the onExit and onEnter of the states
func (c *connection) loadState(s state) {
	c.state.onExit()
	c.state = s
	c.state.onEnter()
}

// sets up a timer to send a telnet NOOP to the client, in an attempt to keep tcp connections alive
func (c *connection) createHeartbeat(d time.Duration) {
	ticker := time.NewTicker(d)
	c.stopHeartbeat = make(chan struct{})
	go func() {
		for {
			select {
			// has the ticker ticked?
			case <-ticker.C:
				c.write([]byte{255, 241})
			// have we got the signal to stop?
			case <-c.stopHeartbeat:
				ticker.Stop()
				close(c.stopHeartbeat)
				return
			}
		}
	}()
}

// these are a bunch of write methods, pretty boring

func (c *connection) write(b []byte) {
	_, err := c.conn.Write(b)
	if err != nil {
		logging.Error(err.Error())
	}
}

func (c *connection) writeln(b []byte) {
	_, err := c.conn.Write(b)
	if err != nil {
		logging.Error(err.Error())
	}

	_, err = c.conn.Write([]byte{charLF, charCR})
	if err != nil {
		logging.Error(err.Error())
	}
}

func (c *connection) writeString(str ...string) {
	s := strings.Join(str, " ")

	_, err := c.conn.Write([]byte(s))
	if err != nil {
		logging.Error(err.Error())
	}
}

func (c *connection) writelnString(str ...string) {
	s := strings.Join(str, " ")
	_, err := c.conn.Write([]byte(s))
	if err != nil {
		logging.Error(err.Error())
	}
	_, err = c.conn.Write([]byte{charLF, charCR})
	if err != nil {
		logging.Error(err.Error())
	}
}

func (c *connection) writePrompt() {
	_, err := c.conn.Write([]byte{charLF, charCR})
	if err != nil {
		logging.Error(err.Error())
	}
	_, err = c.conn.Write([]byte("> "))
	if err != nil {
		logging.Error(err.Error())
	}
}
