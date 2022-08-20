package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johnHPX/blog-hard-backend/internal/utils/authn"
)

type Router struct {
	TokenIsReq bool
	Path       string
	EndPointer http.HandlerFunc
	Method     string
}

type WebService interface {
	Init()
	GetRouters() http.Handler
}

type webServiceImpl struct {
	Router *mux.Router
}

func (s *webServiceImpl) configuration() {
	routers := []Router{}
	routers = append(routers, userRoutes...)
	routers = append(routers, postRoutes...)
	routers = append(routers, numberLikes...)
	routers = append(routers, commentRoutes...)
	routers = append(routers, categoryRoutes...)
	routers = append(routers, postCategory...)
	routers = append(routers, responseComment...)
	routers = append(routers, configsRoutes...)
	for _, router := range routers {
		if router.TokenIsReq {
			s.Router.HandleFunc(router.Path, authn.HeaderMethods(authn.Authenticate(router.EndPointer), router.Method))
		}
		s.Router.HandleFunc(router.Path, authn.HeaderMethods(router.EndPointer, router.Method))
	}
}

func (s *webServiceImpl) Init() {
	s.configuration()
}

func (s *webServiceImpl) GetRouters() http.Handler {
	return s.Router
}

func NewWebService() WebService {
	return &webServiceImpl{
		Router: mux.NewRouter(),
	}
}
