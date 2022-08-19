package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/repository"
	"github.com/johnHPX/validator-hard/pkg/validator"
)

type postServieInterface interface {
	Store(title, content string) error
	List(offset, limit, page int) ([]model.Post, error)
	Count() (int, error)
	Find(id string) (*model.Post, error)
	ListTitle(title string, offeset, limit, page int) ([]model.Post, error)
	CountTitle(title string) (int, error)
	ListByCategory(categoryName string, offset, limit, page int) ([]model.Post, int, error)
	Update(id, title, content string) error
	Remove(id string) error
}

type postServiceImpl struct {
	UserID string
	Kind   string
}

func (s *postServiceImpl) Store(title, content string) error {

	if s.Kind != "adm" {
		return errors.New("apenas usuario admin pode utilizar essa função")
	}

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

func (s *postServiceImpl) List(offset, limit, page int) ([]model.Post, error) {

	repPost := repository.NewPostRepository()
	posts, err := repPost.List(offset, limit, page)
	if err != nil {
		return nil, err
	}

	repNumberLikes := repository.NewNumberLikerRepository()

	entities := make([]model.Post, 0)
	for _, v := range posts {
		countLikes, err := repNumberLikes.CountLikes(v.PostID)
		if err != nil {
			return nil, err
		}

		entities = append(entities, model.Post{
			PostID:  v.PostID,
			Title:   v.Title,
			Content: v.Content,
			Likes:   countLikes,
		})
	}

	return entities, nil
}

func (s *postServiceImpl) Count() (int, error) {

	repPost := repository.NewPostRepository()
	count, err := repPost.Count()
	if err != nil {
		return 0, err
	}

	return count, nil

}

func (s *postServiceImpl) Find(id string) (*model.Post, error) {
	if s.Kind != "adm" {
		return nil, errors.New("apenas usuario admin pode utilizar essa função")
	}

	val := validator.NewValidator()
	IdVal, err := val.CheckAnyData("id", 255, id, true)
	if err != nil {
		return nil, err
	}

	repPost := repository.NewPostRepository()
	post, err := repPost.Find(IdVal.(string))
	if err != nil {
		return nil, err
	}

	repNumberLikes := repository.NewNumberLikerRepository()
	countLikes, err := repNumberLikes.CountLikes(post.PostID)
	if err != nil {
		return nil, err
	}

	post.Likes = countLikes

	return post, nil
}

func (s *postServiceImpl) ListTitle(title string, offset, limit, page int) ([]model.Post, error) {
	val := validator.NewValidator()
	TitleVal, err := val.CheckAnyData("titulo", 255, title, true)
	if err != nil {
		return nil, err
	}

	repPost := repository.NewPostRepository()
	posts, err := repPost.ListTitle(TitleVal.(string), offset, limit, page)
	if err != nil {
		return nil, err
	}

	repNumberLikes := repository.NewNumberLikerRepository()

	entities := make([]model.Post, 0)
	for _, v := range posts {
		countLikes, err := repNumberLikes.CountLikes(v.PostID)
		if err != nil {
			return nil, err
		}

		entities = append(entities, model.Post{
			PostID:  v.PostID,
			Title:   v.Title,
			Content: v.Content,
			Likes:   countLikes,
		})
	}

	return entities, nil
}

func (s *postServiceImpl) CountTitle(title string) (int, error) {
	val := validator.NewValidator()
	TitleVal, err := val.CheckAnyData("titulo", 255, title, true)
	if err != nil {
		return 0, err
	}

	repPost := repository.NewPostRepository()
	count, err := repPost.CountTitle(TitleVal.(string))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *postServiceImpl) ListByCategory(categoryName string, offset, limit, page int) ([]model.Post, int, error) {
	val := validator.NewValidator()
	categoryVal, err := val.CheckAnyData("titulo", 255, categoryName, true)
	if err != nil {
		return nil, 0, err
	}

	repPost := repository.NewPostRepository()
	posts, err := repPost.ListCategory(categoryVal.(string), offset, limit, page)
	if err != nil {
		return nil, 0, err
	}

	repNumberLikes := repository.NewNumberLikerRepository()
	entities := make([]model.Post, 0)
	for _, v := range posts {
		countLikes, err := repNumberLikes.CountLikes(v.PostID)
		if err != nil {
			return nil, 0, err
		}

		entities = append(entities, model.Post{
			PostID:  v.PostID,
			Title:   v.Title,
			Content: v.Content,
			Likes:   countLikes,
		})
	}

	count, err := repPost.CountCategory(categoryVal.(string))
	if err != nil {
		return nil, 0, err
	}

	return entities, count, nil
}

func (s *postServiceImpl) Update(id, title, content string) error {
	if s.Kind != "adm" {
		return errors.New("apenas usuario admin pode utilizar essa função")
	}

	val := validator.NewValidator()
	IdVal, err := val.CheckAnyData("id", 36, id, true)
	if err != nil {
		return err
	}
	TitleVal, err := val.CheckAnyData("titulo", 255, title, true)
	if err != nil {
		return err
	}
	ContentVal, err := val.CheckAnyData("conteudo", 9999, content, true)
	if err != nil {
		return err
	}

	post := new(model.Post)
	post.PostID = IdVal.(string)
	post.Title = TitleVal.(string)
	post.Content = ContentVal.(string)

	repPost := repository.NewPostRepository()
	err = repPost.Update(post)
	if err != nil {
		return err
	}

	return nil
}

func (s *postServiceImpl) Remove(id string) error {
	if s.Kind != "adm" {
		return errors.New("apenas usuario admin pode utilizar essa função")
	}
	val := validator.NewValidator()
	IdVal, err := val.CheckAnyData("id", 36, id, true)
	if err != nil {
		return err
	}

	repPost := repository.NewPostRepository()
	err = repPost.Remove(IdVal.(string))
	if err != nil {
		return err
	}

	return nil
}

func NewPostService(userID, kind string) postServieInterface {
	return &postServiceImpl{
		UserID: userID,
		Kind:   kind,
	}
}
