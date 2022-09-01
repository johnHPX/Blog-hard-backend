package repository

import (
	"database/sql"
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/domain/models"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/messages"
)

type codeRecoveryInterface interface {
	Store(entity *models.CodeRecovery) error
	Find(code string) (*models.CodeRecovery, error)
	Remove(code string) error
}

type codeRecoveryImpl struct{}

func (r *codeRecoveryImpl) Store(entity *models.CodeRecovery) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `

		INSERT INTO tb_code_recovery
		(code, user_uid, expired_at)
		VALUES
		($1, $2, $3)
	
	`
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(entity.Code, entity.UserID, entity.ExpiredAt)
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
func (r *codeRecoveryImpl) Find(code string) (*models.CodeRecovery, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			code,
			user_uid,
			expired_at
		FROM tb_code_recovery
		WHERE code = $1
	`

	codeE := sql.NullString{}
	userIDE := sql.NullString{}
	expiredAT := sql.NullTime{}

	row := db.QueryRow(sqlText, code)
	err = row.Scan(
		&codeE,
		&userIDE,
		&expiredAT,
	)
	if err != nil {
		return nil, err
	}

	codeRecovery := new(models.CodeRecovery)

	if codeE.Valid {
		codeRecovery.Code = codeE.String
	}

	if userIDE.Valid {
		codeRecovery.UserID = userIDE.String
	}

	if expiredAT.Valid {
		codeRecovery.ExpiredAt = expiredAT.Time
	}

	return codeRecovery, nil

}

func (r *codeRecoveryImpl) Remove(code string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
		DELETE FROM tb_code_recovery
		WHERE code = $1;
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(code)
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

func NewCodeRecoveryRepository() codeRecoveryInterface {
	return &codeRecoveryImpl{}
}
