package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
)

type postRepositoryInterface interface {
	Store(post *model.Post) error
	List(offset, limit, page int) ([]model.Post, error)
	Count() (int, error)
	Find(id string) (*model.Post, error)
	ListTitle(title string, offset, limit, page int) ([]model.Post, error)
	CountTitle(title string) (int, error)
	ListCategory(category string, offset, limit, page int) ([]model.Post, error)
	CountCategory(category string) (int, error)
	Update(post *model.Post) error
	Remove(id string) error
}

type postRepositoryImpl struct{}

func (r *postRepositoryImpl) scanIterator(rows *sql.Rows) (*model.Post, error) {
	postId := sql.NullString{}
	title := sql.NullString{}
	content := sql.NullString{}

	err := rows.Scan(
		&postId,
		&title,
		&content,
	)

	if err != nil {
		return nil, err
	}

	post := new(model.Post)

	if postId.Valid {
		post.PostID = postId.String
	}

	if title.Valid {
		post.Title = title.String
	}

	if content.Valid {
		post.Content = content.String
	}

	return post, nil
}

func (r *postRepositoryImpl) Store(post *model.Post) error {
	db, err := databaseConn.Connect()
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

func (r *postRepositoryImpl) List(offset, limit, page int) ([]model.Post, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content
	FROM tb_post
	WHERE deleted_at is null
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		post, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *postRepositoryImpl) Count() (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_post
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

func (r *postRepositoryImpl) Find(id string) (*model.Post, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			id, 
			title,
			content
		FROM tb_post
		WHERE deleted_at is null and id = $1
	`

	rows, err := db.Query(sqlText, id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		post, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}

		return post, nil
	}

	return nil, errors.New("error finding")
}

func (r *postRepositoryImpl) ListTitle(title string, offset, limit, page int) ([]model.Post, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content
	FROM tb_post
	WHERE deleted_at is null and title like $1
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	t := "%" + title + "%"

	rows, err := db.Query(sqlText, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		post, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *postRepositoryImpl) CountTitle(title string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_post
		WHERE deleted_at is null and title like $1
	`

	t := "%" + title + "%"

	var count int
	row := db.QueryRow(sqlText, t)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *postRepositoryImpl) ListCategory(category string, offset, limit, page int) ([]model.Post, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		p.id,
		p.title,
		p.content
	FROM tb_post p
	INNER JOIN tb_post_category pc ON pc.post_pid = p.id
	INNER JOIN tb_category c ON c.id = pc.category_cid
	WHERE p.deleted_at is null and pc.deleted_at is null and c.deleted_at is null
	and c.name = $1
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		post, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *postRepositoryImpl) CountCategory(category string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
	SELECT 
		COUNT(*)
	FROM tb_post p
	INNER JOIN tb_post_category pc ON pc.post_pid = p.id
	INNER JOIN tb_category c ON c.id = pc.category_cid
	WHERE p.deleted_at is null and pc.deleted_at is null and c.deleted_at is null
	 and c.name = $1
	`

	var count int
	row := db.QueryRow(sqlText, category)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *postRepositoryImpl) Update(post *model.Post) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_post SET
			title = $2,
			content = $3,
			updated_at = now()
		WHERE deleted_at is null and id = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(post.PostID, post.Title, post.Content)
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

func (r *postRepositoryImpl) Remove(id string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_post SET
			deleted_at = now()
		WHERE deleted_at is null and id = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("err deleting")
	}

	return nil
}

func NewPostRepository() postRepositoryInterface {
	return &postRepositoryImpl{}
}
