package crawler

import (
	"crush/config"
	"crush/utils"
	"path/filepath"
	"time"
)

type Vulnerability struct {
	Name   string `json:"name"`
	CVE    string `json:"cve"`
	URL    string `json:"url"`
	Source string `json:"source"`
}

var (
	cfg         = config.LoadConfig()
	Yesterday   = time.Now().AddDate(0, 0, -1)
	DownloadPOC = cfg.POC.DownloadPOC
	JsonlogDir  = filepath.Join(utils.GetCurrentPath(), "data", "jsonlog")
	JsonlogPath = filepath.Join(JsonlogDir, Yesterday.Format("2006-01-02"))
)
