package resource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/johnHPX/blog-hard-backend/internal/controller/service"
	"github.com/johnHPX/blog-hard-backend/internal/utils/responseAPI"
)

type userStoreResquest struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Nick      string `json:"nick"`
	Email     string `json:"email"`
	Secret    string `json:"secret"`
	MID       string `json:"mid"`
}

type userStoreResponse struct {
	MID string `json:"mid"`
}

func decodeUserStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(userStoreResquest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func makeUserStoreendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userStoreResquest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		service := service.NewUserService("", "")
		err := service.Store(req.Name, req.Telephone, req.Nick, req.Email, req.Secret)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1001, err, req.MID)
		}

		return &userStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func UserStoreHandler() http.Handler {
	return httptransport.NewServer(
		makeUserStoreendPoint(),
		decodeUserStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type userStoreADMResquest struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Nick      string `json:"nick"`
	Email     string `json:"email"`
	Secret    string `json:"secret"`
	Kind      string `json:"kind"`
	MID       string `json:"mid"`
	Request   *http.Request
}

type userStoreADMResponse struct {
	MID string `json:"mid"`
}

func decodeUserStoreADMRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(userStoreADMResquest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makeUserStoreADMendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userStoreADMResquest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1002, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		service := service.NewUserService(userToken.UserID, userToken.Kind)
		err = service.StoreADM(req.Name, req.Telephone, req.Nick, req.Email, req.Secret, req.Kind)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		return &userStoreADMResponse{
			MID: req.MID,
		}, nil
	}
}

func UserStoreADMHandler() http.Handler {
	return httptransport.NewServer(
		makeUserStoreADMendPoint(),
		decodeUserStoreADMRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type userEntity struct {
	PersonID  string `json:"personID"`
	UserID    string `json:"userID"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Nick      string `json:"nick"`
	Email     string `json:"email"`
	Kind      string `json:"kind"`
}

type userListRequest struct {
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type userListResponse struct {
	Count int          `json:"count"`
	Users []userEntity `json:"users"`
	MID   string       `json:"mid"`
}

func decodeUserListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 1
	}
	mid := r.URL.Query().Get("mid")
	dto := &userListRequest{
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeUserListendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userListRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1006, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1007, err, req.MID)
		}

		service := service.NewUserService(userToken.UserID, userToken.Kind)
		users, err := service.List(req.Offset, req.Limit, req.Page)

		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1008, err, req.MID)
		}

		count, err := service.Count()
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1009, err, req.MID)
		}

		var entities []userEntity
		for _, v := range users {
			entities = append(entities, userEntity{
				PersonID:  v.PersonID,
				UserID:    v.UserID,
				Name:      v.Name,
				Telephone: v.Telephone,
				Nick:      v.Nick,
				Email:     v.Email,
				Kind:      v.Kind,
			})
		}

		return &userListResponse{
			Count: count,
			Users: entities,
			MID:   req.MID,
		}, nil
	}
}

func UserListHandler() http.Handler {
	return httptransport.NewServer(
		makeUserListendPoint(),
		decodeUserListRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type userListNameRequest struct {
	Name    string
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type userListNameResponse struct {
	Count int          `json:"count"`
	Users []userEntity `json:"users"`
	MID   string       `json:"mid"`
}

func decodeUserListNameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name := vars["name"]
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 1
	}
	mid := r.URL.Query().Get("mid")
	dto := &userListNameRequest{
		Name:    name,
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeUserListNameendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userListNameRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1010, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1011, err, req.MID)
		}

		service := service.NewUserService(userToken.UserID, userToken.Kind)
		users, err := service.ListName(req.Name, req.Offset, req.Limit, req.Page)

		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1012, err, req.MID)
		}

		count, err := service.CountName(req.Name)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1013, err, req.MID)
		}

		var entities []userEntity
		for _, v := range users {
			entities = append(entities, userEntity{
				PersonID:  v.PersonID,
				UserID:    v.UserID,
				Name:      v.Name,
				Telephone: v.Telephone,
				Nick:      v.Nick,
				Email:     v.Email,
				Kind:      v.Kind,
			})
		}

		return &userListResponse{
			Count: count,
			Users: entities,
			MID:   req.MID,
		}, nil
	}
}

func UserListNameHandler() http.Handler {
	return httptransport.NewServer(
		makeUserListNameendPoint(),
		decodeUserListNameRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type userFindRequest struct {
	ID      string
	MID     string
	Request *http.Request
}

type userFindResponse struct {
	userEntity
	MID string `json:"mid"`
}

func decodeUserFindRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := &userFindRequest{
		ID:      id,
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeUserFindendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userFindRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1014, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1015, err, req.MID)
		}

		service := service.NewUserService(userToken.UserID, userToken.Kind)
		user, err := service.Find(req.ID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1016, err, req.MID)
		}

		return &userFindResponse{
			userEntity: userEntity{
				PersonID:  user.PersonID,
				UserID:    user.UserID,
				Name:      user.Name,
				Telephone: user.Telephone,
				Nick:      user.Nick,
				Email:     user.Email,
				Kind:      user.Kind,
			},
			MID: req.MID,
		}, nil
	}
}

func UserFindHandler() http.Handler {
	return httptransport.NewServer(
		makeUserFindendPoint(),
		decodeUserFindRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type userUpdateRequest struct {
	ID        string
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Nick      string `json:"nick"`
	Email     string `json:"email"`
	Kind      string `json:"kind"`
	MID       string `json:"mid"`
	Request   *http.Request
}

type userUpdateResponse struct {
	MID string `json:"mid"`
}

func decodeUserUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	dto := new(userUpdateRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.ID = id
	dto.Request = r
	return dto, nil
}

func makeUserUpdateendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userUpdateRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1017, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1018, err, req.MID)
		}

		service := service.NewUserService(userToken.UserID, userToken.Kind)
		err = service.Update(req.ID, req.Name, req.Telephone, req.Nick, req.Email, req.Kind)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1019, err, req.MID)
		}

		return &userUpdateResponse{
			MID: req.MID,
		}, nil
	}
}

func UserUpdateHandler() http.Handler {
	return httptransport.NewServer(
		makeUserUpdateendPoint(),
		decodeUserUpdateRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type userRemoveRequest struct {
	ID      string
	MID     string
	Request *http.Request
}

type userRemoveResponse struct {
	MID string `json:"mid"`
}

func decodeUserRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := &userRemoveRequest{
		ID:      id,
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeUserRemoveendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1020, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1021, err, req.MID)
		}

		service := service.NewUserService(userToken.UserID, userToken.Kind)
		err = service.Remove(req.ID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1022, err, req.MID)
		}

		return &userRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func UserRemoveHandler() http.Handler {
	return httptransport.NewServer(
		makeUserRemoveendPoint(),
		decodeUserRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type userLoginRequest struct {
	EmailOrNick string `json:"nick"`
	Secret      string `json:"password"`
	MID         string `json:"mid"`
}

type userLoginResponse struct {
	Token string `json:"token"`
	MID   string `json:"mid"`
}

func decodeUserLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(userLoginRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func makeUserLoginendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userLoginRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1023, errors.New("invalid request"), "na")
		}

		service := service.NewUserService("", "")
		token, err := service.Login(req.EmailOrNick, req.Secret)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1024, err, req.MID)
		}

		return &userLoginResponse{
			Token: token,
			MID:   req.MID,
		}, nil
	}
}

func UserLoginHandler() http.Handler {
	return httptransport.NewServer(
		makeUserLoginendPoint(),
		decodeUserLoginRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
