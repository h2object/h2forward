package commands

import (
	"os"
	"fmt"
	"path"
	"encoding/json"
	"path/filepath"
	"github.com/codegangsta/cli"
	"github.com/h2object/h2forward/app"
	"github.com/h2object/h2forward/api"
)

func vhostGetCommand(ctx *cli.Context) {
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

	configs, err := app.LoadCONFIG(path.Join(directory, "h2forward.conf"))
	if err != nil {
		fmt.Println("load configure failed:", err)
		return	
	}
	
	addr := ctx.GlobalString("api")

	configs.SetSection("forward")
	AppID := configs.StringDefault("appid", "")
	Secret := configs.StringDefault("secret", "")

	client := api.NewClient(addr)
	auth := api.NewAdminAuth(AppID, Secret)

	hosts := ctx.Args()
	ret := map[string]interface{}{}
	if err := client.GetHost(nil, auth, hosts, &ret); err != nil {
		fmt.Println("get host failed:", err)
		return	
	}

	b, _ := json.MarshalIndent(ret, "", " ")
	fmt.Println(string(b))
	return	
}

func vhostSetCommand(ctx *cli.Context) {
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

	configs, err := app.LoadCONFIG(path.Join(directory, "h2forward.conf"))
	if err != nil {
		fmt.Println("load configure failed:", err)
		return	
	}

	addr := ctx.GlobalString("api")

	configs.SetSection("forward")
	AppID := configs.StringDefault("appid", "")
	Secret := configs.StringDefault("secret", "")

	client := api.NewClient(addr)
	auth := api.NewAdminAuth(AppID, Secret)

	args := ctx.Args()
	if len(args) != 2 {
		fmt.Println("command args: <hostname> <url>")
		return	
	}

	if err := client.SetHost(nil, auth, args[0], args[1]); err != nil {
		fmt.Println("set host failed:", err)
		return	
	}

	fmt.Printf("set host (%s) to url (%s) ok.\n", args[0], args[1])
	return	
}

func vhostDelCommand(ctx *cli.Context) {
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

	configs, err := app.LoadCONFIG(path.Join(directory, "h2forward.conf"))
	if err != nil {
		fmt.Println("load configure failed:", err)
		return	
	}

	addr := ctx.GlobalString("api")

	configs.SetSection("forward")
	AppID := configs.StringDefault("appid", "")
	Secret := configs.StringDefault("secret", "")

	client := api.NewClient(addr)
	auth := api.NewAdminAuth(AppID, Secret)

	args := ctx.Args()
	if len(args) == 0 {
		fmt.Println("command args: <hostname1> <hostname2> ...")
		return	
	}

	for _, host := range args {
		if err := client.DelHost(nil, auth, host); err != nil {
			fmt.Println("det host failed:", err)
			return	
		}
		fmt.Printf("det host (%s) ok.\n", host)
	}
	return
}