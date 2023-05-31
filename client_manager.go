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

func (manager *ClientManager) Count() int {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	return len(manager.clients)
}

func (manager *ClientManager) GetClient(clientUuid string) *Client {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if clientUuid != "" {
		for client := range manager.clients {
			if client.UUID == clientUuid {
				return client
			}
		}
	}

	return nil
}

func (manager *ClientManager) SendToClient(clientUuid string, message []byte) (err error) {
	client := manager.GetClient(clientUuid)
	if client != nil {
		client.Send(message)
		return nil
	}
	return errors.New("Client not found")
}
