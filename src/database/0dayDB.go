package database

import (
	"crush/config"
	"crush/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var zerodaydb *sql.DB

func init() {
	conf := config.LoadConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s", conf.Database.DBUsername, conf.Database.DBPassword, conf.Database.DBPort, conf.Database.Name)
	var err error
	zerodaydb, err = sql.Open("mysql", connectionString)
	if err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
	zerodaydb.SetMaxIdleConns(20)
	zerodaydb.SetMaxOpenConns(50)

	if err := zerodaydb.Ping(); err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
}

func CreateZerodayDB() {
	_, err := zerodaydb.Exec(`
		CREATE TABLE IF NOT EXISTS zeroday_db (
			ID VARCHAR(50),
		    Name TEXT,
		    Date DATE,
			Category VARCHAR(50),
			CVE TEXT,
		    Risk VARCHAR(50),
			POC TEXT	    
		) DEFAULT CHARSET=utf8mb4;
	`)
	if err != nil {
		utils.PrintColor("error", "Error creating table zeroday_db:", err, "\n")
		return
	}

}

func InsertZerodayDB(id string, name string, date time.Time, category string, cve string, risk string, poc string) error {
	_, err := zerodaydb.Exec(`
		INSERT INTO zeroday_db (ID, Name, Date, Category, CVE, Risk, Poc)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, name, date, category, cve, risk, poc)

	if err != nil {
		utils.PrintColor("error", "Error inserting vul %v into table zeroday_db: %v\n", id, err)
		return err
	}

	return nil
}
