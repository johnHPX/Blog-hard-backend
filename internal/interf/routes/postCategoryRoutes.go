 package routes

import (
	"net/http"

	"github.com/johnHPX/blog-hard-backend/internal/interf/resource"
)

var postCategory = []Router{
	{
		TokenIsReq: true,
		Path:       "/post/category/store",
		EndPointer: resource.PostCategoryStoreHandle().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: true,
		Path:       "/post/category/remove",
		EndPointer: resource.PostCategoryRemoveHandle().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
