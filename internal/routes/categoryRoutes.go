package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/controller/resource"
)

var categoryRoutes = []Router{
	{
		TokenIsReq: true,
		Path:       "/category/store",
		EndPointer: resource.CategoryStoreHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: true,
		Path:       "/category/list",
		EndPointer: resource.CategoryListHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/category/list/post/id/{postID}",
		EndPointer: resource.CategoryListPostHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/category/find/id/{id}",
		EndPointer: resource.CategoryFindHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/category/update/id/{id}",
		EndPointer: resource.CategoryUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/category/remove/id/{id}",
		EndPointer: resource.CategoryRemoveHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
