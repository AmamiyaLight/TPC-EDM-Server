package api

import (
	"TPC-H-EDM-Server/api/lineitem_api"
	"TPC-H-EDM-Server/api/orders_api"
	"TPC-H-EDM-Server/api/part_supp_api"
	"TPC-H-EDM-Server/api/user_api"
)

type Api struct {
	UserApi     user_api.UserApi
	OrdersApi   orders_api.OrdersApi
	PartSuppApi part_supp_api.PartSuppApi
	LineItemApi lineitem_api.LineItemApi
}

var App = Api{}
