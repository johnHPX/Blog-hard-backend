package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/repository"
	"github.com/johnHPX/validator-hard/pkg/validator"
)

type numberLikesServiceInterface interface {
	LikePost(postID string) error
	DislikePost(postID string) error
}

type numberLikesServiceImpl struct {
	userID string
}

func (s *numberLikesServiceImpl) LikePost(postID string) error {
	val := validator.NewValidator()
	PostIDVal, err := val.CheckAnyData("id", 36, postID, true)
	if err != nil {
		return err
	}

	repNumberLikes := repository.NewNumberLikerRepository()

	_, err = repNumberLikes.Find(postID, s.userID)
	if err == nil {
		return errors.New("apenas uma curtida por publicação")
	}

	nlid := uuid.New()
	entity := new(model.NumberLikes)
	entity.NumberLikesID = nlid.String()
	entity.PostId = PostIDVal.(string)
	entity.UserId = s.userID

	err = repNumberLikes.Store(entity)
	if err != nil {
		return err
	}

	return nil
}

func (s *numberLikesServiceImpl) DislikePost(postID string) error {
	val := validator.NewValidator()
	PostIDVal, err := val.CheckAnyData("id", 36, postID, true)
	if err != nil {
		return err
	}

	repNumberLikes := repository.NewNumberLikerRepository()
	entity, err := repNumberLikes.Find(PostIDVal.(string), s.userID)
	if err != nil {
		return err
	}

	err = repNumberLikes.Remove(entity.NumberLikesID)
	if err != nil {
		return err
	}

	return nil
}

func NewNumberLikesService(userID string) numberLikesServiceInterface {
	return &numberLikesServiceImpl{
		userID: userID,
	}
}
