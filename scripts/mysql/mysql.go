/*
Package mysql implements a simple way of obtaining a SQL Connection
*/
package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)
import "fmt"

//MySQL ...
type MySQL struct {
	Conn   *sql.DB
	schema string
}

//CreateConnection ...
func CreateConnection() *MySQL {
	var mySQL MySQL
	fmt.Println("Creating new sql connection")
	mySQL.schema = "db10"
	connGlobal, _ := sql.Open("mysql", "root:root@/")
	connGlobal.Exec("CREATE DATABASE IF NOT EXISTS " + mySQL.schema + ";")
	connGlobal.Close()
	conn, _ := sql.Open("mysql", "root:root@tcp(172.17.0.2:3306)/"+mySQL.schema+"?parseTime=true")
	mySQL.Conn = conn

	return &mySQL
}

//Close ...
func (f MySQL) Close() {
	f.Conn.Close()
	fmt.Println("Closing the sql connection")
}

//DoesTableExist ...
func (f MySQL) DoesTableExist(name string) bool {
	var count int
	row := f.Conn.QueryRow(
		`SELECT count(*) 
        FROM information_schema.tables
        WHERE table_schema = '` + f.schema + `'
        AND table_name = '` + name + `'
        LIMIT 1;`)
	row.Scan(&count)
	return count == 1
}
