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

type responseCommentEntity struct {
	ResponseCommentID string `json:"responseCommentId"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	CommentID         string `json:"commentId"`
	UserID            string `json:"userId"`
}

type responseCommentStoreRequest struct {
	CommentID string `json:"commentId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	MID       string `json:"mid"`
	Request   *http.Request
}

type responseCommentStoreResponse struct {
	MID string `json:"mid"`
}

func decodeResponseCommentStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(responseCommentStoreRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makeResponseCommentStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*responseCommentStoreRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewResponseCommentService(userToken.UserID, userToken.Kind)
		err = service.Store(req.CommentID, req.Title, req.Content)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &responseCommentStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func ResponseCommentStoreHandler() http.Handler {
	return httptransport.NewServer(
		makeResponseCommentStoreEndPoint(),
		decodeResponseCommentStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type responseCommentListRequest struct {
	CommentID string
	Offset    int
	Limit     int
	Page      int
	MID       string
	Request   *http.Request
}

type responseCommentListResponse struct {
	Count           int                     `json:"count"`
	ReponseComments []responseCommentEntity `json:"responseComments"`
	MID             string                  `json:"mid"`
}

func decoderesponseCommentListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	commentID := vars["commentID"]
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
	dto := &responseCommentListRequest{
		CommentID: commentID,
		Offset:    int(offset),
		Limit:     int(limit),
		Page:      int(page),
		MID:       mid,
		Request:   r,
	}
	return dto, nil
}

func makeResponseCommentListEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*responseCommentListRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewResponseCommentService(userToken.UserID, userToken.Kind)
		responseComments, count, err := service.List(req.CommentID, req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1004, err, req.MID)
		}

		var entities []responseCommentEntity
		for _, v := range responseComments {
			entities = append(entities, responseCommentEntity{
				CommentID:         v.CommentID,
				Title:             v.Title,
				Content:           v.Content,
				UserID:            v.UserID,
				ResponseCommentID: v.ResponseCommentID,
			})
		}

		return &responseCommentListResponse{
			Count:           count,
			ReponseComments: entities,
			MID:             req.MID,
		}, nil
	}
}

func ResponseCommentListHandler() http.Handler {
	return httptransport.NewServer(
		makeResponseCommentListEndPoint(),
		decoderesponseCommentListRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type responseCommentListUserRequest struct {
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type responseCommentListUserResponse struct {
	Count           int                     `json:"count"`
	ReponseComments []responseCommentEntity `json:"responseComments"`
	MID             string                  `json:"mid"`
}

func decoderesponseCommentListUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
	dto := &responseCommentListUserRequest{
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeResponseCommentListUserEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*responseCommentListUserRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewResponseCommentService(userToken.UserID, userToken.Kind)
		responseComments, count, err := service.ListUser(req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1004, err, req.MID)
		}

		var entities []responseCommentEntity
		for _, v := range responseComments {
			entities = append(entities, responseCommentEntity{
				CommentID:         v.CommentID,
				Title:             v.Title,
				Content:           v.Content,
				UserID:            v.UserID,
				ResponseCommentID: v.ResponseCommentID,
			})
		}

		return &responseCommentListUserResponse{
			Count:           count,
			ReponseComments: entities,
			MID:             req.MID,
		}, nil
	}
}

func ResponseCommentListUserHandler() http.Handler {
	return httptransport.NewServer(
		makeResponseCommentListUserEndPoint(),
		decoderesponseCommentListUserRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type responseCommentUpdateRequest struct {
	ResponseCommentID string
	Title             string `json:"title"`
	Content           string `json:"content"`
	MID               string `json:"mid"`
	Request           *http.Request
}

type responseCommentUpdateResponse struct {
	MID string `json:"mid"`
}

func decoderesponseCommentUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	dto := new(responseCommentUpdateRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.ResponseCommentID = id
	dto.Request = r
	return dto, nil
}

func makeResponseCommentUpdateEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*responseCommentUpdateRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1014, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1015, err, req.MID)
		}

		service := service.NewResponseCommentService(userToken.UserID, userToken.Kind)
		err = service.Update(req.ResponseCommentID, req.Title, req.Content)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1016, err, req.MID)
		}

		return &responseCommentUpdateResponse{
			MID: req.MID,
		}, nil
	}
}

func ResponseCommentUpdateHandler() http.Handler {
	return httptransport.NewServer(
		makeResponseCommentUpdateEndPoint(),
		decoderesponseCommentUpdateRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type responseCommentRemoveRequest struct {
	ID      string
	MID     string
	Request *http.Request
}

type responseCommentRemoveResponse struct {
	MID string `json:"mid"`
}

func decoderesponseCommentRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := new(responseCommentRemoveRequest)
	dto.ID = id
	dto.MID = mid
	dto.Request = r
	return dto, nil
}

func makeResponseCommentRemoveEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*responseCommentRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1017, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1018, err, req.MID)
		}

		service := service.NewResponseCommentService(userToken.UserID, userToken.Kind)
		err = service.Remove(req.ID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1019, err, req.MID)
		}

		return &responseCommentRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func ResponseCommentRemoveHandler() http.Handler {
	return httptransport.NewServer(
		makeResponseCommentRemoveEndPoint(),
		decoderesponseCommentRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
