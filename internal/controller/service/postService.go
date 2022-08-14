package service

import (
	"github.com/google/uuid"
	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/repository"
	"github.com/johnHPX/validator-hard/pkg/validator"
)

type postServieInterface interface {
	Store(title, content string) error
}

type postServiceImpl struct{}

func (s *postServiceImpl) Store(title, content string) error {
	// validator
	val := validator.NewValidator()
	TitleVal, err := val.CheckAnyData("titulo", 255, title, true)
	if err != nil {
		return err
	}
	ContentVal, err := val.CheckAnyData("conteudo", 9999, content, true)
	if err != nil {
		return err
	}
	// generating id
	postID := uuid.New()

	// repository
	repPost := repository.NewPostRepository()

	// create post entity
	postEntity := new(model.Post)
	postEntity.PostID = postID.String()
	postEntity.Title = TitleVal.(string)
	postEntity.Content = ContentVal.(string)
	postEntity.Likes = 0

	// create post rep
	err = repPost.Store(postEntity)
	if err != nil {
		return err
	}

	return nil
}

func NewPostService() postServieInterface {
	return &postServiceImpl{}
}
