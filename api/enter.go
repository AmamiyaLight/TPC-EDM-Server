package api

import "TPC-H-EDM-Server/api/user_api"

type Api struct {
	UserApi user_api.UserApi
}

var App = Api{}
