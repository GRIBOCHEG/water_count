package main

func main() {
	Init()
	startServer()
	defer app.DB.Close()
}
