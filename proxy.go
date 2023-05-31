package updaterproxy

// import "encoding/json"

type Proxy struct {
	cm *ClientManager
	sm *ServerManager
}

// func (proxy *Proxy) SendToClient(message *Message) {
// 	client := proxy.cm.GetClient(message.To, message.To.Sn, message.To.Hostname, message.To.Ip)
// 	if client != nil {
// 		jsonMessage, err := json.Marshal(message)
// 		if err != nil {
// 			return
// 		}
// 		client.Send(jsonMessage)
// 	} else {

// 		resp := *&Response{
// 			Code: "404",
// 			Msg:  "Not Found",
// 		}

// 		jsonMessage, err := json.Marshal(resp)

// 		if err != nil {
// 			return
// 		}

// 		message.Data = jsonMessage
// 		proxy.SendToServer(message)
// 	}
// }

// func (proxy *Proxy) SendToServer(message *Message) {
// 	server := proxy.sm.GetServer()
// 	if server != nil {
// 		jsonMessage, err := json.Marshal(message)
// 		if err != nil {
// 			return
// 		}
// 		server.Send(jsonMessage)
// 	} else {

// 	}
// }
