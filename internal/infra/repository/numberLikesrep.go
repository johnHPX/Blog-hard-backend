package repository

import (
	"database/sql"
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/domain/models"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/messages"
)

type numberLikesRepositoryInterface interface {
	Store(entity *models.NumberLikes) error
	CountLikes(postID string) (int, error)
	Find(postID, userID string) (*models.NumberLikes, error)
	Update(id string, value bool) error
	Remove(id string) error
}

type numberLikesRepositoryImpl struct{}

func (r *numberLikesRepositoryImpl) scanIterator(rows *sql.Rows) (*models.NumberLikes, error) {
	numberLikesID := sql.NullString{}
	userID := sql.NullString{}
	postID := sql.NullString{}
	valueLike := sql.NullBool{}

	err := rows.Scan(
		&numberLikesID,
		&userID,
		&postID,
		&valueLike,
	)

	if err != nil {
		return nil, err
	}

	numberLikesEntity := new(models.NumberLikes)

	if numberLikesID.Valid {
		numberLikesEntity.NumberLikesID = numberLikesID.String
	}

	if userID.Valid {
		numberLikesEntity.UserId = userID.String
	}

	if postID.Valid {
		numberLikesEntity.PostId = postID.String
	}

	if valueLike.Valid {
		numberLikesEntity.ValueLike = valueLike.Bool
	}

	return numberLikesEntity, nil
}

func (r *numberLikesRepositoryImpl) Store(entity *models.NumberLikes) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	
		INSERT INTO tb_number_likes 
		(id, post_pid, user_uid, value_like)
		VALUES
		($1, $2, $3, true)
	
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(entity.NumberLikesID, entity.PostId, entity.UserId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New(messages.StoreError)
	}

	return nil
}

func (r *numberLikesRepositoryImpl) CountLikes(postID string) (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			COUNT(user_uid)
		FROM tb_number_likes 
		WHERE deleted_at is null and
			post_pid = $1 and value_like = true
	`

	row := db.QueryRow(sqlText, postID)
	count := sql.NullInt64{}
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	if count.Valid {
		return int(count.Int64), err
	}

	return 0, errors.New(messages.CountError)

}

func (r *numberLikesRepositoryImpl) Find(postID, userID string) (*models.NumberLikes, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			id,
			user_uid,
			post_pid
		FROM tb_number_likes 
		WHERE deleted_at is null and
			post_pid = $1 and user_uid = $2
	`

	rows, err := db.Query(sqlText, postID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		nl, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}

		return nl, nil
	}

	return nil, errors.New(messages.FindError)
}

func (r *numberLikesRepositoryImpl) Update(id string, value bool) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_number_likes SET
			value_like = $2
			updated_at = now()
		WHERE deleted_at is null and id = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(id, value)
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

func (r *numberLikesRepositoryImpl) Remove(id string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_number_likes SET
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
		return errors.New(messages.RemoveError)
	}

	return nil
}

func NewNumberLikerRepository() numberLikesRepositoryInterface {
	return &numberLikesRepositoryImpl{}
}
