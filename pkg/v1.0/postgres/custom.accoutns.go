package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/demo/pkg/v1.0/utils/converter"
)

const (
	CUSTOM_ACCOUNTS = `custom.accounts`
)

type CustomAccounts struct {
	UserId  string         `db:"user_id,omitempty"`
	Pass    string         `db:"pass,omitempty"`
	DelFlag bool           `db:"del_flag,omitempty"`
	Desc    sql.NullString `db:"desc,omitempty"`
	CreId   string         `db:"cre_id,omitempty"`
	CreTime time.Time      `db:"cre_time,omitempty"`
	ModId   string         `db:"mod_id,omitempty"`
	ModTime time.Time      `db:"mod_time,omitempty"`
}

// func select action history data from postgres
func (pgc *PgConn) CustomAccountsSelect(filter *CustomAccounts) ([]*CustomAccounts, error) {
	// create query syntax
	qsyntax := fmt.Sprintf(`SELECT * FROM %s WHERE`, CUSTOM_ACCOUNTS)

	// get filter struct metadata
	fil := reflect.ValueOf(*filter)

	for i := 0; i < fil.NumField(); i++ {
		if fil.Field(i).Type().String() == "bool" {
			qsyntax = fmt.Sprintf(`%s %s = '%v' AND`, qsyntax, converter.CamelToSnakeCase(fil.Type().Field(i).Name), fil.Field(i))
		} else {
			if !fil.Field(i).IsZero() {
				qsyntax = fmt.Sprintf(`%s %s = '%v' AND`, qsyntax, converter.CamelToSnakeCase(fil.Type().Field(i).Name), fil.Field(i))
			}
		}
	}
	qsyntax = fmt.Sprintf(`%s;`, strings.TrimRight(qsyntax, "AND "))

	// create query syntax
	db, err := sql.Open(POSTGRES, string(*pgc))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(qsyntax)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var scans []*CustomAccounts

	//Fetch data to struct
	for rows.Next() {
		var scan CustomAccounts

		err := rows.Scan(&scan.UserId, &scan.Pass, &scan.DelFlag, &scan.Desc,
			&scan.CreId, &scan.CreTime, &scan.ModId, &scan.ModTime)

		if err != nil {
			return nil, err
		}

		// Append for every next row
		scans = append(scans, &scan)
	}

	if len(scans) == 0 {
		return nil, errors.New("no data found")
	}

	return scans, nil
}

//func to insert custom.accounts
func (pgc *PgConn) CustomAccountsInsert(values *CustomAccounts) (*ResInsertInfo, error) {
	// create query syntax
	qsyntax := fmt.Sprintf(`INSERT INTO %s (`, CUSTOM_ACCOUNTS)

	// get parameter data struct metadata
	param := reflect.ValueOf(*values)
	for i := 0; i < param.NumField(); i++ {
		if param.Field(i).Type().String() == "bool" {
			if param.Field(i).Bool() {
				qsyntax = fmt.Sprintf(`%s %s,`, qsyntax, converter.CamelToSnakeCase(param.Type().Field(i).Name))
			}
		} else {
			if !param.Field(i).IsZero() {
				qsyntax = fmt.Sprintf(`%s %s,`, qsyntax, converter.CamelToSnakeCase(param.Type().Field(i).Name))
			}
		}
	}

	qsyntax = fmt.Sprintf(`%s) VALUES (`, strings.TrimSuffix(qsyntax, ","))

	// get values struct metadata
	val := reflect.ValueOf(*values)

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Type().String() == "bool" {
			if val.Field(i).Bool() {
				qsyntax = fmt.Sprintf(`%s '%v',`, qsyntax, val.Field(i))
			}
		} else {
			if !val.Field(i).IsZero() {
				qsyntax = fmt.Sprintf(`%s '%v',`, qsyntax, val.Field(i))
			}
		}
	}

	qsyntax = fmt.Sprintf(`%s);`, strings.TrimSuffix(qsyntax, ","))

	// execute query
	db, err := sql.Open(POSTGRES, string(*pgc))
	if err != nil {
		return nil, err
	}

	res, err := db.Exec(qsyntax)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	lastInsertetId, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()

	return &ResInsertInfo{
		LastInsertId: lastInsertetId,
		RowsAffected: rowsAffected,
	}, nil
}

// func to update table custom.account
func (pgc *PgConn) CustomAccountsUpdate(setdata, filter *CustomAccounts) (*ResInsertInfo, error) {
	// create query syntax
	qsyntax := fmt.Sprintf(`UPDATE %s SET`, CUSTOM_ACCOUNTS)

	// get set data struct metadata
	sd := reflect.ValueOf(*setdata)

	for i := 0; i < sd.NumField(); i++ {
		if sd.Field(i).Type().String() == "bool" {
			qsyntax = fmt.Sprintf(`%s %s = %v,`, qsyntax, converter.CamelToSnakeCase(sd.Type().Field(i).Name), sd.Field(i))
		} else {
			if !sd.Field(i).IsZero() {
				qsyntax = fmt.Sprintf(`%s %s = '%v',`, qsyntax, converter.CamelToSnakeCase(sd.Type().Field(i).Name), sd.Field(i))
			}
		}
	}

	qsyntax = fmt.Sprintf(`%s WHERE`, strings.TrimSuffix(qsyntax, ","))

	// get filter struct metadata
	fil := reflect.ValueOf(*filter)

	for i := 0; i < fil.NumField(); i++ {
		if fil.Field(i).Type().String() == "bool" {
			qsyntax = fmt.Sprintf(`%s %s = %v AND`, qsyntax, converter.CamelToSnakeCase(fil.Type().Field(i).Name), fil.Field(i))
		} else {
			if !fil.Field(i).IsZero() {
				qsyntax = fmt.Sprintf(`%s %s = '%v' AND`, qsyntax, converter.CamelToSnakeCase(fil.Type().Field(i).Name), fil.Field(i))
			}
		}
	}

	qsyntax = fmt.Sprintf(`%s;`, strings.TrimSuffix(qsyntax, "AND"))

	// execute query
	db, err := sql.Open(POSTGRES, string(*pgc))
	if err != nil {
		return nil, err
	}

	res, err := db.Exec(qsyntax)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	lastInsertetId, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		return nil, errors.New("no data found")
	}

	return &ResInsertInfo{
		LastInsertId: lastInsertetId,
		RowsAffected: rowsAffected,
	}, nil
}
