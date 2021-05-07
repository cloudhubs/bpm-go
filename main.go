package main

import (
	"rad-go/app"
)

func main() {
	app := &app.App{}
	app.Initialize()
	app.Run(":8085")
}
