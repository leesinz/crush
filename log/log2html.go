package log

import (
	"crush/utils"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func str2html(content string) string {
	return strings.Replace(content, "\n", "<br>", -1)
}
func Log2Html(edb, github, msf, seebug, vulhub, zeroday, packetstorm, afrog, nuclei, poc string) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	updatelogDir := filepath.Join(utils.GetCurrentPath(), "data", "updateinfo")
	logpath := filepath.Join(updatelogDir, yesterday)

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

			<div class="section">
		        <h3>0day.today</h3>
		        <p>
		            %s
		        </p>
		    </div>

			<div class="section">
		        <h3>PacketStorm</h3>
		        <p>
		            %s
		        </p>
		    </div>

		<div class="section">
		        <h3>Afrog</h3>
		        <p>
		            %s
		        </p>

		<div class="section">
		        <h3>Nuclei</h3>
		        <p>
		            %s
		        </p>

		<div class="section">
		        <h3>POC</h3>
		        <p>
		            %s
		        </p>

		</body>
		</html>
		`, yesterday, str2html(edb), str2html(github), str2html(msf), str2html(seebug), str2html(vulhub), str2html(zeroday), str2html(packetstorm), str2html(afrog), str2html(nuclei), str2html(poc))
	utils.WriteToLog(htmlTemplate, logpath)
}
