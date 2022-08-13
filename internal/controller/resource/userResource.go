package resource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
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
