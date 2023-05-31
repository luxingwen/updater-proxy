package updaterproxy

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	UUID     string
	Vmuuid   string `json:"vmuuid"`
	Sn       string `json:"sn"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Proxy    *Proxy
}

func NewClient(conn *websocket.Conn, Uuid string, proxy *Proxy) *Client {
	return &Client{
		UUID:  Uuid,
		conn:  conn,
		send:  make(chan []byte),
		Proxy: proxy,
	}
}

func (c *Client) Start() {
	go c.readPump()
	go c.writePump()
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
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		}
		c.Proxy.SendToServer(&msg)
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
