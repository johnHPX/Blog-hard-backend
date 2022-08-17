package repository

import (
	"database/sql"
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
)

type numberLikesRepositoryInterface interface {
	Store(entity *model.NumberLikes) error
	CountLikes(postID string) (int, error)
	Find(postID, userID string) (*model.NumberLikes, error)
	Remove(id string) error
}

type numberLikesRepositoryImpl struct{}

func (r *numberLikesRepositoryImpl) scanIterator(rows *sql.Rows) (*model.NumberLikes, error) {
	numberLikesID := sql.NullString{}
	userID := sql.NullString{}
	postID := sql.NullString{}

	err := rows.Scan(
		&numberLikesID,
		&userID,
		&postID,
	)

	if err != nil {
		return nil, err
	}

	numberLikesEntity := new(model.NumberLikes)

	if numberLikesID.Valid {
		numberLikesEntity.NumberLikesID = numberLikesID.String
	}

	if userID.Valid {
		numberLikesEntity.UserId = userID.String
	}

	if postID.Valid {
		numberLikesEntity.PostId = postID.String
	}

	return numberLikesEntity, nil
}

func (r *numberLikesRepositoryImpl) Store(entity *model.NumberLikes) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	
		INSERT INTO tb_number_likes 
		(id, post_pid, user_uid)
		VALUES
		($1, $2, $3)
	
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
		return errors.New("error when register")
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
			post_pid = $1
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

	return 0, errors.New("erro when count")

}

func (r *numberLikesRepositoryImpl) Find(postID, userID string) (*model.NumberLikes, error) {
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

	return nil, errors.New("erro when finding")
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
		return errors.New("error when deleting")
	}

	return nil
}

func NewNumberLikerRepository() numberLikesRepositoryInterface {
	return &numberLikesRepositoryImpl{}
}