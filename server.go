package updaterproxy

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	conn        *websocket.Conn
	send        chan []byte
	Url         string
	IsConnected bool
	Proxy       *Proxy
}

func NewServer(url string, proxy *Proxy) *Server {
	return &Server{
		Proxy: proxy,
		Url:   url,
		send:  make(chan []byte),
	}
}

func (s *Server) Connect() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(s.Url, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	s.conn = conn
	s.IsConnected = true
	log.Println("Connected to server:", s.Url)
	return nil
}

func (s *Server) Disconnect() {
	s.conn.Close()
	s.IsConnected = false
}

func (s *Server) writePump() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		s.Disconnect()
	}()
	for {
		select {
		case message, ok := <-s.send:
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			s.conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:

			mdata := make(map[string]interface{})
			mdata["time"] = time.Now().Unix()
			bdata, err := json.Marshal(mdata)
			if err != nil {
				log.Println("json marshal error:", err)
				continue
			}

			msg := Message{
				Type: "ProxyHeartBeat",
				Data: json.RawMessage(bdata),
			}

			b, err := json.Marshal(msg)
			if err != nil {
				log.Println("json marshal error:", err)
				continue
			}
			s.conn.WriteMessage(websocket.TextMessage, b)
		}
	}
}

func (s *Server) readPump() {
	defer func() {
		s.Disconnect()
	}()
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			log.Println("reconnecting...")
			time.Sleep(5 * time.Second)
			s.Connect()
			continue
		}

		log.Println("recv: ", string(message))

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		}

		s.Proxy.SendToClient(&msg)
		// 这里你可以添加处理消息的逻辑
	}
}

func (s *Server) Start() error {

	for {
		err := s.Connect()
		if err != nil {
			log.Println("connect to server error:", err)
			log.Println("reconnecting...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	go s.readPump()
	go s.writePump()
	return nil
}

func (s *Server) Send(message []byte) {
	s.send <- message
}
