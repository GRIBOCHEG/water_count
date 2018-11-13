package main

func init() {
	app.Slice = []byte(config.Server.Secret)
	initDB()
	initServer()
}
