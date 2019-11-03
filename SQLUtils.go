package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./torch_peak")
	checkErr(err)
	fmt.Println("查询数据")
	rows, err := db.Query("SELECT * FROM torch_encode")
	checkErr(err)
	for rows.Next() {
		var _id int
		var filename string
		var encode_0 float32
		var encode_1 float32
		err = rows.Scan(&_id, &filename, &encode_0, &encode_1)
		checkErr(err)
		fmt.Println(_id, filename, encode_0, encode_1)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}