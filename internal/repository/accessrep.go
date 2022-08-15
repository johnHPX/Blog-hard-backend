package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
)

type accessRepositoryInterface interface {
	Store(rtoken, userID string, expired time.Time) error
	BlockAcess(userID string, block bool) error
	FindToken(userID string) (*model.Access, error)
	UpdateToken(rtoken, userID string) error
	RemoveToken(userID string) error
}

type accessRepositoryImpl struct{}

func (r *accessRepositoryImpl) scanIterator(rows *sql.Rows) (*model.Access, error) {
	token := sql.NullString{}
	userID := sql.NullString{}
	expired := sql.NullTime{}
	blocked := sql.NullBool{}

	err := rows.Scan(
		&token,
		&userID,
		&expired,
		&blocked,
	)

	if err != nil {
		return nil, err
	}

	access := new(model.Access)

	if token.Valid {
		access.Token = token.String
	}

	if userID.Valid {
		access.UserID = userID.String
	}

	if expired.Valid {
		access.ExpiredAt = expired.Time
	}

	if blocked.Valid {
		access.IsBlocked = blocked.Bool
	}

	return access, nil

}

func (r *accessRepositoryImpl) Store(token, userID string, expired time.Time) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		INSERT INTO tb_access 
		(token, user_uid, expired_at)
		VALUES
		($1, $2, $3)
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(token, userID, expired)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New("error when registering")
	}

	return nil
}

func (r *accessRepositoryImpl) BlockAcess(userID string, block bool) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_access SET
			is_blocked = $2
		WHERE deleted_at is null and user_uid = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID, block)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New("error when block")
	}

	return nil

}

func (r *accessRepositoryImpl) FindToken(userID string) (*model.Access, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		select
		 	token, user_uid, expired_at, is_blocked
		from tb_access
		where deleted_at is null and user_uid = $1
	`

	rows, err := db.Query(sqlText, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		access, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}

		return access, nil
	}

	return nil, errors.New("error finding")

}

func (r *accessRepositoryImpl) UpdateToken(rtoken, userID string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	UPDATE tb_access SET
		token = $2,
		updated_at = now()
	WHERE deleted_at is null and user_uid = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID, rtoken)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New("error when updating")
	}

	return nil

}

func (r *accessRepositoryImpl) RemoveToken(userID string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		UPDATE tb_access SET
			deleted_at = now()
		WHERE deleted_at is null and user_uid = $1
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New("error when deleting")
	}

	return nil
}

func NewAccessRepository() accessRepositoryInterface {
	return &accessRepositoryImpl{}
}
