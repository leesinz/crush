package crawler

import (
	"crush/config"
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

// func fetchGithub(year, page int, wg *sync.WaitGroup, countchan chan struct{}) {
func fetchGithub(year, page int) (int, error) {
	//defer wg.Done()
	poc_dir := cfg.Github.PocDir
	count := 0
	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=CVE-%d&per_page=100&page=%d", year, page)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		utils.PrintColor("error", "Error creating request: %v\n", err)
		return 0, err
	}

	headers := getGithubHeader()
	utils.SetHeaders(req, headers)

	res, err := client.Do(req)

	if err != nil {
		utils.PrintColor("error", "Error establishing connection: %v\n", err)
		return 0, err
	}
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)

	var response GithubRes
	err = json.Unmarshal(content, &response)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return 0, err
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
		existNum, _ := database.CountData(cve)
		if existNum >= 5 {
			continue
		}

		err = database.InsertGithubDB(cve, desc, date_published, cvss2, cvss3, cna, poc_url)
		if err != nil {
			utils.PrintColor("error", "Error insert vul:", cve, "\n")
			continue
		}
		err = utils.GitClone(poc_url+".git", poc_dir+cve+"/"+strings.Split(poc_url, "/")[3])
		if err != nil {
			continue
		}
		//countchan <- struct{}{}
		count++
	}
	return count, nil
}

func CheckGithubUpdate() string {
	var info strings.Builder
	updated := false
	year := time.Now().Year()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	poc_dir := cfg.Github.PocDir
	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=CVE-%d&per_page=100&page=1", year)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		utils.PrintColor("error", "Error creating request: %v\n", err)
		return ""
	}
	headers := getGithubHeader()
	utils.SetHeaders(req, headers)

	res, err := client.Do(req)

	if err != nil {
		utils.PrintColor("error", "Error establishing connection: %v\n", err)
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
		existNum, _ := database.CountData(cve)
		if existNum >= 5 {
			continue
		}
		if date_published.Format("2006-01-02") == yesterday {
			updated = true
			tmp := fmt.Sprintf("%v\npoc_url:%v\ndesc:%v\ncvss2:%v\ncvss3:%v\ncna:%v\n\n", cve, poc_url, desc, cvss2, cvss3, cna)
			info.WriteString(tmp)
			err = database.InsertGithubDB(cve, desc, date_published, cvss2, cvss3, cna, poc_url)
			if err != nil {
				utils.PrintColor("error", "Error insert vul:", cve, "\n")
				continue
			}
			err = utils.GitClone(poc_url+".git", poc_dir+cve+"/"+strings.Split(poc_url, "/")[3])
			if err != nil {
				continue
			}
		}
	}
	if updated {
		utils.PrintColor("info", "Github Updated")
		utils.WriteToLog(yesterday+"\n"+info.String(), githubupdateInfoPath)
		return info.String()
	} else {
		utils.WriteToLog(yesterday+"\nAlready up to date.", githubupdateInfoPath)
		return "Already up to date."
	}
}
func GithubCrawler() {
	startYear := cfg.Github.StartYear
	endYear := cfg.Github.EndYear
	totalCount := 0
	for year := startYear; year <= endYear; year++ {
		yearCount := 0
		for page := 1; page <= 10; page++ {
			count, _ := fetchGithub(year, page)
			utils.PrintColor("success", "year %d page %d successfully insert %d vul info", year, page, count)
			yearCount += count
		}
		utils.PrintColor("success", "Year %d successfully insert %d vul info", year, yearCount)
		totalCount += yearCount
	}
	utils.PrintColor("success", "Github: Successfully insert %d vul info", totalCount)
}

/*
func ConcurrentGitHubCrawler() {
	var wg, mainWG sync.WaitGroup
	var totalCount int
	startYear := cfg.Github.StartYear
	endYear := cfg.Github.EndYear
	for year := startYear; year <= endYear; year++ {
		mainWG.Add(1)
		yearCount := 0
		countChan := make(chan struct{})
		for page := 1; page <= 10; page++ {
			wg.Add(1)
			go fetchGithub(year, page, &wg, countChan)
		}

		go func() {
			wg.Wait()
			close(countChan)
			mainWG.Done()
		}()

		for range countChan {
			yearCount++
		}

		utils.PrintColor("success", "Successfully insert %d CVE-%d vul", yearCount, year)
		totalCount += yearCount
	}
	mainWG.Wait()
	utils.PrintColor("success", "Github: Successfully insert %d vul info", totalCount)
}


*/
