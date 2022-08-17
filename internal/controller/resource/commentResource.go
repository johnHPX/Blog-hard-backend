package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/johnHPX/blog-hard-backend/internal/controller/service"
	"github.com/johnHPX/blog-hard-backend/internal/utils/responseAPI"
)

type commentEntity struct {
	CommentID string `json:"commentID"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserID    string `json:"userID"`
	PostID    string `json:"postID"`
}

type commentStoreRequest struct {
	PostID  string `json:"postID"`
	Title   string `json:"title"`
	Content string `json:"content"`
	MID     string `json:"mid"`
	Request *http.Request
}

type commentStoreResponse struct {
	MID string `json:"mid"`
}

func decodeCommentStoreRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(commentStoreRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.Request = r
	return dto, nil
}

func makeCommentStoreEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*commentStoreRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1001, err, req.MID)
		}

		service := service.NewCommentService(userToken.UserID, userToken.Kind)
		err = service.CreateComment(req.PostID, req.Title, req.Content)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1002, err, req.MID)
		}

		return &postStoreResponse{
			MID: req.MID,
		}, nil
	}
}

func CommentStoreHandler() http.Handler {
	return httptransport.NewServer(
		makeCommentStoreEndPoint(),
		decodeCommentStoreRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type commentListPostRequest struct {
	PostID string
	Offset int
	Limit  int
	Page   int
	MID    string
}

type commentListPostResponse struct {
	Count    int             `json:"count"`
	Comments []commentEntity `json:"comments"`
	MID      string          `json:"mid"`
}

func decodeCommentListPostRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
	dto := &commentListPostRequest{
		PostID: postID,
		Offset: int(offset),
		Limit:  int(limit),
		Page:   int(page),
		MID:    mid,
	}
	return dto, nil
}

func makeCommentListPostEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*commentListPostRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1003, errors.New("invalid request"), "na")
		}

		service := service.NewCommentService("", "")
		comments, count, err := service.ListCommentsPost(req.PostID, req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1004, err, req.MID)
		}

		fmt.Println(comments)

		var entities []commentEntity
		for _, v := range comments {
			entities = append(entities, commentEntity{
				CommentID: v.CommentID,
				Title:     v.Title,
				Content:   v.Content,
				UserID:    v.UserID,
				PostID:    v.PostID,
			})
		}

		return &commentListPostResponse{
			Count:    count,
			Comments: entities,
			MID:      req.MID,
		}, nil
	}
}

func CommentListPostHandler() http.Handler {
	return httptransport.NewServer(
		makeCommentListPostEndPoint(),
		decodeCommentListPostRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type commentListUserRequest struct {
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type commentListUserResponse struct {
	Count    int             `json:"count"`
	Comments []commentEntity `json:"comments"`
	MID      string          `json:"mid"`
}

func decodeCommentListUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {

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
	dto := &commentListUserRequest{
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeCommentListUserEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*commentListUserRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1005, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1006, err, req.MID)
		}

		service := service.NewCommentService(userToken.UserID, userToken.Kind)
		comments, count, err := service.ListCommentsUser(req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1007, err, req.MID)
		}

		var entities []commentEntity
		for _, v := range comments {
			entities = append(entities, commentEntity{
				CommentID: v.CommentID,
				Title:     v.Title,
				Content:   v.Content,
				UserID:    v.UserID,
				PostID:    v.PostID,
			})
		}

		return &commentListUserResponse{
			Count:    count,
			Comments: entities,
			MID:      req.MID,
		}, nil
	}
}

func CommentListUserHandler() http.Handler {
	return httptransport.NewServer(
		makeCommentListUserEndPoint(),
		decodeCommentListUserRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type commentListPostUserRequest struct {
	PostID  string
	Offset  int
	Limit   int
	Page    int
	MID     string
	Request *http.Request
}

type commentListPostUserResponse struct {
	Count    int             `json:"count"`
	Comments []commentEntity `json:"comments"`
	MID      string          `json:"mid"`
}

func decodecommentListPostUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
	dto := &commentListPostUserRequest{
		PostID:  postID,
		Offset:  int(offset),
		Limit:   int(limit),
		Page:    int(page),
		MID:     mid,
		Request: r,
	}
	return dto, nil
}

func makeCommentListPostUserEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*commentListPostUserRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1008, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1009, err, req.MID)
		}

		service := service.NewCommentService(userToken.UserID, userToken.Kind)
		comments, count, err := service.ListCommentsPostUser(req.PostID, req.Offset, req.Limit, req.Page)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1010, err, req.MID)
		}

		var entities []commentEntity
		for _, v := range comments {
			entities = append(entities, commentEntity{
				CommentID: v.CommentID,
				Title:     v.Title,
				Content:   v.Content,
				UserID:    v.UserID,
				PostID:    v.PostID,
			})
		}

		return &commentListPostUserResponse{
			Count:    count,
			Comments: entities,
			MID:      req.MID,
		}, nil
	}
}

func CommentListPostUserHandler() http.Handler {
	return httptransport.NewServer(
		makeCommentListPostUserEndPoint(),
		decodecommentListPostUserRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type commentFindRequest struct {
	commentID string
	MID       string
	Request   *http.Request
}

type commentFindResponse struct {
	commentEntity
	MID string `json:"mid"`
}

func decodeCommentFindRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := &commentFindRequest{
		commentID: id,
		MID:       mid,
		Request:   r,
	}
	return dto, nil
}

func makeCommentFindendPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*commentFindRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1011, errors.New("invalid request"), "na")
		}

		// gets token's informations
		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1012, err, req.MID)
		}

		service := service.NewCommentService(userToken.UserID, userToken.Kind)
		comment, err := service.FindComment(req.commentID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1013, err, req.MID)
		}

		return &commentFindResponse{
			commentEntity: commentEntity{
				CommentID: comment.CommentID,
				Title:     comment.Title,
				Content:   comment.CommentID,
				UserID:    comment.UserID,
				PostID:    comment.PostID,
			},
			MID: req.MID,
		}, nil
	}
}

func CommentFindHandler() http.Handler {
	return httptransport.NewServer(
		makeCommentFindendPoint(),
		decodeCommentFindRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type commentUpdateRequest struct {
	commentID string
	Title     string `json:"title"`
	Content   string `json:"content"`
	MID       string `json:"mid"`
	Request   *http.Request
}

type commentUpdateResponse struct {
	MID string `json:"mid"`
}

func decodeCommentUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	dto := new(commentUpdateRequest)
	docoder := json.NewDecoder(r.Body)
	err := docoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.commentID = id
	dto.Request = r
	return dto, nil
}

func makeCommentUpdateEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*commentUpdateRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1014, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1015, err, req.MID)
		}

		service := service.NewCommentService(userToken.UserID, userToken.Kind)
		err = service.UpdateComment(req.commentID, req.Title, req.Content)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1016, err, req.MID)
		}

		return &commentUpdateResponse{
			MID: req.MID,
		}, nil
	}
}

func CommentUpdateHandler() http.Handler {
	return httptransport.NewServer(
		makeCommentUpdateEndPoint(),
		decodeCommentUpdateRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}

type commentRemoveRequest struct {
	ID      string
	MID     string
	Request *http.Request
}

type commentRemoveResponse struct {
	MID string `json:"mid"`
}

func decodeCommentRemoveRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	mid := r.URL.Query().Get("mid")
	dto := new(commentRemoveRequest)
	dto.ID = id
	dto.MID = mid
	dto.Request = r
	return dto, nil
}

func makeCommentRemoveEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*commentRemoveRequest)
		if !ok {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusBadRequest, 1017, errors.New("invalid request"), "na")
		}

		tokenFunc := service.NewAccessService()
		userToken, err := tokenFunc.ExtractTokenInfo(req.Request)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusUnauthorized, 1018, err, req.MID)
		}

		service := service.NewCommentService(userToken.UserID, userToken.Kind)
		err = service.RemoveComment(req.ID)
		if err != nil {
			return nil, responseAPI.CreateHttpErrorResponse(http.StatusInternalServerError, 1019, err, req.MID)
		}

		return &postRemoveResponse{
			MID: req.MID,
		}, nil
	}
}

func CommentRemoveHandler() http.Handler {
	return httptransport.NewServer(
		makeCommentRemoveEndPoint(),
		decodeCommentRemoveRequest,
		responseAPI.EncodeResponse,
		httptransport.ServerErrorEncoder(responseAPI.ErrorEncoder()),
	)
}
