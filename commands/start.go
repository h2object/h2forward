package commands

import (
	"os"
	"fmt"
	"path"
	"path/filepath"
	"time"
	"syscall"
	"github.com/codegangsta/cli"
	"github.com/h2object/pidfile"
	"github.com/h2object/h2object/log"
	"github.com/h2object/h2forward/app"
)

const (
	success        = "\t\t\t\t\t[  \033[32mOK\033[0m  ]" // Show colored "OK"
	failed         = "\t\t\t\t\t[\033[31mFAILED\033[0m]" // Show colored "FAILED"
)

func startCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	if workdir == "" {
		fmt.Println("unknown working directory, please use -w to provide.")
		os.Exit(1)
	}
	// daemon 
	daemon := ctx.GlobalBool("daemon")
	
	// options
	options := app.NewOptions(ctx.GlobalString("http"),ctx.GlobalString("api"))
	if err := options.Prepare(workdir); err != nil {
		fmt.Println("options prepare failed:", err)
		os.Exit(1)
	}

	refresh := ctx.GlobalString("refresh")
	options.SetRefreshDefault(refresh, time.Minute * 10)

	// configs
	configs, err := app.LoadCONFIG(path.Join(options.Root, "h2forward.conf"))
	if err != nil {
		configs = app.DefaultCONFIG()
		if err := configs.Save(path.Join(options.Root, "h2forward.conf")); err != nil {
			fmt.Println("h2forward.conf saving failed:", err)
			os.Exit(1)
		}
	}
	
	logger := log.NewH2OLogger()
	defer logger.Close()
	logger.SetConsole(true)
	
	configs.SetSection("logs")
	fenable := configs.BoolDefault("file.enable", false)
	fname := configs.StringDefault("file.name", "h2o.log")
	flevel := configs.StringDefault("file.level", "info")
	fsize := configs.IntDefault("file.rotate_max_size", 1024*1024*1024)
	fline := configs.IntDefault("file.rotate_max_line", 102400)
	fdaily := configs.BoolDefault("file.rotate_daily", true)
	fn := path.Join(options.LogsRoot, fname)
	if fenable == true {
		logger.SetFileLog(fn, flevel, fsize, fline, fdaily)	
	}	

	application := app.NewApplication(options, configs, logger)

	if err := application.Init(); err != nil {
		fmt.Println("[h2forward] init failed:", err)
		os.Exit(1)
	}

	start(application, daemon)
}

func stopCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	if workdir == "" {
		fmt.Println("unknown working directory, please use -w to provide.")
		os.Exit(1)
	}

	// directory
	directory, err := filepath.Abs(workdir)
	if err != nil {
		fmt.Println("workdir:", err)
		return
	}

	pid, err := pidfile.Load(path.Join(directory, "h2forward.pid"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := pid.Kill(); err != nil {
		fmt.Println(err.Error())
		return	
	}

	fmt.Println("[h2forward] stop", success)
	return	
}

func run(application *app.Application) {
	pid, err := pidfile.New(path.Join(application.Options.Root, "h2forward.pid"))
	if err != nil {

	}
	defer pid.Kill()

	exitc := make(chan int)
	signc := make(chan os.Signal, 1)

	go func(){
		for {
			sig := <- signc
			switch sig {
			case syscall.SIGHUP:
				application.Refresh()
				continue
			default:
				exitc <- 1
				break
			}
		}
	}()

	application.Main()
	<- exitc 
	application.Exit()
}

func start(application *app.Application, daemon bool) {
	if !daemon {
		run(application)
	} else {
		daemonize(application)
	}
}