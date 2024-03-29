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
	"strings"
)

var (
	VulhublogDir         = filepath.Join(utils.GetCurrentPath(), "data", "vulhub")
	VulhubOldInfoPath    = filepath.Join(VulhublogDir, "vulhub_old_info.log")
	VulhubUpdateInfoPath = filepath.Join(VulhublogDir, "vulhub_update_info.log")
)

func InitVulhubExp() {
	fileList, _ := GetTargetFiles("vulhub", "vulhub", "", "docker-compose.yml")
	contents := strings.Join(fileList, "\n")
	utils.WriteToLog(contents, VulhubOldInfoPath)
}

func CheckVulhubUpdate() string {
	file, _ := os.Open(VulhubOldInfoPath)
	defer file.Close()
	files := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		files[line] = true
	}
	var result strings.Builder
	update := false
	fileList, _ := GetTargetFiles("vulhub", "vulhub", "", "docker-compose.yml")
	contents := strings.Join(fileList, "\n")
	err := ioutil.WriteFile(VulhubOldInfoPath, []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}
	var vulnerabilities []Vulnerability
	for _, data := range fileList {
		if _, ok := files[data]; !ok {
			update = true
			dataParts := strings.Split(data, "/")
			data = dataParts[1] + "/" + dataParts[2]
			var cve string
			if strings.Contains(dataParts[2], "CVE") {
				cve = dataParts[2]
			}
			result.WriteString(fmt.Sprintf("%s\n", data))
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   data,
				CVE:    cve,
				URL:    "https://github.com/vulhub/vulhub",
				Source: "vulhub",
			})

		}
	}

	if update {
		utils.PrintLog("success", "vulhub update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	} else {
		utils.PrintLog("info", "vulhub is up to date")
		result.WriteString("Already up to date.")
	}

	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), VulhubUpdateInfoPath)
	return result.String()
}
