package updaterproxy

import (
	"encoding/json"
	"log"
)

// import "encoding/json"

type Proxy struct {
	cm *ClientManager
	sm *ServerManager
}

func NewProxy(cm *ClientManager, sm *ServerManager) *Proxy {
	return &Proxy{
		cm: cm,
		sm: sm,
	}
}

func (proxy *Proxy) SendToClient(message *Message) {

	b, _ := json.Marshal(message)
	log.Println("SendToClient:", string(b))

	client := proxy.cm.GetClient(message.To)
	if client != nil {

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Println("SendToClient error:", err)
			return
		}
		client.Send(jsonMessage)
	} else {
		log.Println("SendToClient error: client not found")
		message.Code = "404"
		message.Msg = "Not Found"
		message.Method = METHOD_RESPONSE
		proxy.SendToServer(message)
	}
}

func (proxy *Proxy) SendToServer(message *Message) {
	//b, _ := json.Marshal(message)
	//log.Println("SendToServer:", string(b))
	server := proxy.sm.GetServer()
	if server != nil {
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Println("SendToServer error:", err)
			return
		}
		server.Send(jsonMessage)
	} else {
		log.Println("SendToServer error: server not found")
	}
}

func (proxy *Proxy) SendClientOfflineToServer(clientUuid string) {

	if clientUuid == "" {
		return
	}

	message := &Message{
		Method: METHOD_REQUEST,
		To:     "",
		From:   clientUuid,
		Type:   "v1/ClientOffline",
		Msg:    "Client offline",
		Code:   "200",
	}
	proxy.SendToServer(message)
}
