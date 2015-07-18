package commands 

import (
	"github.com/codegangsta/cli"
)

const version = "0.0.1"
const author = ""
const support = "liujianping@h2object.io"

func App() *cli.App {
	app := cli.NewApp()

	//! app settings
	app.Name = "h2forward"
	app.Usage = "http reverse proxy with http api to control the virtual hosts to forward"
	app.Version = version
	app.Author = author
	app.Email = support

	//! app flags
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "http, l",
			Value: "0.0.0.0:80",
			Usage: "proxy server http service address",
		},
		cli.StringFlag{
			Name: "api, a",
			Value: "127.0.0.1:9000",
			Usage: "proxy server api service address",
		},
		cli.StringFlag{
			Name: "workdir, w",
			Value: "",
			Usage: "working directory",
		},
		cli.StringFlag{
			Name: "refresh, r",
			Value: "10m",
			Usage: "refresh interval",
		},
		cli.BoolFlag{
			Name: "daemon, d",
			Usage: "run @ daemon mode",
		},
	}

	//! app commands
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start proxy & api service",
			Action: func(ctx *cli.Context) {
						startCommand(ctx)	
					},
		},
		{
			Name:  "stop",
			Usage: "start proxy & api service",
			Action: func(ctx *cli.Context) {
						stopCommand(ctx)	
					},
		},
	}

	return app
}
