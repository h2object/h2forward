package api

import (
	"sync"
	"github.com/h2object/rpc"
)

type Auth interface{
	rpc.PreRequest
}

type Logger interface{
	rpc.Logger
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{}) 
	Error(format string, args ...interface{}) 
	Critical(format string, args ...interface{})
}

var UserAgent = "Golang h2forward/api package"

type Client struct{
	sync.RWMutex
	addr string
	conn *rpc.Client
}

func NewClient(addr string) *Client {
	connection := rpc.NewClient(rpc.H2OAnalyser{})
	clt := &Client{
		addr: addr,	
		conn: connection,
	}
	return clt
}