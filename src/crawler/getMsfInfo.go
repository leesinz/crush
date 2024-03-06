package crawler

import (
	"bufio"
	"crush/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	newExploitInfo       = regexp.MustCompile(`(?: create mode.{1,12})(modules/[a-zA-z0-9/_]+\.rb)`)
	cvePattern           = regexp.MustCompile(`(?:\')(CVE)(?:\', \')([0-9]{4}-[0-9]{4,10})(?:\')`)
	newCVEInfos          = make(map[string]string)
	msflogDir            = filepath.Join(utils.GetParentPath(), "data", "metasploit", "log")
	msfupdateInfoPath    = filepath.Join(msflogDir, "msf_update_info.log")
	msfupdateHistoryPath = filepath.Join(msflogDir, "msf_update_history.log ")
	msfDir               = cfg.MSF.MsfDir
)

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitMSF() {
	err := utils.GitClone("https://github.com/rapid7/metasploit-framework.git", msfDir)
	if err != nil {
		fmt.Println("git clone msf failed:", err)
		return
	}
}

func UpdateMSF() {
	os.MkdirAll(filepath.Dir(msfupdateInfoPath), 0755)
	err := runCommand("bash", "-c", "date +%Y-%m-%d\\ %H:%M:%S > "+msfupdateInfoPath)
	if err != nil {
		log.Fatal(err)
	}
	err = runCommand("bash", "-c", "cd "+cfg.MSF.MsfDir+"; git pull >> "+msfupdateInfoPath)
	if err != nil {
		log.Fatal(err)
	}
	CheckMSFUpdate()
	err = runCommand("bash", "-c", "cat "+msfupdateInfoPath+" >> "+msfupdateHistoryPath)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckMSFUpdate() string {
	var result strings.Builder
	updated := false
	file, err := os.Open(msfupdateInfoPath)
	if err != nil {
		fmt.Println("open file err:", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if newExploitInfo.MatchString(line) {
			updated = true
			cveFlag := ""
			newFilePath := newExploitInfo.FindStringSubmatch(line)[1]
			file, err := os.Open(msfDir + newFilePath)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return ""
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				rbLine := scanner.Text()
				if cvePattern.MatchString(rbLine) {
					sub := cvePattern.FindStringSubmatch(rbLine)
					cve := sub[1] + "-" + sub[2]
					newCVEInfos[cve] = newFilePath
					cveFlag = cve + ":"
				}
			}
			result.WriteString(cveFlag + newFilePath + "\n")
		}
	}
	if !updated {
		result.WriteString("Already up to date.")
	} else {
		utils.PrintColor("success", "Metasploit Updated")
	}
	return result.String()
}
