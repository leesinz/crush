package loginfo

import (
	"crush/crawler"
	"crush/mail"
	"crush/utils"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

var (
	updatelogDir = filepath.Join(utils.GetParentPath(), "data", "updateinfo", "log")
)

func txt2html(content string) string {
	return strings.Replace(content, "\n", "<br>", -1)
}

func LogScraper() {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	logPath := filepath.Join(updatelogDir, yesterday)
	utils.PrintColor("info", "Monitor Update")
	utils.PrintColor("info", "Check Exploit-db")
	edbInfo := crawler.CheckEdbUpdate()

	utils.PrintColor("info", "Check Github")
	githubInfo := crawler.CheckGithubUpdate()

	utils.PrintColor("info", "Check Metasploit")
	crawler.UpdateMSF()
	msfInfo := crawler.CheckMSFUpdate()

	utils.PrintColor("info", "Check Seebug")
	seebugInfo := crawler.CheckSeebugUpdate()

	utils.PrintColor("info", "Check Vulhub")
	crawler.UpdateVulhub()
	vulhubInfo := crawler.CheckVulhubUpdate()

	utils.PrintColor("info", "Check 0day.today")
	zerodayInfo := crawler.Check0dayUpdate()

	utils.PrintColor("info", "Check PacketStorm")
	packetstormInfo := crawler.CheckPacketstormUpdate()

	//TODO
	htmlTemplate := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vulnerability Update Monitor</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
        }
        h2 {
            color: #333;
        }
        .section {
            margin-bottom: 20px;
        }
        .section h3 {
            margin-bottom: 10px;
            color: #666;
        }
        .section p {
            margin: 0;
        }
    </style>
</head>
<body>
    <h1>Vulnerability Update Monitor</h1>
	<h2>%s</h2>
    <div class="section">
        <h2>Exploit-db</h2>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h2>Github</h2>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h2>Metasploit</h2>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h2>Seebug</h2>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h2>Vulhub</h2>
        <p>
            %s
        </p>
    </div>
	
	<div class="section">
        <h2>0day.today</h2>
        <p>
            %s
        </p>
    </div>

	<div class="section">
        <h2>PacketStorm</h2>
        <p>
            %s
        </p>
    </div>
</body>
</html>
`, yesterday, txt2html(edbInfo), txt2html(githubInfo), txt2html(msfInfo), txt2html(seebugInfo), txt2html(vulhubInfo), txt2html(zerodayInfo), txt2html(packetstormInfo))

	utils.WriteToLog(htmlTemplate, logPath)
	mail.Sendmail()
}
