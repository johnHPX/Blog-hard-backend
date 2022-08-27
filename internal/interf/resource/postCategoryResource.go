package resource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/johnHPX/blog-hard-backend/internal/appl/service"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/responseAPI"
)

type postCategoryStoreRequest struct {
	PostID   string `json:"postId"`
	Category string `json:"categoryId"`
	MID      string `json:"mid"`
	Request  *http.Request
}

type postCategoryStoreResponse struct {
	MID string `json:"mid"`
}

func decodePostCategoryStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(postCategoryStoreRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makePostCategoryStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*postCategoryStoreRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		svcToken := service.NewAccessService()
		userToken, err := svcToken.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		svcpostCategory := service.NewPostCategoryService(userToken.UserID, userToken.Kind)
		err = svcpostCategory.StorePostCategory(req.PostID, req.Category)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &postCategoryStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func PostCategoryStoreHandle() http.Handler {
	return httptransport.NewServer(
		makePostCategoryStoreEndPoint(),
		decodePostCategoryStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type postCategoryRemoveRequest struct {
	PostID     string `json:"postId"`
	CategoryID string `json:"categoryId"`
	MID        string `json:"mid"`
	Request    *http.Request
}

type postCategoryRemoveResponse struct {
	MID string `json:"mid"`
}

func decodePostCategoryRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(postCategoryRemoveRequest)
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makePostCategoryRemoveEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*postCategoryRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		svcToken := service.NewAccessService()
		userToken, err := svcToken.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		svcpostCategory := service.NewPostCategoryService(userToken.UserID, userToken.Kind)
		err = svcpostCategory.RemovePostCategory(req.PostID, req.CategoryID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		return &postCategoryRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func PostCategoryRemoveHandle() http.Handler {
	return httptransport.NewServer(
		makePostCategoryRemoveEndPoint(),
		decodePostCategoryRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
