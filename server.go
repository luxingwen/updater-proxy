package updaterproxy

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	conn        *websocket.Conn
	send        chan []byte
	Url         string
	IsConnected bool
}

func NewServer(url string) *Server {
	return &Server{
		Url: url,
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
}

func (s *Server) Disconnect() {
	s.conn.Close()
	s.IsConnected = false
}

func (s *Server) writePump() {
	defer func() {
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
		// 这里你可以添加处理消息的逻辑
	}
}

func (s *Server) Start() error {
	err := s.Connect()
	if err != nil {
		return err
	}
	go s.readPump()
	go s.writePump()
	return nil
}

func (s *Server) Send(message []byte) {
	s.send <- message
}
