package main

func main() {
	startServer()
	defer app.DB.Close()
}
