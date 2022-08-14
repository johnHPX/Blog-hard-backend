package repository

import (
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils"
)

type postRepositoryInterface interface {
	Store(post *model.Post) error
}

type postRepositoryImpl struct{}

func (r *postRepositoryImpl) Store(post *model.Post) error {
	db, err := utils.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		insert into tb_post 
		(id, title, content)
		values
		($1,$2,$3)
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(post.PostID, post.Title, post.Content)

	if err != nil {
		return err
	}
	defer stmt.Close()

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("error when registering")
	}

	return nil
}

func NewPostRepository() postRepositoryInterface {
	return &postRepositoryImpl{}
}
