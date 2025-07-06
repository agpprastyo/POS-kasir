package main

import (
	"POS-kasir/server"
)

func main() {
	app := server.InitApp()
	defer server.Cleanup(app)

	go server.StartServer(app)

	server.WaitForShutdown(app)
}
