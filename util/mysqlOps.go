package util

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlOps struct {
	UserName   string
	PassWord   string
	Ip         string
	Port       string
	DataBase   string
	DbInstance *sql.DB
}

func (m MysqlOps) Open() error {
	sourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s", m.UserName, m.PassWord, m.Ip, m.Port, m.DataBase)
	db, err := sql.Open("mysql", sourceName)
	if err != nil {
		m.DbInstance = db
	}
	return err
}

func (m MysqlOps) Close() error {
	if m.DbInstance == nil {
		return nil
	}
	return m.DbInstance.Close()
}

func (m MysqlOps) Exec(s string) (sql.Result, error) {
	if m.DbInstance == nil {
		return nil, errors.New("mysqlops.DbInstance is nil")
	}
	return m.DbInstance.Exec(s)
}

func (m MysqlOps) Ping() error {
	if m.DbInstance == nil {
		return errors.New("mysqlops.DbInstance is nil")
	}
	return m.DbInstance.Ping()
}

func (m MysqlOps) Query(s string) (*sql.Rows, error) {
	if m.DbInstance == nil {
		return nil, errors.New("mysqlops.DbInstance is nil")
	}
	return m.DbInstance.Query(s)
}

func (m MysqlOps) QueryRow(s string) *sql.Row {
	if m.DbInstance == nil {
		return nil
	}
	return m.DbInstance.QueryRow(s)
}

func (m MysqlOps) Prepare(s string) (*sql.Stmt, error) {
	if m.DbInstance == nil {
		return nil, errors.New("mysqlops.DbInstance is nil")
	}
	return m.DbInstance.Prepare(s)
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
	if m.DbInstance == nil {
		return nil, errors.New("mysqlops.DbInstance is nil")
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
