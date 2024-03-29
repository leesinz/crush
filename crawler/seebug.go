package crawler

import (
	"context"
	"crush/database"
	"crush/utils"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"path/filepath"
	"strings"
	"time"
)

var (
	SeebuglogDir         = filepath.Join(utils.GetCurrentPath(), "data", "seebug")
	SeebugupdateInfoPath = filepath.Join(SeebuglogDir, "seebug_update_info.log")
)

type SeebugInfo struct {
	SSVID      string
	SubmitTime string
	Severity   string
	Name       string
	CVE        string
	POC        string
}

func getHTML(page int) (string, error) {
	ctx, _ := chromedp.NewExecAllocator(context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-blink-features", "AutomationControlled"),
		)...,
	)
	ctx, _ = chromedp.NewContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	var result string
	url := fmt.Sprintf("https://www.seebug.org/vuldb/vulnerabilities?page=%d", page)
	xpath := "/html/body/div[2]/div/div/div/div/table/tbody/tr[1]/td[4]/a"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(xpath, chromedp.BySearch),
		chromedp.InnerHTML("html", &result),
	)
	if err != nil {
		errMsg := fmt.Sprintf("crawling %v err:%v", url, err)
		log.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}
	return result, nil
}

func extractHTML(html string) map[string]SeebugInfo {
	infomap := make(map[string]SeebugInfo)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		ssvID := s.Find("td").Eq(0).Text()
		submitTime := s.Find("td").Eq(1).Text()
		severity := s.Find("td").Eq(2).Find("div").AttrOr("data-original-title", "")
		vulName := s.Find("td").Eq(3).Find("a").Text()
		cve := s.Find("td").Eq(4).Find("i").First().AttrOr("data-original-title", "")
		poc := s.Find("td").Eq(4).Find("i").Eq(1).AttrOr("data-original-title", "")

		seebugInfo := SeebugInfo{
			ssvID,
			submitTime,
			severity,
			vulName,
			cve,
			poc,
		}
		infomap[ssvID] = seebugInfo
	})
	return infomap
}

func string2date(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

func checkPOC(poc string) bool {
	return !strings.Contains(poc, "无")
}

func checkSeverity(severity string) string {
	severityTmp := strings.TrimSpace(severity)
	switch severityTmp {
	case "低危":
		return "low"
	case "中危":
		return "medium"
	case "高危":
		return "high"
	default:
		return "err severity"
	}
}

func checkCVE(cve string) string {
	if strings.Contains(cve, "无") {
		return ""
	}
	return cve
}

func CheckSeebugUpdate() string {
	var result strings.Builder
	updated := false
	htmlContent, _ := getHTML(1)
	items := extractHTML(htmlContent)
	var vulnerabilities []Vulnerability
	for _, item := range items {
		if item.SubmitTime == Yesterday.Format("2006-01-02") {

			updated = true
			database.InsertSeebugDB(item.SSVID, string2date(item.SubmitTime), checkSeverity(item.Severity), item.Name, checkCVE(item.CVE), checkPOC(item.POC))

			url := "https://www.seebug.org/vuldb/ssvid-" + strings.Split(item.SSVID, "-")[1]
			result.WriteString(fmt.Sprintf("%s %s\n", item.Name, checkCVE(item.CVE)))
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   item.Name,
				CVE:    checkCVE(item.CVE),
				URL:    url,
				Source: "seebug",
			})

		}
	}
	if updated {
		utils.PrintLog("success", "seebug update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	} else {
		utils.PrintLog("info", "seebug is up to date")
		result.WriteString("Already up to date.")
	}

	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), SeebugupdateInfoPath)
	return result.String()
}
