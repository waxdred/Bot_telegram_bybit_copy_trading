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
	CreateTable("admin", "admin", "admin1", order.Db, 100)
	CreateTable("api", "api", "api_secret", order.Db, 100)
	SelectApi("db.api", order.Db, api)
	SelectAdmin("db.admin", order.Db, api)
	if len(api.Api) > 0 {
		if CheckApi("db.api", order.Db, api.Api[0].Api) == true {
			InsertApi(api.Api[0].Api, api.Api[0].Api_secret, "api", order.Db)
		}
	}
	if len(api.Admin) > 0 {
		if CheckAdmin("db.admin", order.Db, api.Admin[0]) == true {
			InsertAdmin(api.Admin[0], "admin", order.Db)
		}
	}
	return nil
}

func InsertApi(api string, api_secret string, dbname string, db *sql.DB) error {
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

func InsertAdmin(adm string, dbname string, db *sql.DB) error {
	dbInsert := fmt.Sprint("Insert INTO `db`.`admin` (admin, admin1) VALUES (?, ?)")
	insert, err := db.Prepare(dbInsert)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = insert.Exec(adm, adm)
	if err != nil {
		log.Println(err)
		return err
	}
	defer insert.Close()
	return nil
}

func SelectApi(dbname string, db *sql.DB, apiEnv *data.Env) (*sql.Rows, error) {
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

func SelectAdmin(dbname string, db *sql.DB, apiEnv *data.Env) (*sql.Rows, error) {
	dbSelect := fmt.Sprint("SELECT * FROM ", dbname)
	lst, err := db.Query(dbSelect)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for lst.Next() {
		var id int64
		var admin, admin1 string
		lst.Scan(&id, &admin, &admin1)
		log.Printf("%s %s", admin, admin1)
		apiEnv.AddAdmin(admin)
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
			log.Print("api exist")
			return false
		}
	}
	return true
}

func CheckAdmin(dbname string, db *sql.DB, adm string) bool {
	dbSelect := fmt.Sprint("SELECT * FROM ", dbname)
	lst, err := db.Query(dbSelect)
	if err != nil {
		return false
	}
	defer lst.Close()
	for lst.Next() {
		var id int64
		var admin, admin1 string
		lst.Scan(&id, &admin, &admin1)
		if adm == admin {
			log.Print("admin exist")
			return false
		}
	}
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
	dbCreate := fmt.Sprint(
		"CREATE TABLE `db`.`",
		tablename,
		"` (`id` INT NOT NULL AUTO_INCREMENT,`",
		db1,
		"` VARCHAR(100) NOT NULL,`",
		db2,
		"` VARCHAR(100) NOT NULL, PRIMARY KEY (`id`));",
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
