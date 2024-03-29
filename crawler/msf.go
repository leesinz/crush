package crawler

import (
	"bufio"
	"crush/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	MsflogDir         = filepath.Join(utils.GetCurrentPath(), "data", "msf")
	MsfOldInfoPath    = filepath.Join(MsflogDir, "msf_old_info.log")
	MsfUpdateInfoPath = filepath.Join(MsflogDir, "msf_update_info.log")
)

func InitMSFExp() {
	fileList, _ := GetTargetFiles("rapid7", "metasploit-framework", "modules/exploits", ".rb")
	contents := strings.Join(fileList, "\n")
	utils.WriteToLog(contents, MsfOldInfoPath)
}

func extractCVE(vul string) string {
	url := "https://raw.githubusercontent.com/rapid7/metasploit-framework/master/" + vul

	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error connecting msf:", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response", err)
	}

	cvePattern := regexp.MustCompile(`(?:\')(CVE)(?:\', \')([0-9]{4}-[0-9]{4,10})(?:\')`)

	cveMatches := cvePattern.FindAllStringSubmatch(string(body), -1)
	var cves []string
	for _, match := range cveMatches {
		cve := match[1] + "-" + match[2]
		cves = append(cves, cve)
	}
	cvesString := strings.Join(cves, ",")
	return cvesString
}

func CheckMSFUpdate() string {
	file, _ := os.Open(MsfOldInfoPath)
	defer file.Close()
	files := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		files[line] = true
	}
	var result strings.Builder
	update := false
	fileList, _ := GetTargetFiles("rapid7", "metasploit-framework", "modules/exploits", ".rb")
	contents := strings.Join(fileList, "\n")
	err := ioutil.WriteFile(MsfOldInfoPath, []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}
	var vulnerabilities []Vulnerability
	for _, data := range fileList {
		if _, ok := files[data]; !ok {
			update = true
			poc := "https://raw.githubusercontent.com/rapid7/metasploit-framework/master/" + data
			cve := extractCVE(data)
			dataParts := strings.Split(strings.TrimSuffix(data, ".rb"), "/")
			data = dataParts[len(dataParts)-1]
			result.WriteString(fmt.Sprintf("%s %s\n", data, cve))

			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   data,
				CVE:    cve,
				URL:    poc,
				Source: "msf",
			})
		}
	}

	if update {
		utils.PrintLog("success", "msf update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	} else {
		utils.PrintLog("info", "msf is up to date")
		result.WriteString("Already up to date.")
	}

	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), MsfUpdateInfoPath)
	return result.String()
}
