package resource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/johnHPX/blog-hard-backend/internal/controller/service"
	"github.com/johnHPX/blog-hard-backend/internal/utils/responseAPI"
)

type numberLikesStoreRequest struct {
	PostID  string `json:"postId"`
	MID     string `json:"mid"`
	Request *http.Request
}

type numberLikesStoreResponse struct {
	MID string `json:"mid"`
}

func decodeNumberLikesStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(numberLikesStoreRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makeNumberLikesStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*numberLikesStoreRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		svcToken := service.NewAccessService()
		userToken, err := svcToken.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		svcNumberLikes := service.NewNumberLikesService(userToken.UserID)
		err = svcNumberLikes.LikePost(req.PostID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &numberLikesStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func NumberLikesStoreHandle() http.Handler {
	return httptransport.NewServer(
		makeNumberLikesStoreEndPoint(),
		decodeNumberLikesStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type numberLikesRemoveRequest struct {
	PostID  string `json:"postId"`
	MID     string `json:"mid"`
	Request *http.Request
}

type numberLikesRemoveResponse struct {
	MID string `json:"mid"`
}

func decodeNumberLikesRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(numberLikesRemoveRequest)
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makeNumberLikesRemoveEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*numberLikesRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		svcToken := service.NewAccessService()
		userToken, err := svcToken.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1004, err, req.MID)
		}

		svcNumberLikes := service.NewNumberLikesService(userToken.UserID)
		err = svcNumberLikes.DislikePost(req.PostID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}

		return &numberLikesRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func NumberLikesRemoveHandle() http.Handler {
	return httptransport.NewServer(
		makeNumberLikesRemoveEndPoint(),
		decodeNumberLikesRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
