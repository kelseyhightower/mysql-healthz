package healthz

import (
	"database/sql"
	"fmt"

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

func (dc *DBChecker) Reachable() error {
	err := dc.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

var tableExistQuery = `SELECT TABLE_NAME
FROM information_schema.tables
WHERE table_schema = ?
  AND table_name = ?
LIMIT 1;`

func (dc *DBChecker) TableExist(database, table string) error {
	var output interface{}
	err := dc.db.QueryRow(tableExistQuery, database, table).Scan(&output)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("table %s does not exist", table)
	case err != nil:
		return err
	default:
		return nil
	}
}
