package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/johnHPX/blog-hard-backend/internal/domain/models"
	"github.com/johnHPX/blog-hard-backend/internal/infra/repository"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/messages"

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

	numberLikesEntity, err := repNumberLikes.Find(postID, s.userID)
	if err == nil {
		if numberLikesEntity.ValueLike {
			return errors.New(messages.LikePost)
		} else {
			err = repNumberLikes.Update(numberLikesEntity.NumberLikesID, true)
			if err != nil {
				return err
			}
			return nil
		}
	}

	nlid := uuid.New()
	entity := new(models.NumberLikes)
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
		return errors.New(messages.DeslikePost)
	}

	if entity.ValueLike {
		err = repNumberLikes.Update(entity.NumberLikesID, false)
		if err != nil {
			return err
		}
	} else {
		return errors.New("esse post já foi discutido por você!")
	}

	return nil
}

func NewNumberLikesService(userID string) numberLikesServiceInterface {
	return &numberLikesServiceImpl{
		userID: userID,
	}
}
