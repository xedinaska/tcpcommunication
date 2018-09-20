package tcp

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

//NewServer returns new tcp.Server instance with initialized fields. It accepts host/port as input params
func NewServer(host string, port int) *Server {
	return &Server{
		host:     host,
		port:     port,
		stopChan: make(chan os.Signal),
		clients:  make(map[string]*Client, 0),
		messages: make([]string, 0),
	}
}

//Server represents TCP server & contains handle methods
type Server struct {
	host     string
	port     int
	clients  map[string]*Client
	stopChan chan os.Signal
	listener net.Listener
	messages []string
}

//Start will start TCP listener
func (s *Server) Start() (err error) {
	//accept TCP connections on provided port
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		log.Printf("[ERROR]host=%s,port=%d; unable to start listener: `%s`", s.host, s.port, err.Error())
		return err
	}

	//listen for os signals to (SIGTERM / SIGINT) for graceful shutdown
	signal.Notify(s.stopChan, syscall.SIGTERM)
	signal.Notify(s.stopChan, syscall.SIGINT)

	go s.shutdown() //gracefully handle server shutdown (close connections, free resources / etc)

	//infinite loop in separate goroutine to handle incoming connections
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Printf("[ERROR]failed to accept incoming connection: `%s`", err.Error())
			}
			c := s.addClient(conn)
			log.Printf("[DEBUG]incoming connection: %s, clients connected: %d", conn.RemoteAddr(), len(s.clients))

			go s.handle(c)
		}
	}()

	return nil
}

//handle used to handle client connection (listen messages, disconnect, etc..)
func (s *Server) handle(c *Client) {
	buf := make([]byte, 1024)
	r := bufio.NewReader(c.conn)

LOOP:
	for {
		n, err := r.Read(buf)
		msg := string(buf[:n])

		switch err {
		case nil:
			s.messages = append(s.messages, msg)
			log.Printf("[DEBUG][%s] received client message #%d: `%s`:", c.address, len(s.messages), msg)
		case io.EOF:
			log.Printf("[DEBUG]EOF, disconnecting client %s..", c.address)

			//close TCP conn / etc
			c.Disconnect()

			//remove client from clients slice
			delete(s.clients, c.id)

			log.Printf("[DEBUG]..done, clients connected: %d", len(s.clients))
			break LOOP
		default:
			c.Disconnect()
			log.Printf("[ERROR]failed to read input message: `%s`", err.Error())
			break
		}
	}
}

//addClient is used to generate init client by provided conn, generate unique id and add client to connected clients slice
func (s *Server) addClient(conn net.Conn) *Client {
	c := &Client{
		id:      hex.EncodeToString(md5.New().Sum([]byte(conn.RemoteAddr().String()))),
		conn:    conn,
		address: conn.RemoteAddr().String(),
	}

	s.clients[c.id] = c
	return c
}

//shutdown used to disconnect all connected clients after shutdown server signal && close TCP connection
func (s *Server) shutdown() {
	<-s.stopChan
	log.Printf("[INFO]received shutdown signal. Stopping %d clients & exit..", len(s.clients))

	for _, c := range s.clients {
		if err := c.Disconnect(); err != nil {
			log.Printf("[ERROR] failed to disconnect client %s: `%s`", c.address, err.Error())
		}
	}

	log.Printf("[INFO]..done. Exit")

	s.listener.Close()
	os.Exit(0)
}
