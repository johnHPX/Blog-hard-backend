package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/domain/models"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/messages"
)

type responseCommentRepositoryInterface interface {
	Store(entity *models.ResponseComment) error
	List(commentID string, offset, limit, page int) ([]models.ResponseComment, error)
	Count(commentID string) (int, error)
	ListUser(userID string, offset, limit, page int) ([]models.ResponseComment, error)
	CountUser(userID string) (int, error)
	Update(entity *models.ResponseComment) error
	Remove(responseCommentID string) error
}

type responseCommentRepositoryImpl struct{}

func (r *responseCommentRepositoryImpl) scanIterator(rows *sql.Rows) (*models.ResponseComment, error) {
	ID := sql.NullString{}
	Title := sql.NullString{}
	Content := sql.NullString{}
	CommentID := sql.NullString{}
	UserID := sql.NullString{}

	err := rows.Scan(
		&ID,
		&Title,
		&Content,
		&CommentID,
		&UserID,
	)

	if err != nil {
		return nil, err
	}

	responseCommentEntity := new(models.ResponseComment)

	if ID.Valid {
		responseCommentEntity.ResponseCommentID = ID.String
	}

	if Title.Valid {
		responseCommentEntity.Title = Title.String
	}

	if Content.Valid {
		responseCommentEntity.Content = Content.String
	}

	if CommentID.Valid {
		responseCommentEntity.CommentID = CommentID.String
	}

	if UserID.Valid {
		responseCommentEntity.UserID = UserID.String
	}

	return responseCommentEntity, nil
}

func (r *responseCommentRepositoryImpl) Store(entity *models.ResponseComment) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		insert into tb_response_comment
		(id, title, content, comment_cid, user_uid)
		values
		($1,$2,	$3, $4, $5)
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(entity.ResponseCommentID, entity.Title, entity.Content, entity.CommentID, entity.UserID)
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
func (r *responseCommentRepositoryImpl) List(commentID string, offset, limit, page int) ([]models.ResponseComment, error) {
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
		comment_cid,
		user_uid
	FROM tb_response_comment
	WHERE deleted_at is null and comment_cid = $1
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	responseComments := make([]models.ResponseComment, 0)
	for rows.Next() {
		responseComment, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		responseComments = append(responseComments, *responseComment)
	}

	return responseComments, nil
}
func (r *responseCommentRepositoryImpl) Count(commentID string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_response_comment
		WHERE deleted_at is null and comment_cid = $1
	`

	var count int
	row := db.QueryRow(sqlText, commentID)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
func (r *responseCommentRepositoryImpl) ListUser(userID string, offset, limit, page int) ([]models.ResponseComment, error) {
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
		comment_cid,
		user_uid
	FROM tb_response_comment
	WHERE deleted_at is null and user_uid = $1
	LIMIT %v OFFSET ((%v - 1) * %v) + %v
`, limit, page, limit, offset)

	rows, err := db.Query(sqlText, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	responseComments := make([]models.ResponseComment, 0)
	for rows.Next() {
		responseComment, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		responseComments = append(responseComments, *responseComment)
	}

	return responseComments, nil
}
func (r *responseCommentRepositoryImpl) CountUser(userID string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(*)
		FROM tb_response_comment
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
func (r *responseCommentRepositoryImpl) Update(entity *models.ResponseComment) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_response_comment SET
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

	result, err := stmt.Exec(entity.ResponseCommentID, entity.Title, entity.Content)
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
func (r *responseCommentRepositoryImpl) Remove(responseCommentID string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_response_comment SET
			deleted_at = now()
		WHERE deleted_at is null and id = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(responseCommentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New(messages.RemoveError)
	}

	return nil
}

func NewResponseCommmentRepository() responseCommentRepositoryInterface {
	return &responseCommentRepositoryImpl{}
}
