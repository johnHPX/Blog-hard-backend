package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/interf/resource"
)

var configsRoutes = []Router{
	{
		TokenIsReq: true,
		Path:       "/config/store",
		EndPointer: resource.ConfigsStoreHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/config/list",
		EndPointer: resource.ConfigsListHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/config/find/id/{id}",
		EndPointer: resource.ConfigsFindHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/config/update/id/{id}",
		EndPointer: resource.ConfigsUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/config/remove/id/{id}",
		EndPointer: resource.ConfigsRemoveHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
