package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
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

	serverAddress := proxy.GetConfig().Servers[0] + proxy.GetConfig().ProxyId

	log.Println("serverAddress:", serverAddress)

	server := proxy.NewServer(serverAddress, proxyServer)
	go server.Start()

	serverManager.Add(server)

	router := gin.Default()

	tagerUrls := make([]*url.URL, 0)

	for _, item := range proxy.GetConfig().PkgServers {
		u, err := url.Parse(item)
		if err != nil {
			log.Fatal(err)
		}
		tagerUrls = append(tagerUrls, u)
	}

	// 定义代理转发的路由
	router.GET("/api/v1/pkg/*path", func(c *gin.Context) {

		dstpath := "/api/v1/pkg" + getFilePath(c)

		// 随机选择一个代理
		targetURL := tagerUrls[rand.Intn(len(tagerUrls))]

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = dstpath
			req.Host = ""
			req.Header.Del("Host")
			req.Header.Del("X-Forwarded-For")
		}

		// 使用选定的代理进行转发
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/api/v1/ws/:uuid", func(c *gin.Context) {

		clientip := c.ClientIP()

		uid := c.Param("uuid")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			// Handle the error
			return
		}

		count := clientManager.Count()

		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", count)))

		client := proxy.NewClient(conn, uid, proxyServer, clientip)

		clientManager.Add(client)
		go client.Start()

	})

	log.Println("proxy server started on port:", proxy.GetConfig().Port)
	router.Run(proxy.GetConfig().Port)
}

func getFilePath(c *gin.Context) (r string) {
	dstpath := c.Param("path")
	if !strings.HasPrefix(dstpath, "/") {
		dstpath = "/" + dstpath
	}
	dstpath = filepath.Clean(dstpath)
	return dstpath
}
