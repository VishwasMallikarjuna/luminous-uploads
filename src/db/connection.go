package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/VishwasMallikarjuna/luminous-uploads/utils"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		utils.AppConfig.DBUser, utils.AppConfig.DBPassword, utils.AppConfig.DBHost, utils.AppConfig.DBPort, utils.AppConfig.DBName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Dberror:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB Ping", err)
	}

	fmt.Println("Successfully connected to the database")
}
