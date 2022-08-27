package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/domain/models"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/messages"
)

type userRepositoryInterface interface {
	Store(entity *models.User) error
	List(offset, limit, page int) ([]models.User, error)
	Count() (int, error)
	ListName(name string, offset, limit, page int) ([]models.User, error)
	CountListName(name string) (int, error)
	Find(id string) (*models.User, error)
	Update(user *models.User) error
	Remove(id string) error
	CheckEmail(email string) error
	CheckNick(nick string) error
	FindByEmailOrNick(emailOrNick string) (*models.User, error)
}

type userRepositoryImpl struct{}

func (r *userRepositoryImpl) scanIterator(rows *sql.Rows, secretIsReq bool) (*models.User, error) {
	userID := sql.NullString{}
	personID := sql.NullString{}
	name := sql.NullString{}
	telephone := sql.NullString{}
	nick := sql.NullString{}
	email := sql.NullString{}
	kind := sql.NullString{}
	secret := sql.NullString{}

	var err error
	if secretIsReq {
		err = rows.Scan(
			&userID,
			&personID,
			&name,
			&telephone,
			&nick,
			&email,
			&kind,
			&secret,
		)
	} else {
		err = rows.Scan(
			&userID,
			&personID,
			&name,
			&telephone,
			&nick,
			&email,
			&kind,
		)
	}

	if err != nil {
		return nil, err
	}

	userEntity := new(models.User)

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

	if secretIsReq {
		if secret.Valid {
			userEntity.Secret = secret.String
		}
	}

	return userEntity, nil
}

func (r *userRepositoryImpl) CheckEmail(email string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	select
	 	email
	from tb_user
	where email = $1
	`

	row, err := db.Query(sqlText, email)
	if err != nil {
		return err
	}
	defer row.Close()

	if row.Next() {
		return errors.New(messages.EmailIsRegister)
	}

	return nil
}

func (r *userRepositoryImpl) CheckNick(nick string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	select
	 	nick
	from tb_user
	where nick = $1
	`

	row, err := db.Query(sqlText, nick)
	if err != nil {
		return err
	}
	defer row.Close()

	if row.Next() {
		return errors.New(messages.NickIsRegister)
	}

	return nil
}

func (r *userRepositoryImpl) Store(entity *models.User) error {
	db, err := databaseConn.Connect()
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
		return errors.New(messages.StoreError)
	}

	return nil
}

func (r *userRepositoryImpl) List(offset, limit, page int) ([]models.User, error) {
	db, err := databaseConn.Connect()
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

	entities := make([]models.User, 0)
	for rows.Next() {
		e, err := r.scanIterator(rows, false)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *e)
	}

	return entities, nil
}

func (r *userRepositoryImpl) Count() (int, error) {
	db, err := databaseConn.Connect()
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

func (r *userRepositoryImpl) ListName(name string, offset, limit, page int) ([]models.User, error) {
	db, err := databaseConn.Connect()
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

	entities := make([]models.User, 0)
	for rows.Next() {
		e, err := r.scanIterator(rows, false)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *e)
	}

	return entities, nil
}

func (r *userRepositoryImpl) CountListName(name string) (int, error) {
	db, err := databaseConn.Connect()
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

func (r *userRepositoryImpl) Find(id string) (*models.User, error) {
	db, err := databaseConn.Connect()
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
		user, err := r.scanIterator(row, false)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, errors.New(messages.FindError)
}

func (r *userRepositoryImpl) Update(user *models.User) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	update tb_user set
		nick = $2,
		email = $3,
		kind = $4,
		updated_at = now()
	where deleted_at is null and id = $1
	`

	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(user.UserID, user.Nick, user.Email, user.Kind)
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

func (r *userRepositoryImpl) Remove(id string) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_user set
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

func (r *userRepositoryImpl) FindByEmailOrNick(emailOrNick string) (*models.User, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
	select 
		u.id,
		p.id,
		p.name,
		p.telephone,
		u.nick,
		u.email,
		u.kind,
		u.secret
	FROM tb_person p
	INNER JOIN tb_user u ON u.id = p.user_uid
	WHERE p.deleted_at is null and u.deleted_at is null
		 and (email = $1 or nick = $1)
	`

	row, err := db.Query(sqlText, emailOrNick)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		user, err := r.scanIterator(row, true)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, errors.New(messages.UserNotExists)
}

func NewUserRepository() userRepositoryInterface {
	return &userRepositoryImpl{}
}
