package util

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//TODO:
type database interface {
	Init() error
	Open() error
	/*.....*/
}

func RunMysqlSample() {
	db, _ := sql.Open("mysql", "root:123456@(127.0.0.1:3306)/provincecity")
	defer db.Close()
	err := db.Ping()
	if err != nil {
		fmt.Println("database connect error! ", err)
		return
	}

	//一：execute insert sql
	/*
	   sql:="insert into province values (2,'江苏')"
	   result,_:=db.Exec(sql)
	   n,_:=result.RowsAffected();
	   fmt.Println("affected row number is:",n)
	*/
	//二：execute prepare
	/*
	   province:=[2][2] string{{"3","上海"},{"4","北京"}}
	   stmt,_:=db.Prepare("insert into province values (?,?)")
	   for _,s:=range province{
	       stmt.Exec(s[0],s[1])
	   }
	*/

	//三：execute single row select sql
	/*
	   var id,name string
	   rows:=db.QueryRow("select * from province where id=4")
	   rows.Scan(&id,&name)
	   fmt.Println(id,"--",name)
	*/

	//四：execute rows select sql

	rows, _ := db.Query("select * from province")
	var id, name string
	for rows.Next() {
		rows.Scan(&id, &name)
		fmt.Println(id, "--", name)
	}
}
