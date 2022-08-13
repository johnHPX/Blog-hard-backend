package repository

import (
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils"
)

type personRepositoryInterface interface {
	Store(entity *model.Person, userID string) error
}

type personRepositoryImpl struct{}

func (r *personRepositoryImpl) Store(entity *model.Person, userID string) error {
	db, err := utils.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `INSERT INTO tb_person 
		(id, user_uid, name, telephone)
		VALUES
		($1, $2, $3, $4)
	 `
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(entity.PersonID, userID, entity.Name, entity.Telephone)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New("error when registering")
	}

	return nil
}

func NewPersonRepository() personRepositoryInterface {
	return &personRepositoryImpl{}
}
