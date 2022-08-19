package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
)

type categoryRepositoryInterface interface {
	Store(entity *model.Category) error
	List(offset, limit, page int) ([]model.Category, error)
	Count() (int, error)
	Find(categoryID string) (*model.Category, error)
	Update(entity *model.Category) error
	Remove(categoryID string) error
}

type categoryRepositoryImpl struct{}

func (r *categoryRepositoryImpl) scanIterator(rows *sql.Rows) (*model.Category, error) {
	categoryID := sql.NullString{}
	name := sql.NullString{}

	err := rows.Scan(
		&categoryID,
		&name,
	)

	if err != nil {
		return nil, err
	}

	categoryEntity := new(model.Category)

	if categoryID.Valid {
		categoryEntity.CategoryID = categoryID.String
	}

	if name.Valid {
		categoryEntity.Name = name.String
	}

	return categoryEntity, nil
}

func (r *categoryRepositoryImpl) Store(entity *model.Category) error {
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
		return errors.New("error when registering")
	}

	return nil
}

func (r *categoryRepositoryImpl) List(offset, limit, page int) ([]model.Category, error) {
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

	categorys := make([]model.Category, 0)
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

func (r *categoryRepositoryImpl) Find(categoryID string) (*model.Category, error) {
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

	return nil, errors.New("error finding")
}

func (r *categoryRepositoryImpl) Update(entity *model.Category) error {
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
		return errors.New("err updating")
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
		return errors.New("errpr deleting")
	}

	return nil
}

func NewCategoryRepository() categoryRepositoryInterface {
	return &categoryRepositoryImpl{}
}
