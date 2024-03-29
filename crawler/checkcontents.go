package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var accessToken = cfg.Github.GithubToken

type File struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func getFileList(auther, repository, path string) ([]File, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", auther, repository, path)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files []File
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

func GetTargetFiles(auther, repository, path, suffix string) ([]string, error) {
	files, err := getFileList(auther, repository, path)
	if err != nil {
		return nil, err
	}

	var targetFiles []string
	for _, file := range files {
		if file.Type == "file" && strings.HasSuffix(file.Name, suffix) {
			targetFiles = append(targetFiles, path+"/"+file.Name)
		} else if file.Type == "dir" {
			subDirFiles, err := GetTargetFiles(auther, repository, path+"/"+file.Name, suffix)
			if err != nil {
				return nil, err
			}
			targetFiles = append(targetFiles, subDirFiles...)
		}
	}
	return targetFiles, nil
}
