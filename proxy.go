package updaterproxy

import "encoding/json"

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
	client := proxy.cm.GetClient(message.To)
	if client != nil {
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			return
		}
		client.Send(jsonMessage)
	} else {
		message.Code = "404"
		message.Msg = "Not Found"
		message.Method = METHOD_RESPONSE
		proxy.SendToServer(message)
	}
}

func (proxy *Proxy) SendToServer(message *Message) {
	server := proxy.sm.GetServer()
	if server != nil {
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			return
		}
		server.Send(jsonMessage)
	} else {

	}
}
