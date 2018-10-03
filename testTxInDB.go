package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // driver for mysql
	"github.com/jmoiron/sqlx"
)

var createTable = `
CREATE TABLE account (
	id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	balance int DEFAULT '0'
)
`

func create() {
	db, err := sqlx.Open("mysql", "root:@tcp(127.0.0.1)/")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE testdb")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec("USE testdb")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalln(err)
	}

	sqlStr := `INSERT INTO account (balance) VALUES ('0')`
	_, err = db.Exec(sqlStr)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	//create()
	done := make(chan bool)
	db, err := sqlx.Connect("mysql", "root:@tcp(127.0.0.1)/testdb")
	if err != nil {
		log.Fatalln(err)
	}

	num := 2
	for i := 0; i < num; i++ {
		go func() {
			var balance int
			tx, err := db.Begin()
			if err != nil {
				log.Fatalln(err)
			}

			defer tx.Rollback()

			selectQuery := `SELECT balance FROM account WHERE id = 1`
			err = tx.QueryRow(selectQuery).Scan(&balance)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("balance:", balance)
			// each go-routine add 100 to account
			balance += 100

			updateQuery := `UPDATE account SET balance = ? WHERE id = 1`
			_, err = tx.Exec(updateQuery, balance)
			if err != nil {
				log.Fatalln(err)
			}
			// for waiting another goroutine
			time.Sleep(time.Second)
			tx.Commit()
			done <- true
		}()
	}

	//wait until database finish
	for i := 0; i < num; i++ {
		<-done
	}
}
