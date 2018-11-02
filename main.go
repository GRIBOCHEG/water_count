package main

// DB - Глобальная переменная содержащая подключение к БД
var DB = connectDB()

func main() {
	DB := connectDB()
	createSchema(DB)
	startServer()
}
