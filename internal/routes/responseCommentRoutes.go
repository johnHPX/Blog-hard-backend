package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/controller/resource"
)

var responseComment = []Router{
	{
		TokenIsReq: true,
		Path:       "/response/comment/store",
		EndPointer: resource.ResponseCommentStoreHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/response/comment/list/comment/id/{commentID}",
		EndPointer: resource.ResponseCommentListHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/response/comment/list/user/id/{userID}",
		EndPointer: resource.ResponseCommentListUserHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: true,
		Path:       "/response/comment/update/id/{id}",
		EndPointer: resource.ResponseCommentUpdateHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: true,
		Path:       "/response/comment/remove/id/{id}",
		EndPointer: resource.ResponseCommentRemoveHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
