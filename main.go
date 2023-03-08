package main

import (
	"github.com/kanatsanan6/go-test/api"
	"github.com/kanatsanan6/go-test/configs"
	"github.com/kanatsanan6/go-test/db"
)

func init() {
	configs.LoadEnv()
	db.ConnectDB()
}

func main() {
	server, err := api.NewServer()
	if err != nil {
		panic(err)
	}

	server.Start()
}
