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
	IP         string
	Port       string
	DataBase   string
	DbInstance *sql.DB
}

func (m *MysqlOps) Open() error {
	sourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s", m.UserName, m.PassWord, m.IP, m.Port, m.DataBase)
	db, err := sql.Open("mysql", sourceName)
	if err == nil {
		m.DbInstance = db
	}
	return err
}

func (m *MysqlOps) Close() error {
	if m.DbInstance == nil {
		return nil
	}
	return m.DbInstance.Close()
}

func (m *MysqlOps) Exec(s string) (sql.Result, error) {
	if m.DbInstance == nil {
		return nil, errors.New("mysqlops.DbInstance is nil")
	}
	return m.DbInstance.Exec(s)
}

func (m *MysqlOps) Ping() error {
	if m.DbInstance == nil {
		return errors.New("mysqlops.DbInstance is nil")
	}
	return m.DbInstance.Ping()
}

func (m *MysqlOps) Query(s string) (*sql.Rows, error) {
	if m.DbInstance == nil {
		return nil, errors.New("mysqlops.DbInstance is nil")
	}
	return m.DbInstance.Query(s)
}

func (m *MysqlOps) QueryRow(s string) *sql.Row {
	if m.DbInstance == nil {
		return nil
	}
	return m.DbInstance.QueryRow(s)
}

func (m *MysqlOps) Prepare(s string) (*sql.Stmt, error) {
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
func (m *MysqlOps) BatchInserts(table string, fields []string, values [][]interface{}) (sql.Result, error) {
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
				tmp = append(tmp, fmt.Sprintf("%v", e))
			}
		}
		tmpString := fmt.Sprintf("(%s)", strings.Join(tmp, ","))
		keyValueFormatStrings = append(keyValueFormatStrings, tmpString)
	}

	valueSql := fmt.Sprintf("insert into %s (%s) values %s", table, strings.Join(fields, ","), strings.Join(keyValueFormatStrings, ","))
	return m.Exec(valueSql)
}

func RunMysqlSample() {
	db := MysqlOps{UserName: "root", PassWord: "123456", IP: "127.0.0.1", Port: "3306", DataBase: "provincecity"}
	err := db.Open()
	if err != nil {
		fmt.Println("database instance create err!")
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("database connect err! ", err)
		return
	}

	fieldsList := []string{"city_id", "city_name"}
	fieldValues := [][]interface{}{{0, "a"}, {0, "b"}}
	db.BatchInserts("city", fieldsList, fieldValues)
	//output:insert into city (city_id,city_name) values (0,'a'),(0,'b')
}
