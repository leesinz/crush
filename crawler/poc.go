package crawler

import (
	"bufio"
	"crush/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	POClogDir         = filepath.Join(utils.GetCurrentPath(), "data", "poc")
	POCOldInfoPath    = filepath.Join(POClogDir, "poc_old_info.log")
	POCUpdateInfoPath = filepath.Join(POClogDir, "poc_update_info.log")
)

func InitPOCExp() {
	fileList, _ := GetTargetFiles("wy876", "POC", "", ".md")
	contents := strings.Join(fileList, "\n")
	utils.WriteToLog(contents, POCOldInfoPath)
}

func CheckPOCUpdate() string {
	file, _ := os.Open(POCOldInfoPath)
	defer file.Close()
	files := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		files[line] = true
	}
	var result strings.Builder
	update := false
	fileList, _ := GetTargetFiles("wy876", "POC", "", ".md")
	contents := strings.Join(fileList, "\n")
	err := ioutil.WriteFile(POCOldInfoPath, []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}

	var vulnerabilities []Vulnerability
	for _, data := range fileList {
		if _, ok := files[data]; !ok {
			update = true
			url := "https://raw.githubusercontent.com/wy876/POC/main" + data
			dataParts := strings.Split(strings.TrimSuffix(data, ".md"), "/")
			data = dataParts[1]

			var cve string
			re := regexp.MustCompile(`CVE-\d{4}-\d{4}`)
			cveMatches := re.FindStringSubmatch(data)
			if len(cveMatches) > 0 {
				cve = cveMatches[0]
			}
			result.WriteString(fmt.Sprintf("%s\n", data))
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   data,
				CVE:    cve,
				URL:    url,
				Source: "poc",
			})

		}
	}

	if update {
		utils.PrintLog("success", "poc update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	} else {
		utils.PrintLog("info", "poc is up to date")
		result.WriteString("Already up to date.")
	}

	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), POCUpdateInfoPath)
	return result.String()
}
