package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils"
)

type userRepositoryInterface interface {
	Store(entity *model.User) error
	List(offset, limit, page int) ([]model.User, error)
	Count() (int, error)
	ListName(name string, offset, limit, page int) ([]model.User, error)
	CountListName(name string) (int, error)
	Find(id string) (*model.User, error)
}

type userRepositoryImpl struct{}

func (r *userRepositoryImpl) scanIterator(rows *sql.Rows) (*model.User, error) {
	userID := sql.NullString{}
	personID := sql.NullString{}
	name := sql.NullString{}
	telephone := sql.NullString{}
	nick := sql.NullString{}
	email := sql.NullString{}
	kind := sql.NullString{}

	err := rows.Scan(
		&userID,
		&personID,
		&name,
		&telephone,
		&nick,
		&email,
		&kind,
	)

	if err != nil {
		return nil, err
	}

	userEntity := new(model.User)

	if userID.Valid {
		userEntity.UserID = userID.String
	}

	if personID.Valid {
		userEntity.PersonID = personID.String
	}

	if name.Valid {
		userEntity.Name = name.String
	}

	if telephone.Valid {
		userEntity.Telephone = telephone.String
	}

	if nick.Valid {
		userEntity.Nick = nick.String
	}

	if email.Valid {
		userEntity.Email = email.String
	}

	if kind.Valid {
		userEntity.Kind = kind.String
	}

	return userEntity, nil
}

func (r *userRepositoryImpl) Store(entity *model.User) error {
	db, err := utils.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `INSERT INTO tb_user 
		(id, nick, email, secret, kind)
		VALUES
		($1, $2, $3, $4, $5)
	 `
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(entity.UserID, entity.Nick, entity.Email, entity.Secret, entity.Kind)
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

func (r *userRepositoryImpl) List(offset, limit, page int) ([]model.User, error) {
	db, err := utils.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
			select
				u.id,
				p.id,
				p.name,
				p.telephone,
				u.nick,
				u.email,
				u.kind
			from tb_person p
			INNER JOIN tb_user u ON u.id = p.user_uid
			where p.deleted_at is null and u.deleted_at is null
			LIMIT %v OFFSET ((%v - 1) * (%v)) + %v`, limit, page, limit, offset)

	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]model.User, 0)
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *e)
	}

	return entities, nil
}

func (r *userRepositoryImpl) Count() (int, error) {
	db, err := utils.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT COUNT(user_uid) FROM tb_person WHERE deleted_at is null;
	`
	row, err := db.Query(sqlText)
	if err != nil {
		return 0, err
	}
	defer row.Close()

	var countNumber int
	if row.Next() {
		err := row.Scan(&countNumber)
		if err != nil {
			return 0, err
		}
	}

	return countNumber, nil
}

func (r *userRepositoryImpl) ListName(name string, offset, limit, page int) ([]model.User, error) {
	db, err := utils.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
			select
				u.id,
				p.id,
				p.name,
				p.telephone,
				u.nick,
				u.email,
				u.kind
			from tb_person p
			INNER JOIN tb_user u ON u.id = p.user_uid
			where p.deleted_at is null and u.deleted_at is null and p.name like $1
			LIMIT %v OFFSET ((%v - 1) * (%v)) + %v`, limit, page, limit, offset)

	v := "%" + name + "%"

	rows, err := db.Query(sqlText, v)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]model.User, 0)
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *e)
	}

	return entities, nil
}

func (r *userRepositoryImpl) CountListName(name string) (int, error) {
	db, err := utils.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT COUNT(user_uid) FROM tb_person WHERE deleted_at is null and name like $1;
	`

	v := "%" + name + "%"

	row, err := db.Query(sqlText, v)
	if err != nil {
		return 0, err
	}
	defer row.Close()

	var countNumber int
	if row.Next() {
		err := row.Scan(&countNumber)
		if err != nil {
			return 0, err
		}
	}

	return countNumber, nil
}

func (r *userRepositoryImpl) Find(id string) (*model.User, error) {
	db, err := utils.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		SELECT 
			u.id,
			p.id,
			p.name,
			p.telephone,
			u.nick,
			u.email,
			u.kind
		FROM tb_person p
		INNER JOIN tb_user u ON u.id = p.user_uid
		WHERE p.deleted_at is null and u.deleted_at is null
			 and (p.id = $1 or u.id = $1);
	`

	row, err := db.Query(sqlText, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		user, err := r.scanIterator(row)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, errors.New("error finding")
}

func NewUserRepository() userRepositoryInterface {
	return &userRepositoryImpl{}
}
