package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johnHPX/blog-hard-backend/internal/utils"
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
	for _, router := range routers {
		if router.TokenIsReq {
			s.Router.HandleFunc(router.Path, utils.HeaderMethods(utils.Authenticate(router.EndPointer), router.Method))
		}
		s.Router.HandleFunc(router.Path, utils.HeaderMethods(router.EndPointer, router.Method))
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
