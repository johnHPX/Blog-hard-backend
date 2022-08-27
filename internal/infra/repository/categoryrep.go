package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/domain/models"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/messages"
)

type categoryRepositoryInterface interface {
	Store(entity *models.Category) error
	List(offset, limit, page int) ([]models.Category, error)
	Count() (int, error)
	ListPost(postID string, offset, limit, page int) ([]models.Category, error)
	CountPost(postID string) (int, error)
	Find(categoryID string) (*models.Category, error)
	Update(entity *models.Category) error
	Remove(categoryID string) error
}

type categoryRepositoryImpl struct{}

func (r *categoryRepositoryImpl) scanIterator(rows *sql.Rows) (*models.Category, error) {
	categoryID := sql.NullString{}
	name := sql.NullString{}

	err := rows.Scan(
		&categoryID,
		&name,
	)

	if err != nil {
		return nil, err
	}

	categoryEntity := new(models.Category)

	if categoryID.Valid {
		categoryEntity.CategoryID = categoryID.String
	}

	if name.Valid {
		categoryEntity.Name = name.String
	}

	return categoryEntity, nil
}

func (r *categoryRepositoryImpl) Store(entity *models.Category) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		insert into tb_category
		(id, name)
		values
		($1,$2)
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(entity.CategoryID, entity.Name)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New(messages.StoreError)
	}

	return nil
}

func (r *categoryRepositoryImpl) List(offset, limit, page int) ([]models.Category, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		id,
		name
	FROM tb_category
	WHERE deleted_at is null
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categorys := make([]models.Category, 0)
	for rows.Next() {
		category, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		categorys = append(categorys, *category)
	}

	return categorys, nil
}

func (r *categoryRepositoryImpl) Count() (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_category
		WHERE deleted_at is null
	`

	var count int
	row := db.QueryRow(sqlText)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *categoryRepositoryImpl) ListPost(postID string, offset, limit, page int) ([]models.Category, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		c.id,
		c.name
	FROM tb_category c
	INNER JOIN tb_post_category pc ON pc.category_cid = c.id
	INNER JOIN tb_post p ON p.id = pc.post_pid
	WHERE p.deleted_at is null and pc.deleted_at is null and c.deleted_at is null
	 and p.id = $1
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categorys := make([]models.Category, 0)
	for rows.Next() {
		category, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		categorys = append(categorys, *category)
	}

	return categorys, nil
}

func (r *categoryRepositoryImpl) CountPost(postID string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_category c
		INNER JOIN tb_post_category pc ON pc.category_cid = c.id
		INNER JOIN tb_post p ON p.id = pc.post_pid
		WHERE p.deleted_at is null and pc.deleted_at is null and c.deleted_at is null
		and p.id = $1
	`

	var count int
	row := db.QueryRow(sqlText, postID)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *categoryRepositoryImpl) Find(categoryID string) (*models.Category, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			id, 
			name
		FROM tb_category
		WHERE deleted_at is null and id = $1
	`

	rows, err := db.Query(sqlText, categoryID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		category, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}

		return category, nil
	}

	return nil, errors.New(messages.FindError)
}

func (r *categoryRepositoryImpl) Update(entity *models.Category) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_category SET
			name = $2,
			updated_at = now()
		WHERE deleted_at is null and id = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(entity.CategoryID, entity.Name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New(messages.UpdateError)
	}

	return nil
}

func (r *categoryRepositoryImpl) Remove(categoryID string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_category SET
			deleted_at = now()
		WHERE deleted_at is null and id = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(categoryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New(messages.UpdateError)
	}

	return nil
}

func NewCategoryRepository() categoryRepositoryInterface {
	return &categoryRepositoryImpl{}
}
