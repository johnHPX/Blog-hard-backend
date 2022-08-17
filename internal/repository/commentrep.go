package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
)

type commentRepositoryInterface interface {
	Store(entity *model.Comment) error
	List(postID string, offset, limit, page int) ([]model.Comment, error)
	Count(postID string) (int, error)
	ListUser(userID string, offset, limit, page int) ([]model.Comment, error)
	CountUser(userID string) (int, error)
	ListUserPost(postID, userID string, offset, limit, page int) ([]model.Comment, error)
	CountUserPost(postID, userID string) (int, error)
	Find(commentID string) (*model.Comment, error)
	Update(entity *model.Comment) error
	Remove(commentID string) error
}

type commentRepositoryImpl struct{}

func (r *commentRepositoryImpl) scanIterator(rows *sql.Rows) (*model.Comment, error) {
	commentID := sql.NullString{}
	title := sql.NullString{}
	content := sql.NullString{}
	userID := sql.NullString{}
	postID := sql.NullString{}

	err := rows.Scan(
		&commentID,
		&title,
		&content,
		&userID,
		&postID,
	)

	if err != nil {
		return nil, err
	}

	comment := new(model.Comment)

	if commentID.Valid {
		comment.CommentID = commentID.String
	}

	if title.Valid {
		comment.Title = title.String
	}

	if content.Valid {
		comment.Content = content.String
	}

	if userID.Valid {
		comment.UserID = userID.String
	}

	if postID.Valid {
		comment.PostID = postID.String
	}

	return comment, nil
}

func (r *commentRepositoryImpl) Store(entity *model.Comment) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		insert into tb_comment
		(id, title, content, user_uid, post_pid)
		values
		($1,$2,$3, $4, $5)
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(entity.CommentID, entity.Title, entity.Content, entity.UserID, entity.PostID)
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

func (r *commentRepositoryImpl) List(postID string, offset, limit, page int) ([]model.Comment, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content,
		user_uid,
		post_pid
	FROM tb_comment
	WHERE deleted_at is null and post_pid = $1
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		comment, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, *comment)
	}

	return comments, nil
}

func (r *commentRepositoryImpl) Count(postID string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_comment
		WHERE deleted_at is null and post_pid = $1
	`

	var count int
	row := db.QueryRow(sqlText, postID)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *commentRepositoryImpl) ListUser(userID string, offset, limit, page int) ([]model.Comment, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content,
		user_uid,
		post_pid
	FROM tb_comment
	WHERE deleted_at is null and user_uid = $1
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		comment, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, *comment)
	}

	return comments, nil
}

func (r *commentRepositoryImpl) CountUser(userID string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_comment
		WHERE deleted_at is null and user_uid = $1
	`

	var count int
	row := db.QueryRow(sqlText, userID)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *commentRepositoryImpl) ListUserPost(postID, userID string, offset, limit, page int) ([]model.Comment, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content,
		user_uid,
		post_pid
	FROM tb_comment
	WHERE deleted_at is null and post_pid = $1 and user_uid = $2
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText, postID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		comment, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, *comment)
	}

	return comments, nil
}

func (r *commentRepositoryImpl) CountUserPost(postID, userID string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_comment
		WHERE deleted_at is null and post_pid = $1 and user_uid = $2
	`

	var count int
	row := db.QueryRow(sqlText, postID, userID)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *commentRepositoryImpl) Find(commentID string) (*model.Comment, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			id, 
			title,
			content,
			user_uid, 
			post_pid
		FROM tb_comment
		WHERE deleted_at is null and id = $1
	`

	rows, err := db.Query(sqlText, commentID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		comment, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}

		return comment, nil
	}

	return nil, errors.New("error finding")
}

func (r *commentRepositoryImpl) Update(entity *model.Comment) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_comment SET
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

	result, err := stmt.Exec(entity.CommentID, entity.Title, entity.Content)
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

func (r *commentRepositoryImpl) Remove(commentID string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_comment SET
			deleted_at = now()
		WHERE deleted_at is null and id = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(commentID)
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

func NewCommentRepository() commentRepositoryInterface {
	return &commentRepositoryImpl{}
}
