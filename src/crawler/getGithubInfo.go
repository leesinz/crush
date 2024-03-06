package crawler

import (
	"crush/config"
	"crush/database"
	"crush/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Repo struct {
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
	Fork        bool      `json:"fork"`
	ContentsURL string    `json:"contents_url"`
}

type GithubRes struct {
	Items []Repo `json:"items"`
}

var (
	cfg                  = config.LoadConfig()
	githublogDir         = filepath.Join(utils.GetParentPath(), "data", "github", "log")
	githubupdateInfoPath = filepath.Join(githublogDir, "github_update_info.log")
)

func getGithubHeader() map[string]string {
	token := cfg.Github.GithubToken
	return map[string]string{
		"Authorization": "token " + token,
	}
}

func CheckGithubUpdate() string {
	var result strings.Builder
	updated := false
	year := time.Now().Year()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=CVE-%d&sort=updated&per_page=100&page=1", year)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	headers := getGithubHeader()
	utils.SetHeaders(req, headers)

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)

	var response GithubRes
	err = json.Unmarshal(content, &response)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return ""
	}
	for _, repo := range response.Items {
		if repo.Fork {
			continue
		}
		cve := utils.ExtractCVEFromName(repo.Name)
		if cve == "" {
			continue
		}
		contentURL := strings.Split(repo.ContentsURL, "{")[0]
		res, err := http.Get(contentURL)
		if err != nil {
			utils.PrintColor("error", "Error request", contentURL, "\n")
			continue
		}
		body, _ := ioutil.ReadAll(res.Body)
		if strings.Contains(string(body), "This repository is empty") {
			continue
		}
		date_published := repo.CreatedAt
		exists, desc, cvss2_tmp, cvss3_tmp, cna_tmp, err := utils.GetCVEInfo(cve)
		if err != nil || !exists {
			continue
		}
		cvss2, _ := strconv.ParseFloat(cvss2_tmp, 64)
		cvss3, _ := strconv.ParseFloat(cvss3_tmp, 64)
		cna, _ := strconv.ParseFloat(cna_tmp, 64)
		poc_url := repo.HTMLURL
		if date_published.Format("2006-01-02") == yesterday && !database.CheckGithubDuplicate(poc_url) {
			updated = true
			tmp := fmt.Sprintf("%v\npoc_url:%v\ndesc:%v\ncvss2:%v\ncvss3:%v\ncna:%v\n\n", cve, poc_url, desc, cvss2, cvss3, cna)
			result.WriteString(tmp)
			err = database.InsertGithubDB(cve, desc, date_published, cvss2, cvss3, cna, poc_url)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if updated {
		utils.PrintColor("success", "Github Updated")
		utils.WriteToLog(yesterday+"\n"+result.String(), githubupdateInfoPath)
		return result.String()
	} else {
		utils.WriteToLog(yesterday+"\nAlready up to date.", githubupdateInfoPath)
		return "Already up to date."
	}
}
