package commands

import (
	"fmt"
	"github.com/h2object/h2forward/app"
)

func daemonize(application *app.Application) {
	fmt.Println("[h2forward] warn: ", "windows not support daemon mode")
	run(application)	
}