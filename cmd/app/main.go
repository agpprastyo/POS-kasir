package main

import (
	"POS-kasir/server"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	app := server.InitApp()
	defer server.Cleanup(app)

	// Start app
	go server.StartServer(app)

	// Wait for interrupt signal
	server.WaitForShutdown(app)
}
