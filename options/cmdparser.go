package options

import (
	"crush/crawler"
	"crush/database"
	"crush/log"
	"crush/mail"
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
		initDB()
	case "monitor":
		monitor()
	default:
		utils.Help()
	}
}

func initDB() {
	utils.PrintLog("info", "Start initializing tasks")
	database.CreateZerodayDB()
	database.CreateEDB()
	database.CreateGithubDB()
	database.CreatePacketStormDB()
	database.CreateSeebugDB()
	utils.PrintLog("success", "Create database successfully")
	utils.PrintLog("info", "Start crawling msf history info")
	crawler.InitMSFExp()
	utils.PrintLog("info", "Start crawling vulhub history info")
	crawler.InitVulhubExp()
	utils.PrintLog("info", "Start crawling nuclei history info")
	crawler.InitNucleiExp()
	utils.PrintLog("info", "Start crawling afrog history info")
	crawler.InitAfrogExp()
	utils.PrintLog("info", "Start crawling poc history info")
	crawler.InitPOCExp()
	utils.PrintLog("success", "Init history info successfully")
}

func monitor() {
	utils.PrintLog("info", "Start monitoring tasks")
	edbInfo := crawler.CheckEdbUpdate() + "\n"
	githubInfo := crawler.CheckGithubUpdate() + "\n"
	msfInfo := crawler.CheckMSFUpdate() + "\n"
	seebugInfo := crawler.CheckSeebugUpdate() + "\n"
	vulhubInfo := crawler.CheckVulhubUpdate() + "\n"
	zerodayInfo := crawler.Check0dayUpdate() + "\n"
	packetstormInfo := crawler.CheckPacketstormUpdate() + "\n"
	aforgInfo := crawler.CheckAfrogUpdate() + "\n"
	nucleiInfo := crawler.CheckNucleiUpdate() + "\n"
	pocInfo := crawler.CheckPOCUpdate() + "\n"
	log.Log2Html(edbInfo, githubInfo, msfInfo, seebugInfo, vulhubInfo, zerodayInfo, packetstormInfo, aforgInfo, nucleiInfo, pocInfo)
	mail.Sendmail()
}
