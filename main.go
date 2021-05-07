package main

import (
	"rad-go/app"
)

func main() {
	radApp := &app.App{}
	radApp.InitializeAndRun(":8085")
}
