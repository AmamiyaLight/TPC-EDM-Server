package api

import (
	"TPC-EDM-Server/api/customer_api"
	"TPC-EDM-Server/api/lineitem_api"
	"TPC-EDM-Server/api/nation_api"
	"TPC-EDM-Server/api/orders_api"
	"TPC-EDM-Server/api/part_api"
	"TPC-EDM-Server/api/part_supp_api"
	"TPC-EDM-Server/api/region_api"
	"TPC-EDM-Server/api/supplier_api"
	"TPC-EDM-Server/api/tpcc_api"
	"TPC-EDM-Server/api/user_api"
)

type Api struct {
	UserApi     user_api.UserApi
	OrdersApi   orders_api.OrdersApi
	PartSuppApi part_supp_api.PartSuppApi
	LineItemApi lineitem_api.LineItemApi
	CustomerApi customer_api.CustomerApi
	NationApi   nation_api.NationApi
	PartApi     part_api.PartApi
	SupplierApi supplier_api.SupplierApi
	RegionApi   region_api.RegionApi
	TpccApi     tpcc_api.TpccApi
}

var App = Api{}
