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
	logpath := filepath.Join(updatelogDir, yesterday)
	utils.PrintColor("info", "Monitor update")
	utils.PrintColor("info", "Check Exploitdb")
	edbinfo := crawler.CheckEdbUpdate() + "\n"
	utils.PrintColor("info", "Check Github")
	githubinfo := crawler.CheckGithubUpdate() + "\n"
	utils.PrintColor("info", "Check Metasploit")
	crawler.UpdateMSF()
	msfinfo := crawler.CheckMSFUpdate() + "\n\n"
	utils.PrintColor("info", "Check Seebug")
	seebuginfo := crawler.CheckSeebugUpdate() + "\n\n"
	utils.PrintColor("info", "Check Vulhub")
	crawler.UpdateVulhub()
	vulhubinfo := crawler.CheckVulhubUpdate() + "\n\n"
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
    <h2>Vulnerability Update Monitor</h2>
	<h3>%s</h3>
    <div class="section">
        <h3>Exploit-db</h3>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h3>Github</h3>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h3>Metasploit</h3>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h3>Seebug</h3>
        <p>
            %s
        </p>
    </div>
    <div class="section">
        <h3>Vulhub</h3>
        <p>
            %s
        </p>
    </div>
</body>
</html>
`, yesterday, txt2html(edbinfo), txt2html(githubinfo), txt2html(msfinfo), txt2html(seebuginfo), txt2html(vulhubinfo))

	utils.WriteToLog(htmlTemplate, logpath)
	mail.Sendmail()
}
