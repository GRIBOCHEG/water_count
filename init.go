package main

var initFunctions = []func() error{
	initConfig,
	initDB,
	initServer,
}

func panicIfFailed(callback func() error) {
	err := callback()
	if err != nil {
		panic(err)
	}
}

func init() {
	for _, initFunc := range initFunctions {
		panicIfFailed(initFunc)
	}
}
