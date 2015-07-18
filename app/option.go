package app

import (
	"os"
	"time"
	"path"
	"path/filepath"
)

type Options struct{
	HTTPAddress  		string
	APIAddress   		string
	Root 				string
	LogsRoot 			string
	RefreshInterval 	time.Duration
}

func NewOptions(http, api string) *Options {
	return &Options{
		HTTPAddress: http,
		APIAddress: api,
		RefreshInterval: 10*time.Minute,
	}
}

func (opt *Options) Prepare(workdir string) error {
	directory, err := filepath.Abs(workdir)
	if err != nil {
		return err
	}
	opt.Root = directory
	opt.LogsRoot = path.Join(directory, "logs")
	if err := os.MkdirAll(opt.Root, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.LogsRoot, os.ModePerm); err != nil {
		return err
	}
	return nil
}


func (opt *Options) SetRefreshDefault(s string, default_refresh time.Duration) {
	if d, err := time.ParseDuration(s); err == nil {
		opt.RefreshInterval = d
	} else {
		opt.RefreshInterval = default_refresh
	}
}