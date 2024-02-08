package crawler

import (
	"bufio"
	"crush/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	newVulhubInfo           = regexp.MustCompile(`(?: create mode.+docker-compose.yml)`)
	vulhublogDir            = filepath.Join(utils.GetParentPath(), "data", "vulhub", "log")
	vulhubupdateInfoPath    = filepath.Join(vulhublogDir, "vulhub_update_info.log")
	vulhubupdateHistoryPath = filepath.Join(vulhublogDir, "vulhub_update_history.log ")
	vulhubDir               = cfg.Vulhub.VulhubDir
)

func Test() {
	fmt.Println(msfDir)
	fmt.Println(vulhubDir)

}

func InitVulhub() {
	fmt.Println(vulhubDir)
	err := utils.GitClone("https://github.com/vulhub/vulhub.git", vulhubDir)
	if err != nil {
		fmt.Println("git clone vulhub failed:", err)
		return
	}
}

func UpdateVulhub() error {
	os.MkdirAll(filepath.Dir(vulhubupdateInfoPath), 0755)
	err := runCommand("bash", "-c", "date +%Y-%m-%d\\ %H:%M:%S > "+vulhubupdateInfoPath)
	if err != nil {
		return err
	}
	err = runCommand("bash", "-c", "cd "+cfg.Vulhub.VulhubDir+"; git pull >> "+vulhubupdateInfoPath)
	if err != nil {
		return err
	}
	//check new exploits
	CheckVulhubUpdate()
	err = runCommand("bash", "-c", "cat "+vulhubupdateInfoPath+" >> "+vulhubupdateHistoryPath)
	if err != nil {
		return err
	}
	return nil
}

func CheckVulhubUpdate() string {
	var result strings.Builder
	updated := false
	file, err := os.Open(vulhubupdateInfoPath)
	if err != nil {
		fmt.Println("open file err:", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if newVulhubInfo.MatchString(line) {
			updated = true
			patterns := strings.Split(line, " ")
			last := patterns[len(patterns)-1]
			newVul := strings.TrimSuffix(last, "/docker-compose.yml")
			result.WriteString(newVul + "\n")
		}
		if !updated {
			result.WriteString("Already up to date.")
		}
	}
	return result.String()
}
