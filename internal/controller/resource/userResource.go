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
	"github.com/johnHPX/blog-hard-backend/internal/utils"
)

type userStoreResquest struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Nick      string `json:"nick"`
	Email     string `json:"email"`
	Secret    string `json:"secret"`
	Kind      string `json:"kind"`
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
			return nil, utils.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		service := service.NewUserService()

		err := service.Store(req.Name, req.Telephone, req.Nick, req.Email, req.Secret, req.Kind)

		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1001, err, req.MID)
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
		utils.EncodeResponse,
		httptransport.ServerErrorEncoder(utils.ErrorEncoder()),
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
	Offset int
	Limit  int
	Page   int
	MID    string
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
		Offset: int(offset),
		Limit:  int(limit),
		Page:   int(page),
		MID:    mid,
	}
	return dto, nil
}

func makeUserListendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userListRequest)
		if !ok {
			return nil, utils.CreateHttpErrorResponse(http.StatusBadRequest, 1002, errors.New("invalid request"), "na")
		}

		service := service.NewUserService()
		users, err := service.List(req.Offset, req.Limit, req.Page)

		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1003, err, req.MID)
		}

		count, err := service.Count()
		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1004, err, req.MID)
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
		utils.EncodeResponse,
		httptransport.ServerErrorEncoder(utils.ErrorEncoder()),
	)
}

type userListNameRequest struct {
	Name   string
	Offset int
	Limit  int
	Page   int
	MID    string
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
		Name:   name,
		Offset: int(offset),
		Limit:  int(limit),
		Page:   int(page),
		MID:    mid,
	}
	return dto, nil
}

func makeUserListNameendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userListNameRequest)
		if !ok {
			return nil, utils.CreateHttpErrorResponse(http.StatusBadRequest, 1005, errors.New("invalid request"), "na")
		}

		service := service.NewUserService()
		users, err := service.ListName(req.Name, req.Offset, req.Limit, req.Page)

		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1006, err, req.MID)
		}

		count, err := service.CountName(req.Name)
		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1007, err, req.MID)
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
		utils.EncodeResponse,
		httptransport.ServerErrorEncoder(utils.ErrorEncoder()),
	)
}

type userFindRequest struct {
	ID  string
	MID string
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
		ID:  id,
		MID: mid,
	}
	return dto, nil
}

func makeUserFindendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userFindRequest)
		if !ok {
			return nil, utils.CreateHttpErrorResponse(http.StatusBadRequest, 1008, errors.New("invalid request"), "na")
		}

		service := service.NewUserService()
		user, err := service.Find(req.ID)
		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1009, err, req.MID)
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
		utils.EncodeResponse,
		httptransport.ServerErrorEncoder(utils.ErrorEncoder()),
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
	return dto, nil
}

func makeUserUpdateendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userUpdateRequest)
		if !ok {
			return nil, utils.CreateHttpErrorResponse(http.StatusBadRequest, 1010, errors.New("invalid request"), "na")
		}

		service := service.NewUserService()
		err := service.Update(req.ID, req.Name, req.Telephone, req.Nick, req.Email, req.Kind)
		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1011, err, req.MID)
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
		utils.EncodeResponse,
		httptransport.ServerErrorEncoder(utils.ErrorEncoder()),
	)
}

type userRemoveRequest struct {
	ID  string
	MID string
}

type userRemoveResponse struct {
	MID string `json:"mid"`
}

func decodeUserRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := &userRemoveRequest{
		ID:  id,
		MID: mid,
	}
	return dto, nil
}

func makeUserRemoveendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*userRemoveRequest)
		if !ok {
			return nil, utils.CreateHttpErrorResponse(http.StatusBadRequest, 1012, errors.New("invalid request"), "na")
		}

		service := service.NewUserService()
		err := service.Remove(req.ID)
		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1013, err, req.MID)
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
		utils.EncodeResponse,
		httptransport.ServerErrorEncoder(utils.ErrorEncoder()),
	)
}
