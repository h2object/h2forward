package app

import (
	"errors"
	"strings"
 	"net/http"
 	"github.com/h2object/h2object/util"
 	"github.com/h2object/h2object/httpext"
)

type HostURL struct{
	Host string `json:"host"`
	URL  string `json:"url"`
}

type APIController struct{
	vhost *VirtualHost
	signature string
	acl_enable bool
}

func NewAPIController(vhost *VirtualHost) *APIController {
	return &APIController{
		vhost: vhost,
		acl_enable: false,
	}
}

func (api *APIController) sign(secret, appid string) string {
	api.signature = util.SignString(secret, appid)
	return api.signature
}

func (api *APIController) acl(flag bool) {
	api.acl_enable = flag
}

func (api *APIController)  ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := ext.NewRequest(req)
	response := ext.NewResponse(w)
	controller := ext.NewController(request, response)

	if strings.ToLower(request.URI()) != "/virtualhost.json" {
		controller.JsonError(http.StatusBadRequest, errors.New("bad request"))
		return 
	}

	// check token for security
	if api.acl_enable == true {
		token := request.Param("token")
		if token == "" {
			authorization := request.Header.Get("Authorization")
			if strings.HasPrefix(authorization, "H2FORWARD ") {
				token = authorization[len("H2FORWARD "):]
			}
		}
		if token != api.signature {
			controller.JsonError(http.StatusUnauthorized, errors.New("request unauthorized"))
			return
		}		
	}

	switch request.MethodToLower() {
	case "get":
		hosts := request.Params("host")
		if len(hosts) == 0 {
			val, err := api.vhost.HostData()
			if err != nil {
				controller.JsonError(http.StatusInternalServerError, err)
				return 
			}

			controller.Json(val)
			return
		}

		results := map[string]interface{}{}
		for _, host := range hosts {
			val, err := api.vhost.GetHost(host)
			if err != nil {
				controller.JsonError(http.StatusInternalServerError, err)
				return 
			}
			results[host] = val
		}
		controller.Json(results)
		return

	case "put":
		// data format: {"host":"", "url":""}
		var hurl HostURL
		if err := request.JsonData(&hurl); err != nil {
			controller.JsonError(http.StatusInternalServerError, err)
			return 
		}
		
		if err := api.vhost.SetHost(hurl.Host, hurl.URL); err != nil {
			controller.JsonError(http.StatusInternalServerError, err)
			return 	
		}

		controller.Json(hurl)
		return
	case "delete":
		hosts := request.Params("host")
		if len(hosts) == 0 {
			controller.JsonError(http.StatusNotImplemented, errors.New("none delete host param"))
			return
		}

		results := map[string]interface{}{}		
		for _, host := range hosts {
			if err := api.vhost.RemoveHost(host); err != nil {
				controller.JsonError(http.StatusInternalServerError, err)
				return
			}
			results[host] = "deleted"
		}
		controller.Json(results)
		return
	}

	controller.JsonError(http.StatusBadRequest, errors.New("unsupport method"))
	return
}