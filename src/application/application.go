package main

import (
	"example.com/backend"
)

func main() {
	a := backend.App{
		Port:   ":9003",
		DBType: "sqlite3",
		DBPath: "./practiceit.db",
	}

	a.Initialize()
	a.GetData()

	a.Run()

}
