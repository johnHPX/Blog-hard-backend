package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/repository"
	"github.com/johnHPX/blog-hard-backend/internal/utils/messages"
	"github.com/johnHPX/validator-hard/pkg/validator"
)

type categoryServiceInterface interface {
	CreateCategory(name string) error
	ListCategory(offset, limit, page int) ([]model.Category, int, error)
	ListCategoryByPost(postID string, offset, limit, page int) ([]model.Category, int, error)
	FindCategory(categoryID string) (*model.Category, error)
	UpdateCategory(categoryID, name string) error
	RemoveCategory(categoryID string) error
}

type categoryServiceImpl struct {
	userID string
	kind   string
}

func (s *categoryServiceImpl) CreateCategory(name string) error {

	if s.kind != "adm" {
		return errors.New(messages.AdmMessage)
	}

	val := validator.NewValidator()
	nameVal, err := val.CheckAnyData("nome", 255, name, true)
	if err != nil {
		return err
	}

	categoryID := uuid.New()
	categoryEntity := new(model.Category)
	categoryEntity.CategoryID = categoryID.String()
	categoryEntity.Name = nameVal.(string)

	repCategory := repository.NewCategoryRepository()
	err = repCategory.Store(categoryEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryServiceImpl) ListCategory(offset, limit, page int) ([]model.Category, int, error) {

	if s.kind != "adm" {
		return nil, 0, errors.New(messages.AdmMessage)
	}

	repCategory := repository.NewCategoryRepository()
	categoryEntities, err := repCategory.List(offset, limit, page)
	if err != nil {
		return nil, 0, err
	}
	count, err := repCategory.Count()
	if err != nil {
		return nil, 0, err
	}

	return categoryEntities, count, nil
}

func (s *categoryServiceImpl) ListCategoryByPost(postID string, offset, limit, page int) ([]model.Category, int, error) {

	if s.kind != "adm" {
		return nil, 0, errors.New(messages.AdmMessage)
	}

	val := validator.NewValidator()
	postIDval, err := val.CheckAnyData("id da postagem", 36, postID, true)
	if err != nil {
		return nil, 0, err
	}

	repCategory := repository.NewCategoryRepository()
	categoryEntities, err := repCategory.ListPost(postIDval.(string), offset, limit, page)
	if err != nil {
		return nil, 0, err
	}
	count, err := repCategory.CountPost(postIDval.(string))
	if err != nil {
		return nil, 0, err
	}

	return categoryEntities, count, nil
}

func (s *categoryServiceImpl) FindCategory(categoryID string) (*model.Category, error) {

	if s.kind != "adm" {
		return nil, errors.New(messages.AdmMessage)
	}

	val := validator.NewValidator()
	categoryIDval, err := val.CheckAnyData("id da categoria", 36, categoryID, true)
	if err != nil {
		return nil, err
	}

	repCategory := repository.NewCategoryRepository()
	comment, err := repCategory.Find(categoryIDval.(string))
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *categoryServiceImpl) UpdateCategory(categoryID, name string) error {

	if s.kind != "adm" {
		return errors.New(messages.AdmMessage)
	}

	val := validator.NewValidator()
	categoryIDVal, err := val.CheckAnyData("id da categoria", 36, categoryID, true)
	if err != nil {
		return err
	}
	nameVal, err := val.CheckAnyData("nome da categoria", 255, name, true)
	if err != nil {
		return err
	}

	categoryEntity := new(model.Category)
	categoryEntity.CategoryID = categoryIDVal.(string)
	categoryEntity.Name = nameVal.(string)

	repCategory := repository.NewCategoryRepository()
	err = repCategory.Update(categoryEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryServiceImpl) RemoveCategory(categoryID string) error {

	if s.kind != "adm" {
		return errors.New(messages.AdmMessage)
	}

	val := validator.NewValidator()
	categoryIDVal, err := val.CheckAnyData("id da categoria", 36, categoryID, true)
	if err != nil {
		return err
	}

	repCategory := repository.NewCategoryRepository()
	err = repCategory.Remove(categoryIDVal.(string))
	if err != nil {
		return err
	}

	return nil
}

func NewCategoryService(userID, kind string) categoryServiceInterface {
	return &categoryServiceImpl{
		userID: userID,
		kind:   kind,
	}
}
