package main

import (
	"github.com/kanatsanan6/go-test/configs"
	"github.com/kanatsanan6/go-test/db"
	"github.com/kanatsanan6/go-test/routes"
)

func init() {
	configs.LoadEnv()
	db.ConnectDB()
}

func main() {
	r := routes.Router()

	r.Run()
}
