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

type postEntity struct {
	PostID  string `json:"postID"`
	Title   string `json:"title"`
	Content string `json:"content`
	Likes   int    `json:"likes"`
}

type postStoreRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	MID     string `json:"mid"`
	Request *http.Request
}

type postStoreResponse struct {
	MID string `json:"mid"`
}

func decodePostStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(postStoreRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makePostStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*postStoreRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewPostService(userToken.UserID, userToken.Kind)
		err = service.Store(req.Title, req.Content)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &postStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func PostStoreHandler() http.Handler {
	return httptransport.NewServer(
		makePostStoreEndPoint(),
		decodePostStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type postListRequest struct {
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type postListResponse struct {
	Count int          `json:"count"`
	Posts []postEntity `json:"posts"`
	MID   string       `json:"mid"`
}

func decodePostListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
	dto := &postListRequest{
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makePostListEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*postListRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		service := service.NewPostService(userToken.UserID, userToken.Kind)
		posts, err := service.List(req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		count, err := service.Count()
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1006, err, req.MID)
		}

		var entities []postEntity
		for _, v := range posts {
			entities = append(entities, postEntity{
				PostID:  v.PostID,
				Title:   v.Title,
				Content: v.Content,
				Likes:   v.Likes,
			})
		}

		return &postListResponse{
			Count: count,
			Posts: entities,
			MID:   req.MID,
		}, nil
	}
}

func PostListHandler() http.Handler {
	return httptransport.NewServer(
		makePostListEndPoint(),
		decodePostListRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type postFindRequest struct {
	ID      string
	MID     string
	Request *http.Request
}

type postFindResponse struct {
	postEntity
	MID string `json:"mid"`
}

func decodePostFindRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := &postFindRequest{
		ID:      id,
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makePostFindendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*postFindRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1007, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1008, err, req.MID)
		}

		service := service.NewPostService(userToken.UserID, userToken.Kind)
		post, err := service.Find(req.ID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1009, err, req.MID)
		}

		return &postFindResponse{
			postEntity: postEntity{
				PostID:  post.PostID,
				Title:   post.Title,
				Content: post.Content,
				Likes:   post.Likes,
			},
			MID: req.MID,
		}, nil
	}
}

func PostFindHandler() http.Handler {
	return httptransport.NewServer(
		makePostFindendPoint(),
		decodePostFindRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type postListTitleRequest struct {
	title   string
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type postListTitleResponse struct {
	Count int          `json:"count"`
	Posts []postEntity `json:"posts"`
	MID   string       `json:"mid"`
}

func decodePostListTitleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	title := vars["title"]
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
	dto := &postListTitleRequest{
		title:   title,
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makePostListTitleEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*postListTitleRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		service := service.NewPostService(userToken.UserID, userToken.Kind)
		posts, err := service.ListTitle(req.title, req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		count, err := service.CountTitle(req.title)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1006, err, req.MID)
		}

		var entities []postEntity
		for _, v := range posts {
			entities = append(entities, postEntity{
				PostID:  v.PostID,
				Title:   v.Title,
				Content: v.Content,
				Likes:   v.Likes,
			})
		}

		return &postListTitleResponse{
			Count: count,
			Posts: entities,
			MID:   req.MID,
		}, nil
	}
}

func PostListTitleHandler() http.Handler {
	return httptransport.NewServer(
		makePostListTitleEndPoint(),
		decodePostListTitleRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type postUpdateRequest struct {
	ID      string
	Title   string `json:"title"`
	Content string `json:"content"`
	MID     string `json:"mid"`
	Request *http.Request
}

type postUpdateResponse struct {
	MID string `json:"mid"`
}

func decodePostUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	dto := new(postUpdateRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.ID = id
	dto.Request = r
	return dto, nil
}

func makePostUpdateEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*postUpdateRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewPostService(userToken.UserID, userToken.Kind)
		err = service.Update(req.ID, req.Title, req.Content)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &postUpdateResponse{
			MID: req.MID,
		}, nil
	}
}

func PostUpdateHandler() http.Handler {
	return httptransport.NewServer(
		makePostUpdateEndPoint(),
		decodePostUpdateRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type postRemoveRequest struct {
	ID      string
	MID     string
	Request *http.Request
}

type postRemoveResponse struct {
	MID string `json:"mid"`
}

func decodePostRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := new(postRemoveRequest)
	dto.ID = id
	dto.MID = mid
	dto.Request = r
	return dto, nil
}

func makePostRemoveEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*postRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewPostService(userToken.UserID, userToken.Kind)
		err = service.Remove(req.ID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &postRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func PostRemoveHandler() http.Handler {
	return httptransport.NewServer(
		makePostRemoveEndPoint(),
		decodePostRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
