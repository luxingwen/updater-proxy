package main

import (
	"fmt"
	"net/http"
	proxy "updater-proxy"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by not checking the origin
		return true
	},
}

func main() {

	serverManager := proxy.NewServerManager()

	clientManager := proxy.NewClientManager()

	proxyServer := proxy.NewProxy(clientManager, serverManager)

	server := proxy.NewServer("ws://127.0.0.1:8080/api/v1/ws/proxy-uuid1", proxyServer)
	server.Start()

	serverManager.Add(server)

	router := gin.Default()

	router.GET("/api/v1/ws/:uuid", func(c *gin.Context) {
		uid := c.Param("uuid")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			// Handle the error
			return
		}

		count := clientManager.Count()

		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", count)))

		client := proxy.NewClient(conn, uid, proxyServer)
		go client.Start()

	})

	router.Run(":8081")
}
