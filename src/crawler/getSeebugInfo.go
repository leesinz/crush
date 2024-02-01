package crawler

import (
	"context"
	"crush/database"
	"crush/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SeebugInfo struct {
	SSVID      string
	SubmitTime string
	Severity   string
	Name       string
	CVE        string
	POC        string
}

var (
	seebuglogDir         = filepath.Join(utils.GetParentPath(), "data", "seebug", "log")
	seebugupdateInfoPath = filepath.Join(seebuglogDir, "seebug_update_info.log")
)

func getHTML(page int) (string, error) {
	ctx, _ := chromedp.NewExecAllocator(context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			//chromedp.Flag("headless", true),
			chromedp.Flag("enable-automation", false),
			//chromedp.Flag("disable-gpu", true),
			chromedp.Flag("disable-blink-features", "AutomationControlled"),
			//chromedp.Flag("no-sandbox", true),
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

func string2date(date_str string) time.Time {
	t, _ := time.Parse("2006-01-02", date_str)
	return t
}

func checkPOC(poc string) bool {
	return !strings.Contains(poc, "无")
}

func checkSeverity(severity string) string {
	severity_tmp := strings.TrimSpace(severity)
	switch severity_tmp {
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

func FetchSeebug() {
	logFile, err := os.Create("error.log")
	if err != nil {
		log.Fatalf("failed to create error log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	for i := 1; i <= 2936; i++ {
		log.Printf("starting crawling page %d", i)
		htmlContent, err := getHTML(i)
		if err != nil {
			logger.Println(err)
			continue // 继续下一个页面的爬取
		}

		items := extractHTML(htmlContent)
		for _, item := range items {
			err := database.InsertSeebugDB(item.SSVID, string2date(item.SubmitTime), checkSeverity(item.Severity), item.Name, checkCVE(item.CVE), checkPOC(item.POC))
			if err != nil {
				logger.Println(err)
			}
		}
	}
}

func CheckSeebugUpdate() string {
	var info strings.Builder
	updated := false
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	htmlContent, _ := getHTML(1)
	items := extractHTML(htmlContent)
	for _, item := range items {
		if item.SubmitTime == yesterday {
			updated = true
			database.InsertSeebugDB(item.SSVID, string2date(item.SubmitTime), checkSeverity(item.Severity), item.Name, checkCVE(item.CVE), checkPOC(item.POC))
			tmp := fmt.Sprintf("%v %v %v %v %d\n", item.SSVID, checkSeverity(item.Severity), item.Name, checkCVE(item.CVE), checkPOC(item.POC))
			info.WriteString(tmp)
		}

		//database.InsertSeebugDB(item.SSVID, string2date(item.SubmitTime), checkSeverity(item.Severity), item.Name, checkCVE(item.CVE), checkPOC(item.POC))
	}
	if updated {
		utils.WriteToLog(yesterday+"\n"+info.String(), seebugupdateInfoPath)
		return info.String()
	} else {
		utils.WriteToLog(yesterday+"\nAlready up to date.", seebugupdateInfoPath)
		return "Already up to date."
	}
}
