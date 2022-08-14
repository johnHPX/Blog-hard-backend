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
}
