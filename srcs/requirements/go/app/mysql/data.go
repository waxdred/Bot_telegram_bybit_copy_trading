package mysql

import (
	"bot/data"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectionDb(order *data.Bot, api *data.Env) error {
	var DbErr error

	order.Db, DbErr = DbConnect("root:bot@tcp(mysql:3306)/db")
	if DbErr != nil {
		return errors.New("Coudn't connect to database")
	}
	// check if table exits if not create it
	CreateTable("api", "api", "api_secret", order.Db, 100)
	Select("db.api", order.Db, api)
	if len(api.Api) > 0 {
		if CheckApi("db.api", order.Db, api.Api[0].Api) == true {
			Insert(api.Api[0].Api, api.Api[0].Api_secret, "db.api", order.Db)
		}
	}
	return nil
}

func Insert(api string, api_secret string, dbname string, db *sql.DB) error {
	dbInsert := fmt.Sprint("Insert INTO `db`.`api` (api, api_secret) VALUES (?, ?)")
	insert, err := db.Prepare(dbInsert)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = insert.Exec(api, api_secret)
	if err != nil {
		log.Println(err)
		return err
	}
	defer insert.Close()
	return nil
}

func Select(dbname string, db *sql.DB, apiEnv *data.Env) (*sql.Rows, error) {
	dbSelect := fmt.Sprint("SELECT * FROM ", dbname)
	lst, err := db.Query(dbSelect)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for lst.Next() {
		var id int64
		var api, api_secret string
		lst.Scan(&id, &api, &api_secret)
		log.Printf("%s %s", api, api_secret)
		apiEnv.AddApi(api, api_secret)
	}
	defer lst.Close()
	return lst, nil
}

func CheckApi(dbname string, db *sql.DB, api string) bool {
	dbSelect := fmt.Sprint("SELECT * FROM ", dbname)
	lst, err := db.Query(dbSelect)
	if err != nil {
		return false
	}
	defer lst.Close()
	for lst.Next() {
		var id int64
		var db_api, db_api_secret string
		lst.Scan(&id, &db_api, &db_api_secret)
		if api == db_api {
			log.Print("api exit")
			return false
		}
	}
	log.Print("api not found")
	return true
}

func DbConnect(connection string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connection)
	log.Println(db)
	if err != nil {
		log.Panic(err)
	} else {
		log.Print("Connection to database Success\n")
	}
	return db, nil
}

func DbDelete(tablename string, api string, db *sql.DB) error {
	dbDelette := fmt.Sprintf("DELETE FROM `db`.`api` WHERE `api` = ?")
	drop, err := db.Prepare(dbDelette)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = drop.Exec(api)
	if err != nil {
		log.Println(err)
		return err
	}
	defer drop.Close()
	return nil
}

func CreateTable(tablename string, db1 string, db2 string, db *sql.DB, size int64) {
	dbCreate := fmt.Sprintf(
		"CREATE TABLE `db`.`api` (`id` INT NOT NULL AUTO_INCREMENT,`api` VARCHAR(100) NOT NULL,`api_secret` VARCHAR(100) NOT NULL, PRIMARY KEY (`id`));",
	)

	data, err := db.Query(dbCreate)
	if err != nil {
		log.Print("error: mysql")
		log.Println(err)
		return
	} else {
		log.Printf("Create table")
	}
	defer data.Close()
}
