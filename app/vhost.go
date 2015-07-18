package app

import (
	"fmt"
	"sync"
	"path"
	"strings"
	"net/http"
	"github.com/h2object/h2object/log"
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/object"
    "github.com/h2object/oxy/reverse"
)

type VirtualHost struct{
	sync.RWMutex 
	log.Logger
	root 	string
	hosts 	map[string]reverse.Endpoint
	store 	object.Store
}

func NewVirtualHost(root string, l log.Logger) *VirtualHost {
	return &VirtualHost{
		Logger: l,
		root: root,
		hosts: make(map[string]reverse.Endpoint),
	}
}

func (vh *VirtualHost) Load() error {
	store := object.NewBoltStore(vh.root, "vhost.dat", object.BoltCoder{})
	if err := store.Load(); err != nil {
		return err
	}
	vh.store = store
	return nil
}


func (vh *VirtualHost) Route(req *http.Request) (reverse.Endpoint, error) {
	if vh.store == nil {
		return nil, fmt.Errorf("virtual hosts not loaded.")
	}
	vh.RLock()
	defer vh.RUnlock()

	hostname := strings.Split(strings.ToLower(req.Host), ":")[0]
	matcher, exists := vh.hosts[hostname]
	if !exists {
		val , err := vh.store.Get(path.Join("/", hostname), true)
		if err != nil {
			return nil, err
		}

		var url string
		if err := util.Convert(val, &url); err != nil {
			return nil, err
		}

		endpoint, err := reverse.ParseUrl(url)
		if err != nil {
			return nil, err
		}

		vh.hosts[hostname] = endpoint
		return endpoint, nil
	}
	return matcher, nil
}

func (vh *VirtualHost) HostData() (interface{}, error) {
	if vh.store == nil {
		return nil, fmt.Errorf("virtual hosts not loaded.")
	}
	vh.RLock()
	defer vh.RUnlock()

	return vh.store.Get("/", true)
}

func (vh *VirtualHost) SetHost(hostname, url string) error {
	if vh.store == nil {
		return fmt.Errorf("virtual hosts not loaded.")
	}
	vh.Lock()
	defer vh.Unlock()

	hostname = strings.ToLower(hostname)

	endpoint, err := reverse.ParseUrl(url)
	if err != nil {
		return err
	}

	if err := vh.store.Put(path.Join("/", hostname), url); err != nil {
		return err
	}

	vh.hosts[hostname] = endpoint
	vh.Info("virtual hosts set (%s) upstream to (%s)", hostname, url)
	return nil
}

func (vh *VirtualHost) GetHost(hostname string) (string, error) {
	if vh.store == nil {
		return "", fmt.Errorf("virtual hosts not loaded.")
	}
	vh.RLock()
	defer vh.RUnlock()

	hostname = strings.ToLower(hostname)
	val, err := vh.store.Get(path.Join("/", hostname), true)
	if err != nil {
		return "", err
	}

	var url string
	if err := util.Convert(val, &url); err != nil {
		return "", err
	}
	return url, nil
}


func (vh *VirtualHost) RemoveHost(hostname string) error {
	if vh.store == nil {
		return fmt.Errorf("virtual hosts not loaded.")
	}
	vh.Lock()
	defer vh.Unlock()

	hostname = strings.ToLower(hostname)

	if err := vh.store.Del(path.Join("/", hostname)); err != nil {
		return err
	}

	delete(vh.hosts, hostname)
	vh.Info("virtual hosts del (%s)", hostname)
	return nil
}	
