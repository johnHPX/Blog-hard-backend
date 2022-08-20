package repository

import (
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/utils/messages"
)

type postCategoryRepositoryInterface interface {
	Store(entity *model.PostCategory) error
	Update(entity *model.PostCategory) error
	Remove(postID, categoryID string) error
}

type postCategoryRepositoryImpl struct{}

func (r *postCategoryRepositoryImpl) Store(entity *model.PostCategory) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `INSERT INTO tb_post_category 
		(id, post_pid, category_cid)
		VALUES
		($1, $2, $3)
	 `
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(entity.PostCategoryId, entity.PostId, entity.CategoryId)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New(messages.StoreError)
	}

	return nil
}

func (r *postCategoryRepositoryImpl) Update(entity *model.PostCategory) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	update tb_post_category set
		post_pid = $2,
		category_cid = $3,
		updated_at = now()
	where deleted_at is null and id = $1
	`

	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(entity.PostCategoryId, entity.PostId, entity.CategoryId)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New(messages.UpdateError)
	}

	return nil
}

func (r *postCategoryRepositoryImpl) Remove(postID, categoryID string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_post_category set
		deleted_at = now()
	where deleted_at is null and post_pid = $1 and category_cid = $2
	`
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(postID, categoryID)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New(messages.RemoveError)
	}

	return nil
}

func NewPostCategoryRepository() postCategoryRepositoryInterface {
	return &postCategoryRepositoryImpl{}
}
