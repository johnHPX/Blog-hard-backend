package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/controller/resource"
)

var commentRoutes = []Router{
	{
		TokenIsReq: true,
		Path:       "/comment/store",
		EndPointer: resource.CommentStoreHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/comment/list/post/id/{postID}",
		EndPointer: resource.CommentListPostHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/comment/list/user/id/{userID}",
		EndPointer: resource.CommentListUserHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/comment/list/post/id/{postID}/user/id/{userID}",
		EndPointer: resource.CommentListPostUserHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/comment/find/id/{id}",
		EndPointer: resource.CommentFindHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/comment/update/id/{id}",
		EndPointer: resource.CommentUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/comment/remove/id/{id}",
		EndPointer: resource.CommentRemoveHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
