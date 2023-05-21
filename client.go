package updaterproxy

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
	send chan []byte

	Vmuuid   string `json:"vmuuid"`
	Sn       string `json:"sn"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.send <- message
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) Send(message []byte) {
	c.send <- message
}
