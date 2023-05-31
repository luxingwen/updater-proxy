package main

import (
	proxy "updater-proxy"
)

func main() {
	server := proxy.NewServer("ws://127.0.0.1:8080/api/v1/ws/proxy-uuid1")
	server.Start()
	select {}
}
