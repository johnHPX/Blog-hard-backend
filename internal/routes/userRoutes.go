package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/controller/resource"
)

var userRoutes = []Router{
	{
		TokenIsReq: false,
		Path:       "/user/store",
		EndPointer: resource.UserStoreHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/user/list",
		EndPointer: resource.UserListHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/user/list/name/{name}",
		EndPointer: resource.UserListNameHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/user/find/id/{id}",
		EndPointer: resource.UserFindHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/user/update/id/{id}",
		EndPointer: resource.UserUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: false,
		Path:       "/user/remove/id/{id}",
		EndPointer: resource.UserRemoveHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
