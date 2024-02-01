package options

import (
	"crush/crawler"
	"crush/database"
	"crush/loginfo"
	"crush/utils"
	"os"
)

func ArgParser() {
	args := os.Args
	if len(args) < 2 {
		utils.Help()
		return
	}
	command := args[1]
	switch command {
	case "init":
		initTasks()
	case "monitor":
		monitor()
	}
}

func initTasks() {

	//exploitdb
	utils.PrintColor("info", "Commencing initialization of the task")
	database.CreateEDB()
	utils.PrintColor("info", "Create exploit-db database successfully")
	utils.PrintColor("info", "Start crawling exploit-db history info")
	crawler.ConcurrentEDBCrawler()

	database.CreateGithubDB()
	/*
		//github
		utils.PrintColor("info", "Create github database successfully")
		utils.PrintColor("info", "Start crawling github history info")
		crawler.GithubCrawler()

	*/

	database.CreateSeebugDB()
	/*
		//seebug
		utils.PrintColor("info", "Create seebug database successfully")
		utils.PrintColor("info", "Start crawling seebug history info")
		crawler.FetchSeebug()
	*/

	//msf
	utils.PrintColor("info", "Start git clone metasploit")
	crawler.InitMSF()

	//vulhub
	utils.PrintColor("info", "Start git clone vulhub")
	crawler.InitVulhub()

}

func monitor() {
	loginfo.LogScraper()
}
