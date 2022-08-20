package repository

import (
	"errors"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/utils/messages"
)

type personRepositoryInterface interface {
	Store(entity *model.Person, userID string) error
	Update(entity *model.Person) error
	Remove(id string) error
}

type personRepositoryImpl struct{}

func (r *personRepositoryImpl) Store(entity *model.Person, userID string) error {
	db, err := databaseConn.Connect()
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
		return errors.New(messages.StoreError)
	}

	return nil
}

func (r *personRepositoryImpl) Update(entity *model.Person) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	update tb_person set
		name = $2,
		telephone = $3,
		updated_at = now()
	where deleted_at is null and id = $1
	`

	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(entity.PersonID, entity.Name, entity.Telephone)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New(messages.UpdateError)
	}

	return nil
}

func (r *personRepositoryImpl) Remove(id string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_person set
		deleted_at = now()
	where id = $1
	`
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(id)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New(messages.RemoveError)
	}

	return nil
}

func NewPersonRepository() personRepositoryInterface {
	return &personRepositoryImpl{}
}
