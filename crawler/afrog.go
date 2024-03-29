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
	AfroglogDir         = filepath.Join(utils.GetCurrentPath(), "data", "afrog")
	AfrogOldInfoPath    = filepath.Join(AfroglogDir, "afrog_old_info.log")
	AfrogUpdateInfoPath = filepath.Join(AfroglogDir, "afrog_update_info.log")
)

func InitAfrogExp() {
	fileList, _ := GetTargetFiles("zan8in", "afrog", "pocs/afrog-pocs", ".yaml")
	contents := strings.Join(fileList, "\n")
	utils.WriteToLog(contents, AfrogOldInfoPath)
}

func CheckAfrogUpdate() string {
	file, _ := os.Open(AfrogOldInfoPath)
	defer file.Close()
	files := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		files[line] = true
	}
	var result strings.Builder
	update := false
	fileList, _ := GetTargetFiles("zan8in", "afrog", "pocs/afrog-pocs", ".yaml")
	contents := strings.Join(fileList, "\n")
	err := ioutil.WriteFile(AfrogOldInfoPath, []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}
	var vulnerabilities []Vulnerability
	for _, data := range fileList {
		if _, ok := files[data]; !ok {
			update = true
			url := "https://raw.githubusercontent.com/zan8in/afrog/main/" + data
			dataParts := strings.Split(strings.TrimSuffix(data, ".yaml"), "/")
			data = dataParts[len(dataParts)-1]
			result.WriteString(fmt.Sprintf("%s\n", data))
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   data,
				CVE:    "",
				URL:    url,
				Source: "afrog",
			})

		}
	}

	if update {
		utils.PrintLog("success", "afrog update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	} else {
		utils.PrintLog("info", "afrog is up to date")
		result.WriteString("Already up to date.")
	}

	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), AfrogUpdateInfoPath)
	return result.String()
}
