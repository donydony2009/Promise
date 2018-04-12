package main
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)
import "fmt"

func main() {
    var count int
	conn, _ := sql.Open("mysql", "root:root@/world")
    defer conn.Close()

    row := conn.QueryRow(
        `SELECT count(*) 
        FROM information_schema.tables
        WHERE table_schema = 'world' 
        AND table_name = 'city'
        LIMIT 1;`)
    row.Scan(&count)
    fmt.Printf("%d", count)
}