package service

import (
	"github.com/google/uuid"
	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/repository"
	"github.com/johnHPX/validator-hard/pkg/validator"
)

type responseCommentServiceInterface interface {
	Store(commentID, title, content string) error
	List(commentID string, offset, limit, page int) ([]model.ResponseComment, int, error)
	ListUser(offset, limit, page int) ([]model.ResponseComment, int, error)
	Update(responseCommentID, title, content string) error
	Remove(responseCommentID string) error
}
type responseCommentServiceImpl struct {
	userID string
	kind   string
}

func (s *responseCommentServiceImpl) Store(commentID, title, content string) error {
	val := validator.NewValidator()

	commentIDval, err := val.CheckAnyData("id do comentario", 36, commentID, true)
	if err != nil {
		return err
	}
	titleVal, err := val.CheckAnyData("titulo", 255, title, true)
	if err != nil {
		return err
	}
	contentVal, err := val.CheckAnyData("conteudo", 2024, content, true)
	if err != nil {
		return err
	}

	responseCommentID := uuid.New()
	responseCommentEntity := new(model.ResponseComment)
	responseCommentEntity.ResponseCommentID = responseCommentID.String()
	responseCommentEntity.Title = titleVal.(string)
	responseCommentEntity.Content = contentVal.(string)
	responseCommentEntity.CommentID = commentIDval.(string)
	responseCommentEntity.UserID = s.userID

	repResponseComment := repository.NewResponseCommmentRepository()
	err = repResponseComment.Store(responseCommentEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *responseCommentServiceImpl) List(commentID string, offset, limit, page int) ([]model.ResponseComment, int, error) {
	val := validator.NewValidator()
	commentIDval, err := val.CheckAnyData("id do comentario", 36, commentID, true)
	if err != nil {
		return nil, 0, err
	}

	repResponseComment := repository.NewResponseCommmentRepository()
	responseCommentsEntities, err := repResponseComment.List(commentIDval.(string), offset, limit, page)
	if err != nil {
		return nil, 0, err
	}
	count, err := repResponseComment.Count(commentIDval.(string))
	if err != nil {
		return nil, 0, err
	}

	return responseCommentsEntities, count, nil
}

func (s *responseCommentServiceImpl) ListUser(offset, limit, page int) ([]model.ResponseComment, int, error) {
	repResponseComment := repository.NewResponseCommmentRepository()
	commentsEntities, err := repResponseComment.ListUser(s.userID, offset, limit, page)
	if err != nil {
		return nil, 0, err
	}
	count, err := repResponseComment.CountUser(s.userID)
	if err != nil {
		return nil, 0, err
	}

	return commentsEntities, count, nil
}

func (s *responseCommentServiceImpl) Update(responseCommentID, title, content string) error {
	val := validator.NewValidator()
	responseCommentIDVal, err := val.CheckAnyData("id do comentario", 36, responseCommentID, true)
	if err != nil {
		return err
	}
	titleVal, err := val.CheckAnyData("titulo", 255, title, true)
	if err != nil {
		return err
	}
	contentVal, err := val.CheckAnyData("conteudo", 2024, content, true)
	if err != nil {
		return err
	}

	responseCommentEntity := new(model.ResponseComment)
	responseCommentEntity.ResponseCommentID = responseCommentIDVal.(string)
	responseCommentEntity.Title = titleVal.(string)
	responseCommentEntity.Content = contentVal.(string)

	repComment := repository.NewResponseCommmentRepository()
	err = repComment.Update(responseCommentEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *responseCommentServiceImpl) Remove(responseCommentID string) error {
	val := validator.NewValidator()
	responseCommentIDVal, err := val.CheckAnyData("id do comentario", 36, responseCommentID, true)
	if err != nil {
		return err
	}

	repComment := repository.NewResponseCommmentRepository()
	err = repComment.Remove(responseCommentIDVal.(string))
	if err != nil {
		return err
	}

	return nil
}

func NewResponseCommentService(userID, kind string) responseCommentServiceInterface {
	return &responseCommentServiceImpl{
		userID: userID,
		kind:   kind,
	}
}
