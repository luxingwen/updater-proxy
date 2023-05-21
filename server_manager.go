package updaterproxy

import "sync"

type ServerManager struct {
	servers map[*Server]bool
	mu      sync.Mutex
}

func NewServerManager() *ServerManager {
	return &ServerManager{
		servers: make(map[*Server]bool),
	}
}

func (manager *ServerManager) Add(server *Server) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.servers[server] = true
}

func (manager *ServerManager) Remove(server *Server) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if _, ok := manager.servers[server]; ok {
		delete(manager.servers, server)
	}
}

func (manager *ServerManager) GetServer() *Server {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for server := range manager.servers {
		if server.IsConnected {
			return server
		}
	}
	return nil
}
