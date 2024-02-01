package database

import (
	"crush/config"
	"crush/utils"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var githubdb *sql.DB

func init() {
	conf := config.LoadConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", conf.Database.DBUsername, conf.Database.DBPassword, conf.Database.Name)
	var err error
	githubdb, err = sql.Open("mysql", connectionString)
	if err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
	githubdb.SetMaxIdleConns(20)
	githubdb.SetMaxOpenConns(50)

	if err := githubdb.Ping(); err != nil {
		utils.PrintColor("error", "Error connecting mysql database:", err, "\n")
		return
	}
}

func CreateGithubDB() {
	// 创建 exploit_db 表
	_, err := githubdb.Exec(`
		CREATE TABLE IF NOT EXISTS github_db (
			CVE VARCHAR(50),
		    description TEXT,
			date_published DATE,
			CVSS2 float,
			CVSS3 float,
			CNA3 float,
			poc_url VARCHAR(512)
		) DEFAULT CHARSET=utf8mb4;
	`)
	if err != nil {
		utils.PrintColor("error", "Error creating table github_db:", err, "\n")
		return
	}

}

func InsertGithubDB(cve string, description string, date time.Time, cvss2, cvss3, cna3 float64, poc_url string) error {
	_, err := githubdb.Exec(`
		INSERT INTO github_db (CVE, description, date_published, CVSS2, CVSS3, CNA3, poc_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, cve, description, date, cvss2, cvss3, cna3, poc_url)

	if err != nil {
		utils.PrintColor("error", "Error inserting vul %v into table github_db: %v\n", poc_url, err)
		return err
	}

	return nil
}

// select count(*) from github_db where CVE=%s
func CountData(CVE string) (int, error) {
	var count int
	err := githubdb.QueryRow(`SELECT count(*) from github_db where CVE= ?`, CVE).Scan(&count)
	if err != nil {
		utils.PrintColor("error", "Error query %s from github_db: %v", CVE, err)
		return 0, err
	}
	return count, nil
}
