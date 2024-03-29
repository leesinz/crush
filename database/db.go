package database

import (
	"crush/config"
	"crush/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	zerodaydb, exploitdb, githubdb, packetstormdb, seebugdb *sql.DB
)

func initDB() (*sql.DB, error) {

	conf := config.LoadConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s", conf.Database.DBUsername, conf.Database.DBPassword, conf.Database.DBPort, conf.Database.Name)
	var err error
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		utils.PrintLog("error", "Error connecting mysql database:", err, "\n")
		return nil, err
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(50)

	if err := db.Ping(); err != nil {
		utils.PrintLog("error", "Error connecting mysql database:", err, "\n")
		return nil, err
	}

	return db, nil
}

func init() {
	zerodaydb, _ = initDB()
	exploitdb, _ = initDB()
	githubdb, _ = initDB()
	packetstormdb, _ = initDB()
	seebugdb, _ = initDB()
}

func createTable(db *sql.DB, tableName, createQuery string) {
	_, err := db.Exec(createQuery)
	if err != nil {
		utils.PrintLog("error", fmt.Sprintf("Error creating table %s: %v", tableName, err))
		return
	}
}

func insertData(db *sql.DB, insertQuery string, args ...interface{}) error {
	_, err := db.Exec(insertQuery, args...)
	if err != nil {
		utils.PrintLog("error", fmt.Sprintf("Error inserting data into table: %v\n", err))
		return err
	}
	return nil
}

func CreateZerodayDB() {
	createQuery := `
		CREATE TABLE IF NOT EXISTS zeroday (
			ID VARCHAR(50),
		    Name TEXT,
		    Date DATE,
			Category VARCHAR(50),
			CVE TEXT,
		    Risk VARCHAR(50),
			POC TEXT	    
		) DEFAULT CHARSET=utf8mb4;
	`
	createTable(zerodaydb, "zeroday", createQuery)
}

func CreateEDB() {
	createQuery := `
		CREATE TABLE IF NOT EXISTS exploit_db (
		    id INT,
			description VARCHAR(255),
			type VARCHAR(50),
			platform VARCHAR(50),
			date_published DATE,
			verified INT,
			cve VARCHAR(2048),
			osvdb VARCHAR(2048),
			otherNum VARCHAR(2048)
		) DEFAULT CHARSET=utf8mb4;
	`
	createTable(exploitdb, "exploit_db", createQuery)
}

func CreateGithubDB() {
	createQuery := `
		CREATE TABLE IF NOT EXISTS github (
			CVE VARCHAR(50),
		    description TEXT,
			date_published DATE,
			CVSS2 float,
			CVSS3 float,
			CNA3 float,
			poc_url VARCHAR(512)
		) DEFAULT CHARSET=utf8mb4;
	`
	createTable(githubdb, "github", createQuery)

}

func CreatePacketStormDB() {
	createQuery := `
		CREATE TABLE IF NOT EXISTS packetstorm (
			ID VARCHAR(50),
			Date DATE,
			Name TEXT,
			CVE TEXT,
			POC TEXT,
		    Description TEXT
			) DEFAULT CHARSET=utf8mb4;
`
	createTable(packetstormdb, "packetstorm", createQuery)
}

func CreateSeebugDB() {
	createQuery := `
		CREATE TABLE IF NOT EXISTS seebug (
			SSVID VARCHAR(50),
			SubmitTime DATE,
		    Severity VARCHAR(50),
			Name TEXT,
			CVE TEXT,
			HasPOC BOOLEAN
		) DEFAULT CHARSET=utf8mb4;
	`
	createTable(seebugdb, "seebug", createQuery)
}

func InsertZerodayDB(id string, name string, date time.Time, category string, cve string, risk string, poc string) error {
	insertQuery := `
		INSERT INTO zeroday (ID, Name, Date, Category, CVE, Risk, Poc)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	return insertData(zerodaydb, insertQuery, id, name, date, category, cve, risk, poc)
}

func InsertEDB(id int, description, exploitType, platform string, datePublished time.Time, verified int, cve, osvdb, otherNum string) error {
	insertQuery := `
		INSERT INTO exploit_db (id, description, type, platform, date_published, verified, cve, osvdb, otherNum)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	return insertData(exploitdb, insertQuery, id, description, exploitType, platform, datePublished, verified, cve, osvdb, otherNum)
}

func InsertGithubDB(cve string, description string, date time.Time, cvss2, cvss3, cna3 float64, poc_url string) error {
	insertQuery := `
		INSERT INTO github (CVE, description, date_published, CVSS2, CVSS3, CNA3, poc_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`
	return insertData(githubdb, insertQuery, cve, description, date, cvss2, cvss3, cna3, poc_url)
}

func CountGithubInfo(CVE string) (int, error) {
	var count int
	err := githubdb.QueryRow(`SELECT count(*) from github where CVE= ?`, CVE).Scan(&count)
	if err != nil {
		utils.PrintLog("error", "Error query %s from github: %v", CVE, err)
		return 0, err
	}
	return count, nil
}

func InsertPacketstormDB(id string, date time.Time, name string, cve string, poc string, desc string) error {
	insertQuery := `
		INSERT INTO packetstorm (ID, Date, Name, CVE, POC, Description)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	return insertData(packetstormdb, insertQuery, id, date, name, cve, poc, desc)

}

func InsertSeebugDB(ssvid string, submittime time.Time, severity string, name string, cve string, haspoc bool) error {
	insertQuery := `
		INSERT INTO seebug (SSVID, SubmitTime, Severity, Name, CVE, HasPoc)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	return insertData(seebugdb, insertQuery, ssvid, submittime, severity, name, cve, haspoc)

}
