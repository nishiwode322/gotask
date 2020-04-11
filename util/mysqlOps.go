package util

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlOps struct {
	userName   string
	passWord   string
	ip         string
	port       string
	dataBase   string
	dbInstance *sql.DB
}

func (m MysqlOps) Open() error {
	sourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s", m.userName, m.passWord, m.ip, m.port, m.dataBase)
	db, err := sql.Open("mysql", sourceName)
	m.dbInstance = db
	return err
}

func (m MysqlOps) Close() error {
	if m.dbInstance == nil {
		return nil
	}
	return m.dbInstance.Close()
}

func (m MysqlOps) Exec(s string) (sql.Result, error) {
	if m.dbInstance == nil {
		return nil, errors.New("mysqlops.dbInstance is nil")
	}
	return m.dbInstance.Exec(s)
}

func (m MysqlOps) Ping() error {
	if m.dbInstance == nil {
		return errors.New("mysqlops.dbInstance is nil")
	}
	return m.dbInstance.Ping()
}

func (m MysqlOps) Query(s string) (*sql.Rows, error) {
	if m.dbInstance == nil {
		return nil, errors.New("mysqlops.dbInstance is nil")
	}
	return m.dbInstance.Query(s)
}

func (m MysqlOps) QueryRow(s string) *sql.Row {
	if m.dbInstance == nil {
		return nil
	}
	return m.dbInstance.QueryRow(s)
}

func (m MysqlOps) Prepare(s string) (*sql.Stmt, error) {
	if m.dbInstance == nil {
		return nil, errors.New("mysqlops.dbInstance is nil")
	}
	return m.dbInstance.Prepare(s)
}

/*
*params
*table:
*fields:
*values:
*
*return value
*sql.result
*error
 */
func (m MysqlOps) BatchInserts(table string, fields []string, values [][]interface{}) (sql.Result, error) {
	if m.dbInstance == nil {
		return nil, errors.New("mysqlops.dbInstance is nil")
	}

	keyValueFormatStrings := []string{}
	for _, x := range values {
		var tmp []string
		for _, e := range x {
			if v, ok := e.(string); ok {
				tmp = append(tmp, fmt.Sprintf("'%s'", v))
			} else {
				tmp = append(tmp, fmt.Sprintf("%v", x))
			}
		}
		tmpString := fmt.Sprintf("(%s)", strings.Join(tmp, ","))
		keyValueFormatStrings = append(keyValueFormatStrings, tmpString)
	}

	valueSql := fmt.Sprintf("insert into %s (%s) values %s", table, strings.Join(fields, ","), strings.Join(keyValueFormatStrings, ","))
	return m.Exec(valueSql)
}

func RunMysqlSample() {
	db, _ := sql.Open("mysql", "root:123456@(127.0.0.1:3306)/provincecity")
	defer db.Close()
	err := db.Ping()
	if err != nil {
		fmt.Println("database connect error! ", err)
		return
	}

	sql := "insert into province values (0,'江苏')"
	result, _ := db.Exec(sql)
	n, _ := result.RowsAffected()

	fmt.Println("affected row number is:", n)

	rows, _ := db.Query("select * from province")
	var id, name string
	for rows.Next() {
		rows.Scan(&id, &name)
		fmt.Println(id, "--", name)
	}
}
