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

	utils.PrintColor("info", "Commencing initialization of the task")
	database.CreateEDB()
	utils.PrintColor("info", "Create exploit_db database")

	database.CreateGithubDB()
	utils.PrintColor("info", "Create github_db database")

	database.CreateSeebugDB()
	utils.PrintColor("info", "Create seebug_db database")

	utils.PrintColor("info", "Git clone metasploit")
	crawler.InitMSF()

	utils.PrintColor("info", "Git clone vulhub")
	crawler.InitVulhub()

	database.CreateZerodayDB()
	utils.PrintColor("info", "Create zeroday_db database")

	database.CreatePacketstormDB()
	utils.PrintColor("info", "Create packetstorm_db database")

	utils.PrintColor("success", "Tasks init successfully")
}

func monitor() {
	loginfo.LogScraper()
}
