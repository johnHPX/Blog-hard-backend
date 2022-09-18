package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/interf/resource"
)

var userRoutes = []Router{
	{
		TokenIsReq: false,
		Path:       "/user/store",
		EndPointer: resource.UserStoreHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: true,
		Path:       "/user/adm/store",
		EndPointer: resource.UserStoreADMHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: true,
		Path:       "/user/list",
		EndPointer: resource.UserListHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/user/list/name/{name}",
		EndPointer: resource.UserListNameHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/user/find/id/{id}",
		EndPointer: resource.UserFindHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/user/update/id/{id}",
		EndPointer: resource.UserUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/user/remove/id/{id}",
		EndPointer: resource.UserRemoveHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
	{
		TokenIsReq: false,
		Path:       "/user/login",
		EndPointer: resource.UserLoginHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/user/recor/email",
		EndPointer: resource.UserSendEmailHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/user/verific/code",
		EndPointer: resource.UserVerificCodeHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: true,
		Path:       "/user/password/recovery",
		EndPointer: resource.UserPasswordRecoveryHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/user/password/update",
		EndPointer: resource.UserPasswordUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/user/logout",
		EndPointer: resource.UserLogoutHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
}
