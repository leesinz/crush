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
	NucleilogDir         = filepath.Join(utils.GetCurrentPath(), "data", "nuclei")
	NucleiOldInfoPath    = filepath.Join(NucleilogDir, "nuclei_old_info.log")
	NucleiUpdateInfoPath = filepath.Join(NucleilogDir, "nuclei_update_info.log")
)

func InitNucleiExp() {
	fileList, _ := GetTargetFiles("projectdiscovery", "nuclei-templates", "", ".yaml")
	contents := strings.Join(fileList, "\n")
	utils.WriteToLog(contents, NucleiOldInfoPath)
}

func CheckNucleiUpdate() string {
	file, _ := os.Open(NucleiOldInfoPath)
	defer file.Close()
	files := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		files[line] = true
	}
	var result strings.Builder
	update := false
	fileList, _ := GetTargetFiles("projectdiscovery", "nuclei-templates", "", ".yaml")
	contents := strings.Join(fileList, "\n")
	err := ioutil.WriteFile(NucleiOldInfoPath, []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}
	var vulnerabilities []Vulnerability
	for _, data := range fileList {
		if _, ok := files[data]; !ok {
			update = true
			dataParts := strings.Split(strings.TrimSuffix(data, ".yaml"), "/")
			data = dataParts[len(dataParts)-2] + "/" + dataParts[len(dataParts)-1]
			var cve string
			if strings.Contains(dataParts[len(dataParts)-1], "CVE") {
				cve = dataParts[len(dataParts)-1]
			}
			result.WriteString(fmt.Sprintf("%s\n", data))
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   data,
				CVE:    cve,
				URL:    "https://github.com/projectdiscovery/nuclei-templates",
				Source: "nuclei",
			})

		}
	}

	if update {
		utils.PrintLog("success", "nuclei update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	} else {
		utils.PrintLog("info", "nuclei is up to date")
		result.WriteString("Already up to date.")
	}

	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), NucleiUpdateInfoPath)
	return result.String()
}
