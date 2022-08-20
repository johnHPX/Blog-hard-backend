package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/model"
	"github.com/johnHPX/blog-hard-backend/internal/utils/databaseConn"
	"github.com/johnHPX/blog-hard-backend/internal/utils/messages"
	"github.com/lib/pq"
)

type configsRepositoryInterface interface {
	Store(configs *model.Configs) error
	List(offset, limit, page int) ([]model.Configs, error)
	Count() (int, error)
	Find(configsID int) (*model.Configs, error)
	Update(configs *model.Configs) error
	Remove(configsID int) error
}

type configsRepositoryImpl struct{}

func (r *configsRepositoryImpl) scanIterator(rows *sql.Rows) (*model.Configs, error) {
	configsID := sql.NullInt32{}
	collors := []sql.NullString{}
	links := []sql.NullString{}
	menuAs := []sql.NullString{}
	bannerURL := sql.NullString{}

	rows.Scan(
		&configsID,
		pq.Array(&collors),
		pq.Array(&links),
		pq.Array(&menuAs),
		&bannerURL,
	)

	configs := new(model.Configs)
	if configsID.Valid {
		configs.ConfigID = uint(configsID.Int32)
	}

	for _, v := range collors {
		if v.Valid {
			configs.Collors = append(configs.Collors, v.String)
		}
	}

	for _, v := range links {
		if v.Valid {
			configs.Links = append(configs.Links, v.String)
		}
	}

	for _, v := range menuAs {
		if v.Valid {
			configs.MenuAs = append(configs.MenuAs, v.String)
		}
	}

	if bannerURL.Valid {
		configs.BannerURL = bannerURL.String
	}

	return configs, nil

}

func (r *configsRepositoryImpl) Store(configs *model.Configs) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	
		INSERT INTO tb_configs
		(collors, links, menuAs, bannerURL)
		VALUES
		($1, $2, $3, $4)
	
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(pq.Array(configs.Collors), pq.Array(configs.Links), pq.Array(configs.MenuAs), configs.BannerURL)
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

func (r *configsRepositoryImpl) List(offset, limit, page int) ([]model.Configs, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := fmt.Sprintf(`
		SELECT
			id, collors, links, menuAs, bannerURL
		FROM tb_configs
		WHERE deleted_at is null
		LIMIT %v OFFSET ((%v - 1) * %v) + %v
	`, limit, page, limit, offset)

	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}

	configsEntities := make([]model.Configs, 0)
	for rows.Next() {
		config, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		configsEntities = append(configsEntities, *config)
	}

	return configsEntities, nil
}
func (r *configsRepositoryImpl) Count() (int, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlText := `
		SELECT
			COUNT(*)
		FROM tb_configs
		WHERE deleted_at is null
	`

	var count int
	row := db.QueryRow(sqlText)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *configsRepositoryImpl) Find(configsID int) (*model.Configs, error) {
	db, err := databaseConn.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlText := `
		SELECT
			id, collors, links, menuAs, bannerURL
		FROM tb_configs
		WHERE deleted_at is null and id = $1
	`

	rows, err := db.Query(sqlText, configsID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		config, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return nil, errors.New(messages.FindError)
}

func (r *configsRepositoryImpl) Update(configs *model.Configs) error {

	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	
		UPDATE tb_configs SET
			collors = $2,
			links = $3,
			menuAs = $4,
			bannerURL = $5,
			updated_at = NOW()
		WHERE deleted_at is null and id = $1
	
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(configs.ConfigID, pq.Array(configs.Collors), pq.Array(configs.Links), pq.Array(configs.MenuAs), configs.BannerURL)
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

func (r *configsRepositoryImpl) Remove(configsID int) error {
	db, err := databaseConn.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlText := `
	
		UPDATE tb_configs SET
			deleted_at = NOW()
		WHERE deleted_at is null and id = $1
	
	`

	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(configsID)
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

func NewConfigsRepository() configsRepositoryInterface {
	return &configsRepositoryImpl{}
}
