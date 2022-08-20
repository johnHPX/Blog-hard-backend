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

type categoryEntity struct {
	CategoryID string `json:"categoryID"`
	Name       string `json:"name"`
}

type categoryStoreRequest struct {
	Name    string `json:"name"`
	MID     string `json:"mid"`
	Request *http.Request
}

type categoryStoreResponse struct {
	MID string `json:"mid"`
}

func decodeCategoryStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(categoryStoreRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makeCategoryStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*categoryStoreRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewCategoryService(userToken.UserID, userToken.Kind)
		err = service.CreateCategory(req.Name)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &postStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func CategoryStoreHandler() http.Handler {
	return httptransport.NewServer(
		makeCategoryStoreEndPoint(),
		decodeCategoryStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type categoryListRequest struct {
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type categoryListResponse struct {
	Count     int              `json:"count"`
	Categorys []categoryEntity `json:"categorys"`
	MID       string           `json:"mid"`
}

func decodeCategoryListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
	dto := &categoryListRequest{
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeCategoryListEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*categoryListRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		service := service.NewCategoryService(userToken.UserID, userToken.Kind)
		category, count, err := service.ListCategory(req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		var entities []categoryEntity
		for _, v := range category {
			entities = append(entities, categoryEntity{
				CategoryID: v.CategoryID,
				Name:       v.Name,
			})
		}

		return &categoryListResponse{
			Count:     count,
			Categorys: entities,
			MID:       req.MID,
		}, nil
	}
}

func CategoryListHandler() http.Handler {
	return httptransport.NewServer(
		makeCategoryListEndPoint(),
		decodeCategoryListRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type categoryListPostRequest struct {
	PostID  string
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type categoryListPostResponse struct {
	Count     int              `json:"count"`
	Categorys []categoryEntity `json:"categorys"`
	MID       string           `json:"mid"`
}

func decodeCategoryListPostRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	postID := vars["postID"]
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
	dto := &categoryListPostRequest{
		PostID:  postID,
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeCategoryListPostEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*categoryListPostRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		service := service.NewCategoryService(userToken.UserID, userToken.Kind)
		category, count, err := service.ListCategoryByPost(req.PostID, req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		var entities []categoryEntity
		for _, v := range category {
			entities = append(entities, categoryEntity{
				CategoryID: v.CategoryID,
				Name:       v.Name,
			})
		}

		return &categoryListResponse{
			Count:     count,
			Categorys: entities,
			MID:       req.MID,
		}, nil
	}
}

func CategoryListPostHandler() http.Handler {
	return httptransport.NewServer(
		makeCategoryListPostEndPoint(),
		decodeCategoryListPostRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type categoryFindRequest struct {
	categoryID string
	MID        string
	Request    *http.Request
}

type categoryFindResponse struct {
	categoryEntity
	MID string `json:"mid"`
}

func decodeCategoryFindRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := &categoryFindRequest{
		categoryID: id,
		MID:        mid,
		Request:    r,
	}
	return dto, nil
}

func makeCategoryFindEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*categoryFindRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1006, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1007, err, req.MID)
		}

		service := service.NewCategoryService(userToken.UserID, userToken.Kind)
		category, err := service.FindCategory(req.categoryID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1008, err, req.MID)
		}

		return &categoryFindResponse{
			categoryEntity: categoryEntity{
				CategoryID: category.CategoryID,
				Name:       category.Name,
			},
			MID: req.MID,
		}, nil
	}
}

func CategoryFindHandler() http.Handler {
	return httptransport.NewServer(
		makeCategoryFindEndPoint(),
		decodeCategoryFindRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type categoryUpdateRequest struct {
	categoryID string
	Name       string `json:"name"`
	MID        string `json:"mid"`
	Request    *http.Request
}

type categoryUpdateResponse struct {
	MID string `json:"mid"`
}

func decodeCategoryUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	dto := new(categoryUpdateRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.categoryID = id
	dto.Request = r
	return dto, nil
}

func makeCategoryUpdateEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*categoryUpdateRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1009, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1010, err, req.MID)
		}

		service := service.NewCategoryService(userToken.UserID, userToken.Kind)
		err = service.UpdateCategory(req.categoryID, req.Name)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1011, err, req.MID)
		}

		return &commentUpdateResponse{
			MID: req.MID,
		}, nil
	}
}

func CategoryUpdateHandler() http.Handler {
	return httptransport.NewServer(
		makeCategoryUpdateEndPoint(),
		decodeCategoryUpdateRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type categoryRemoveRequest struct {
	ID      string
	MID     string
	Request *http.Request
}

type categoryRemoveResponse struct {
	MID string `json:"mid"`
}

func decodeCategoryRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := new(categoryRemoveRequest)
	dto.ID = id
	dto.MID = mid
	dto.Request = r
	return dto, nil
}

func makeCategoryRemoveEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*categoryRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1012, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1013, err, req.MID)
		}

		service := service.NewCategoryService(userToken.UserID, userToken.Kind)
		err = service.RemoveCategory(req.ID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1014, err, req.MID)
		}

		return &postRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func CategoryRemoveHandler() http.Handler {
	return httptransport.NewServer(
		makeCategoryRemoveEndPoint(),
		decodeCategoryRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
