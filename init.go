package main

import (
	"encoding/binary"
	"time"
)

func init() {
	stamp := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(stamp, time.Now().UnixNano())
	app.Slice = append(stamp, []byte(config.Server.Secret)...)
	initDB()
	initServer()
}
