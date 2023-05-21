package updaterproxy

import (
	"errors"
	"sync"
)

type ClientManager struct {
	clients map[*Client]bool
	mu      sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*Client]bool),
	}
}

func (manager *ClientManager) Add(client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients[client] = true
}

func (manager *ClientManager) Remove(client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if _, ok := manager.clients[client]; ok {
		delete(manager.clients, client)
		close(client.send)
	}
}

func (manager *ClientManager) GetClient(vmuuid, sn, hostname, ip string) *Client {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if vmuuid == "" && sn == "" && hostname == "" && ip == "" {
		for client := range manager.clients {
			if (client.Vmuuid == vmuuid && vmuuid != "") && (client.Sn == sn && sn != "") && (client.Hostname == hostname && hostname != "") && (client.Ip == ip && ip != "") {
				return client
			}
		}
	}

	if vmuuid != "" {
		for client := range manager.clients {
			if client.Vmuuid == vmuuid {
				return client
			}
		}
	}

	if sn != "" {
		for client := range manager.clients {
			if client.Sn == sn {
				return client
			}
		}
	}

	if hostname != "" {
		for client := range manager.clients {
			if client.Hostname == hostname {
				return client
			}
		}
	}

	if ip != "" {
		for client := range manager.clients {
			if client.Ip == ip {
				return client
			}
		}

	}
	return nil
}

func (manager *ClientManager) SendToClient(vmuuid, sn, hostname, ip string, message []byte) (err error) {
	client := manager.GetClient(vmuuid, sn, hostname, ip)
	if client != nil {
		client.Send(message)
		return nil
	}
	return errors.New("Client not found")
}
