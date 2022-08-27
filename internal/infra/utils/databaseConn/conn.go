package databaseConn

import (
	"database/sql"
	"fmt"

	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/configsAPI"
	_ "github.com/lib/pq"
)

// Connect: return database startement
func Connect() (*sql.DB, error) {
	config := configsAPI.NewConfigs()
	database, err := config.DatabaseConfigs()
	if err != nil {
		return nil, err
	}
	stringConnect := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", database.User, database.Pswd, database.Host, database.Port, database.Dbnm)
	db, err := sql.Open("postgres", stringConnect)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
