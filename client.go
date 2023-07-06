package updaterproxy

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn        *websocket.Conn
	send        chan []byte
	UUID        string
	Vmuuid      string `json:"vmuuid"`
	Sn          string `json:"sn"`
	Hostname    string `json:"hostname"`
	Ip          string `json:"ip"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	Version     string `json:"version"`
	HostIp      string `json:"hostIp"`
	Proxy       *Proxy
	Nonce       string `json:"nonce"`
	ConnectTime int64  `json:"connectTime"`
}

func NewClient(conn *websocket.Conn, Uuid string, proxy *Proxy, hostip string, nonce string) *Client {
	return &Client{
		UUID:        Uuid,
		conn:        conn,
		send:        make(chan []byte),
		Proxy:       proxy,
		HostIp:      hostip,
		Nonce:       nonce,
		ConnectTime: time.Now().Unix(),
	}
}

func (c *Client) Start() {
	go c.readPump()
	go c.writePump()
}

func (c *Client) Disconnect() {

	if time.Now().Unix()-c.ConnectTime > 5 {
		c.Proxy.SendClientOfflineToServer(c.UUID)
	}

	c.conn.Close()
	c.Proxy.cm.Remove(c)
}

func (c *Client) readPump() {
	defer func() {
		c.Disconnect()
	}()
	for {
		log.Println("readPump: uuid:", c.UUID, "nonce:", c.Nonce)
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err, "uuid:", c.UUID)
			return
		}
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("json unmarshal err:", err)
			continue
		}
		msg.ClientIP = c.HostIp
		c.Proxy.SendToServer(&msg)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Disconnect()

	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("writePump: !ok")
				return
			}
			log.Println("writePump:", string(message), "uuid:", c.UUID, "nonce:", c.Nonce)

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("writePump err:", err)
			}
		}
	}
}

func (c *Client) Send(message []byte) {
	c.send <- message
}
