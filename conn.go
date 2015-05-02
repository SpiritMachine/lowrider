// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package lowrider

import (
	"bufio"
	//"bytes"
	"errors"
	//"fmt"
	//"io"
	"net"
	//"strconv"
	"sync"
	"time"
)

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string { return string(err) }

// Conn represents a connection to a Redis server.
type Conn interface {
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value if the connection is broken. The returned
	// value is either the first non-nil value returned from the underlying
	// network connection or a protocol parsing error. Applications should
	// close broken connections.
	Err() error

	// Do sends a command to the server and returns the received reply.
	Do(commandName string, args ...interface{}) (reply interface{}, err error)

	// Send writes the command to the client's output buffer.
	// Send(commandName string, args ...interface{}) error

	// Flush flushes the output buffer to the Redis server.
	// Flush() error

	// Receive receives a single reply from the Redis server
	// Receive() (reply interface{}, err error)
}

// conn is the low-level implementation of Conn
type conn struct {

	// Shared
	mu      sync.Mutex
	pending int
	err     error
	conn    net.Conn

	// Read
	readTimeout time.Duration
	br          *bufio.Reader

	// Write
	writeTimeout time.Duration
	bw           *bufio.Writer

	// Scratch space for formatting argument length.
	// '*' or '$', length, "\r\n"
	lenScratch [32]byte

	// Scratch space for formatting integers and floats.
	numScratch [40]byte
}

// Dial connects to Infinispan at the given network and address.
func Dial(network, address string) (Conn, error) {
	dialer := Dialer{}
	return dialer.Dial(network, address)
}

// DialTimeout acts like Dial but takes timeouts for establishing the
// connection to the server, writing a command and reading a reply.
func DialTimeout(network, address string, connectTimeout, readTimeout, writeTimeout time.Duration) (Conn, error) {
	netDialer := net.Dialer{Timeout: connectTimeout}
	dialer := Dialer{
		NetDial:      netDialer.Dial,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	return dialer.Dial(network, address)
}

// A Dialer specifies options for connecting to Inifispan.
type Dialer struct {
	// NetDial specifies the dial function for creating TCP connections. If
	// NetDial is nil, then net.Dial is used.
	NetDial func(network, addr string) (net.Conn, error)

	// ReadTimeout specifies the timeout for reading a single command
	// reply. If ReadTimeout is zero, then no timeout is used.
	ReadTimeout time.Duration

	// WriteTimeout specifies the timeout for writing a single command.  If
	// WriteTimeout is zero, then no timeout is used.
	WriteTimeout time.Duration
}

// Dial connects to Infinispan at address on the named network.
func (d *Dialer) Dial(network, address string) (Conn, error) {
	dial := d.NetDial
	if dial == nil {
		dial = net.Dial
	}
	netConn, err := dial(network, address)
	if err != nil {
		return nil, err
	}
	return &conn{
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  d.ReadTimeout,
		writeTimeout: d.WriteTimeout,
	}, nil
}

// NewConn returns a new Lowrider connection for the given net connection.
func NewConn(netConn net.Conn, readTimeout, writeTimeout time.Duration) Conn {
	return &conn{
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (c *conn) Close() error {
	c.mu.Lock()
	err := c.err
	if c.err == nil {
		c.err = errors.New("Lowrider: closed")
		err = c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) fatal(err error) error {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		// Close connection to force errors on subsequent calls and to unblock other reader or writer.
		c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *conn) Do(cmd string, args ...interface{}) (interface{}, error) {
	c.mu.Lock()
	pending := c.pending
	c.pending = 0
	c.mu.Unlock()

	if cmd == "" && pending == 0 {
		return nil, nil
	}

	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	// TODO: Flesh this out.
	return nil, nil
}
