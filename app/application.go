package app

import (
	"net"
	"sync"
	"time"
	"github.com/h2object/h2object/log"
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/httpext"
	"github.com/h2object/oxy/forward"
	"github.com/h2object/oxy/reverse"
)

type Application struct{
	sync.RWMutex
	log.Logger

	// option & config
	Options		 *Options
	Configs		 *CONFIG
	
	// http
	httpAddr 	 *net.TCPAddr
	httpListener net.Listener

	// api
	apiAddr 	 *net.TCPAddr
	apiListener  net.Listener
	
	// virtual host
	vhosts  	 *VirtualHost

	// forward
	forwarder 	 *forward.Forwarder

	// reverse router
	reverser     *reverse.ReverseRouter

	// api controller
	api   		 *APIController

	// background workers
	background    util.Background
	exitc   	  chan int
}

func NewApplication(options *Options, configs *CONFIG, logger log.Logger) *Application {
	return &Application{
		Logger: logger,
		Options: options,
		Configs: configs,
		exitc: make(chan int),
	}
}


func (app *Application) Init() error {
	httpAddr, err := net.ResolveTCPAddr("tcp", app.Options.HTTPAddress)
	if err != nil {
		return err
	}
	app.httpAddr = httpAddr

	httpListener, err := net.Listen("tcp", app.httpAddr.String())
	if  err != nil {
		return err
	}
	app.httpListener = httpListener

	apiAddr, err := net.ResolveTCPAddr("tcp", app.Options.APIAddress)
	if err != nil {
		return err
	}
	app.apiAddr = apiAddr

	apiListener, err := net.Listen("tcp", app.apiAddr.String())
	if  err != nil {
		return err
	}
	app.apiListener = apiListener

	vhosts := NewVirtualHost(app.Options.Root, app.Logger)
	if err := vhosts.Load(); err != nil {
		return err
	}
	app.vhosts = vhosts

	forwarder, err := forward.New()
	if err != nil {
		return err
	}
	app.forwarder = forwarder

	reverser, err := reverse.New(app.forwarder, reverse.Route(app.vhosts))
	if err != nil {
		return err
	}
	app.reverser = reverser

	app.api = NewAPIController(app.vhosts)
	appid_dft, _ := util.AlphaStringRange(24, 32)
	secret_dft, _ := util.AlphaStringRange(32, 36)

	app.Configs.SetSection("forward")
	appid := app.Configs.StringDefault("appid", appid_dft)
	secret := app.Configs.StringDefault("secret", secret_dft)
	acl := app.Configs.BoolDefault("acl.enable", true)
	signature := app.api.sign(secret, appid)
	app.Info("application signature (%s)", signature)
	app.api.acl(acl)

	if err := app.Configs.Save(""); err != nil {
		return err
	}
	return nil
}

func (app *Application) Main() {
	app.background.Work(func() { 
		ext.Serve(app.httpListener, app.reverser, "http", app.Logger) 
		app.Info("background serving reverse worker exiting")
	})
	app.background.Work(func() { 
		ext.Serve(app.apiListener, app.api, "http", app.Logger) 
		app.Info("background serving api worker exiting")
	})

	app.background.Work(func() { 
		c := time.Tick(app.Options.RefreshInterval)
		for {
			select {
			case <- c:
				app.Refresh()
			case <- app.exitc:
				goto timeExit
			}	
		}
	timeExit:
		app.Info("background refresh worker exiting")
	})
}

func (app *Application) Refresh() {
	app.Info("application refresh ...")	
}

func (app *Application) Exit() {
	if app.httpListener != nil {
		app.httpListener.Close()
	}
	if app.apiListener != nil {
		app.apiListener.Close()
	}
	app.Lock()
	// do something if needed
	app.Unlock()

	// notify app to exit
	close(app.exitc)

	// wait all backgroud workers
	app.background.Wait()
}
