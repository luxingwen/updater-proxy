package updaterproxy

import (
	"errors"
	"sync"
)

type ClientManager struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*Client),
	}
}

func (manager *ClientManager) Add(client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients[client.UUID] = client
}

func (manager *ClientManager) Remove(client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if _, ok := manager.clients[client.UUID]; ok {
		delete(manager.clients, client.UUID)
		close(client.send)
		client.conn.Close()
		client = nil
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

	client, ok := manager.clients[clientUuid]
	if ok {
		return client
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
