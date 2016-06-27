package healthz

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DBChecker struct {
	db *sql.DB
}

func NewDBChecker(driverName, dataSourceName string) (*DBChecker, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DBChecker{db}, nil
}

func (dc *DBChecker) Ping() error {
	err := dc.db.Ping()
	if err != nil {
		return err
	}
	return nil
}
