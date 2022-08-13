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
}
