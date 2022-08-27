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

	"github.com/johnHPX/blog-hard-backend/internal/appl/service"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/responseAPI"
)

type configsEntity struct {
	ConfigID  int      `json:"configID"`
	Collors   []string `json:"collors"`
	Links     []string `json:"links"`
	MenuAs    []string `json:"menuAncoras"`
	BannerURL string   `json:"banner"`
}

type configsStoreRequest struct {
	Collors   []string `json:"collors"`
	Links     []string `json:"links"`
	MenuAs    []string `json:"menuAncoras"`
	BannerURL string   `json:"banner"`
	MID       string   `json:"mid"`
	Request   *http.Request
}

type configsStoreResponse struct {
	MID string `json:"mid"`
}

func decodeConfigsStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(configsStoreRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makeConfigsStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*configsStoreRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewConfigsService(userToken.UserID, userToken.Kind)
		err = service.Store(req.Collors, req.Links, req.MenuAs, req.BannerURL)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &configsStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func ConfigsStoreHandler() http.Handler {
	return httptransport.NewServer(
		makeConfigsStoreEndPoint(),
		decodeConfigsStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type configsListRequest struct {
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type configsListResponse struct {
	Count   int             `json:"count"`
	Configs []configsEntity `json:"configs"`
	MID     string          `json:"mid"`
}

func decodeConfigsListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
	dto := &configsListRequest{
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeConfigsListEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*configsListRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		service := service.NewConfigsService(userToken.UserID, userToken.Kind)
		configs, count, err := service.List(req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		var entities []configsEntity
		for _, v := range configs {
			entities = append(entities, configsEntity{
				ConfigID:  int(v.ConfigID),
				Collors:   v.Collors,
				Links:     v.Links,
				MenuAs:    v.MenuAs,
				BannerURL: v.BannerURL,
			})
		}

		return &configsListResponse{
			Count:   count,
			Configs: entities,
			MID:     req.MID,
		}, nil
	}
}

func ConfigsListHandler() http.Handler {
	return httptransport.NewServer(
		makeConfigsListEndPoint(),
		decodeConfigsListRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type configsFindRequest struct {
	ID      uint
	MID     string
	Request *http.Request
}

type configsFindResponse struct {
	configsEntity
	MID string `json:"mid"`
}

func decodeConfigsFindRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	idConv, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	mid := r.URL.Query().Get("mid")
	dto := &configsFindRequest{
		ID:      uint(idConv),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeConfigsFindendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*configsFindRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1006, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1007, err, req.MID)
		}

		service := service.NewConfigsService(userToken.UserID, userToken.Kind)
		config, err := service.Find(int(req.ID))
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1008, err, req.MID)
		}

		return &configsFindResponse{
			configsEntity: configsEntity{
				ConfigID:  int(config.ConfigID),
				Collors:   config.Collors,
				Links:     config.Links,
				MenuAs:    config.MenuAs,
				BannerURL: config.BannerURL,
			},
			MID: req.MID,
		}, nil
	}
}

func ConfigsFindHandler() http.Handler {
	return httptransport.NewServer(
		makeConfigsFindendPoint(),
		decodeConfigsFindRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type configsUpdateRequest struct {
	ID        uint
	Collors   []string `json:"collors"`
	Links     []string `json:"links"`
	MenuAs    []string `json:"menuAncoras"`
	BannerURL string   `json:"banner"`
	MID       string   `json:"mid"`
	Request   *http.Request
}

type configsUpdateResponse struct {
	MID string `json:"mid"`
}

func decodeConfigsUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	idConv, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	dto := new(configsUpdateRequest)
	docoder := json.NewDecoder(r.Body)
	err = docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.ID = uint(idConv)
	dto.Request = r
	return dto, nil
}

func makeConfigsUpdateEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*configsUpdateRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1009, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1010, err, req.MID)
		}

		service := service.NewConfigsService(userToken.UserID, userToken.Kind)
		err = service.Update(int(req.ID), req.Collors, req.Links, req.MenuAs, req.BannerURL)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1011, err, req.MID)
		}

		return &configsUpdateResponse{
			MID: req.MID,
		}, nil
	}
}

func ConfigsUpdateHandler() http.Handler {
	return httptransport.NewServer(
		makeConfigsUpdateEndPoint(),
		decodeConfigsUpdateRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type configsRemoveRequest struct {
	ID      uint
	MID     string
	Request *http.Request
}

type configsRemoveResponse struct {
	MID string `json:"mid"`
}

func decodeConfigsRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	idConv, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	mid := r.URL.Query().Get("mid")
	dto := new(configsRemoveRequest)
	dto.ID = uint(idConv)
	dto.MID = mid
	dto.Request = r
	return dto, nil
}

func makeConfigsRemoveEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*configsRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1012, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1013, err, req.MID)
		}

		service := service.NewConfigsService(userToken.UserID, userToken.Kind)
		err = service.Remove(int(req.ID))
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1014, err, req.MID)
		}

		return &configsRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func ConfigsRemoveHandler() http.Handler {
	return httptransport.NewServer(
		makeConfigsRemoveEndPoint(),
		decodeConfigsRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
