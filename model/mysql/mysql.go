package mysql

import (
	"database/sql"
	"errors"
)

const (
	mysqlDioxideInsert = iota
)

var (
	errInvalidInsert = errors.New("errInvalidInsert")
)

var (
	dioxideSQLString = []string{
		// `INSERT INTO test (id) VALUES (?)`,
		"INSERT INTO Co2_test (设备id, 状态, 地区名, `二氧化碳浓度(ppm)`) VALUES (?, ?, ?, ?)",
	}
)

func InsertDioxide(db *sql.DB, dioxideDensity, status int, zoneName, deviceId string) error {
	result, err := db.Exec(dioxideSQLString[mysqlDioxideInsert], deviceId, status, zoneName, dioxideDensity)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidInsert
	}

	return nil
}
