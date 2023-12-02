package main

import (
	"example.com/backend"
)

func main() {
	a := backend.App{Port: ":9003", DBType: "sqlite3", DBPath: "./practiceit.sqlite"}

	a.Initialize()

	a.RunServer()

}
