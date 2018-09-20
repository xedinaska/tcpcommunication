package tcp

import "net"

//msgSTOP will be sent to client in case of server stop
const msgSTOP = "STOP"

//Client describes connected client
type Client struct {
	conn    net.Conn
	id      string
	address string
}

//Disconnect is used to close connection / free resources / etc
func (c *Client) Disconnect() error {
	c.conn.Write([]byte(msgSTOP))
	return c.conn.Close()
}
