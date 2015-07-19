package api

import (
	"github.com/h2object/rpc"
	. "github.com/h2object/h2forward/app"
	"net/url"
	"errors"
)

func (h2o *Client) SetHost(l Logger, auth Auth, hostname string, url string) error {
	URL := rpc.BuildHttpURL(h2o.addr, "/virtualhost.json", nil)

	h2o.Lock()
	defer h2o.Unlock()

	h2o.conn.Prepare(auth)
	defer h2o.conn.Prepare(nil)

	var hurl HostURL
	hurl.Host = hostname
	hurl.URL = url

	if err := h2o.conn.PutJson(l, URL, hurl, nil); err != nil {
		return err
	}
	return nil
}

func (h2o *Client) GetHost(l Logger, auth Auth, hosts []string, ret interface{}) error {
	var params url.Values = nil
	if len(hosts) > 0 {
		params = url.Values{
			"host": hosts,
		}
	}

	URL := rpc.BuildHttpURL(h2o.addr, "/virtualhost.json", params)

	h2o.Lock()
	defer h2o.Unlock()

	h2o.conn.Prepare(auth)
	defer h2o.conn.Prepare(nil)

	if err := h2o.conn.Get(l, URL, ret); err != nil {
		return err
	}
	return nil
}

func (h2o *Client) DelHost(l Logger, auth Auth, hostname string) error {
	if hostname == "" {
		return errors.New("hostname absent")
	}
	params := url.Values{
			"host": {hostname},
		}
	URL := rpc.BuildHttpURL(h2o.addr, "/virtualhost.json", params)

	h2o.Lock()
	defer h2o.Unlock()

	h2o.conn.Prepare(auth)
	defer h2o.conn.Prepare(nil)

	if err := h2o.conn.Delete(l, URL, nil); err != nil {
		return err
	}
	return nil
}



