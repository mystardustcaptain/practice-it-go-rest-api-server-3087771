package main

import (
	"example.com/backend"
)

func main() {
	a := backend.App{Port: ":9003"}

	a.Initialize()

	a.RunServer()

}
