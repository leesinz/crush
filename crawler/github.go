package crawler

import (
	"crush/database"
	"crush/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	GithublogDir         = filepath.Join(utils.GetCurrentPath(), "data", "github")
	GithubupdateInfoPath = filepath.Join(GithublogDir, "github_update_info.log")
	GithubPocDir         = cfg.Github.PocDir
)

type GithubRes struct {
	Items []Repo `json:"items"`
}

type Repo struct {
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
	Fork        bool      `json:"fork"`
	ContentsURL string    `json:"contents_url"`
}

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
	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=CVE-%d&sort=updated&per_page=100&page=1", year)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		utils.PrintLog("error", "Error creating request: %v\n", err)
		return ""
	}
	headers := getGithubHeader()
	utils.SetHeaders(req, headers)

	res, err := client.Do(req)

	if err != nil {
		utils.PrintLog("error", "Error establishing connection: %v\n", err)
		return ""
	}
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)

	var response GithubRes
	err = json.Unmarshal(content, &response)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return ""
	}

	vulnerabilities := []Vulnerability{}
	for _, repo := range response.Items {
		if repo.Fork {
			continue
		}
		cve := utils.ExtractCVE(repo.Name)
		if cve == "" {
			continue
		}
		contentURL := strings.Split(repo.ContentsURL, "{")[0]
		res, err := http.Get(contentURL)
		if err != nil {
			utils.PrintLog("error", "Error request", contentURL, "\n")
			continue
		}
		body, _ := ioutil.ReadAll(res.Body)
		if strings.Contains(string(body), "This repository is empty") {
			continue
		}
		datePublished := repo.CreatedAt
		exists, desc, cvss2Tmp, cvss3Tmp, cnaTmp, err := utils.GetCVEInfo(cve)
		if err != nil || !exists {
			continue
		}
		cvss2, _ := strconv.ParseFloat(cvss2Tmp, 64)
		cvss3, _ := strconv.ParseFloat(cvss3Tmp, 64)
		cna, _ := strconv.ParseFloat(cnaTmp, 64)
		pocUrl := repo.HTMLURL
		existNum, _ := database.CountGithubInfo(cve)
		if datePublished.Format("2006-01-02") == Yesterday.Format("2006-01-02") {
			if existNum >= 5 {
				continue
			}
			updated = true
			result.WriteString(fmt.Sprintf("%v\npoc_url:%v\ndesc:%v\ncvss2:%v\ncvss3:%v\ncna:%v\n\n", cve, pocUrl, desc, cvss2, cvss3, cna))
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   cve,
				CVE:    cve,
				URL:    pocUrl,
				Source: "github",
			})

			err = database.InsertGithubDB(cve, desc, datePublished, cvss2, cvss3, cna, pocUrl)
			if err != nil {
				utils.PrintLog("error", "Error insert vul:", cve, "\n")
				continue
			}

			if DownloadPOC == true {
				err = utils.GitClone(pocUrl+".git", GithubPocDir+cve+"/"+strings.Split(pocUrl, "/")[3])
				if err != nil {
					continue
				}
			}
		}

	}

	if updated {
		utils.PrintLog("success", "github update")
		utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), GithubupdateInfoPath)
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	} else {
		utils.WriteToLog(Yesterday.Format("2006-01-02")+"\nAlready up to date.", GithubupdateInfoPath)
		utils.PrintLog("info", "github is up to date")
		result.WriteString("Already up to date.")
	}

	return result.String()
}
