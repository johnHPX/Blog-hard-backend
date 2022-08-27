package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/interf/resource"
)

var numberLikes = []Router{
	{
		TokenIsReq: true,
		Path:       "/user/post/like",
		EndPointer: resource.NumberLikesStoreHandle().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: true,
		Path:       "/user/post/dislike",
		EndPointer: resource.NumberLikesRemoveHandle().ServeHTTP,
		Method:     http.MethodPost,
	},
}
