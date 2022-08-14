package resource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/johnHPX/blog-hard-backend/internal/controller/service"
	"github.com/johnHPX/blog-hard-backend/internal/utils"
)

type postStoreRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	MID     string `json:"mid"`
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
	return dto, nil
}

func makePostStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*postStoreRequest)
		if !ok {
			return nil, utils.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		service := service.NewPostService()
		err := service.Store(req.Title, req.Content)
		if err != nil {
			return nil, utils.CreateHttpErrorResponse(http.StatusInternalServerError, 1001, err, req.MID)
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
		utils.EncodeResponse,
		httptransport.ServerErrorEncoder(utils.ErrorEncoder()),
	)
}
