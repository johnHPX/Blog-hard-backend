package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/controller/resource"
)

var postRoutes = []Router{
	{
		TokenIsReq: true,
		Path:       "/post/store",
		EndPointer: resource.PostStoreHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/post/list",
		EndPointer: resource.PostListHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/post/list/title/{title}",
		EndPointer: resource.PostListTitleHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/post/list/category/name/{category}",
		EndPointer: resource.PostListCategoryHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/post/find/id/{id}",
		EndPointer: resource.PostFindHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/post/update/id/{id}",
		EndPointer: resource.PostUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/post/remove/id/{id}",
		EndPointer: resource.PostRemoveHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
