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
	"regexp"
	"strings"
	"time"
)

var (
	ZerodaylogDir         = filepath.Join(utils.GetCurrentPath(), "data", "zeroday")
	ZerodayUpdateInfoPath = filepath.Join(ZerodaylogDir, "zeroday_update_info.log")
)

func Check0dayUpdate() string {
	ctx, _ := chromedp.NewExecAllocator(context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			//chromedp.Flag("headless", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-blink-features", "AutomationControlled"),
		)...,
	)
	ctx, _ = chromedp.NewContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	var content string
	var result strings.Builder
	dateFormat := Yesterday.Format("02-01-2006")
	url := fmt.Sprintf("https://0day.today/date/%s", dateFormat)
	agreeXpath := "/html/body/div/div[1]/div[14]/div[3]/form/input"
	closeXpath := "/html/body/div[5]/div/div/a"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(agreeXpath, chromedp.BySearch),
		chromedp.Click(agreeXpath),
		chromedp.WaitVisible(closeXpath, chromedp.BySearch),
		chromedp.Click(closeXpath),
		chromedp.InnerHTML("html", &content),
	)
	if err != nil {
		errMsg := fmt.Sprintf("crawling %v err:%v", url, err)
		log.Println(errMsg)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	newVul := make(map[string]bool)
	vulnerabilities := []Vulnerability{}
	doc.Find(".ExploitTableContent").Each(func(i int, s *goquery.Selection) {
		dateStr := s.Find(".td a").First().Text()
		dateTmp, _ := time.Parse("02-01-2006", dateStr)
		date := dateTmp.Format("2006-01-02")
		idSelection := s.Find("h3 a[href^='/exploit/description']")
		idTmp, _ := idSelection.Attr("href")
		idParts := strings.Split(idTmp, "/")
		id := idParts[len(idParts)-1]
		name := s.Find("h3 a").Text()
		name = strings.TrimSpace(name)
		category := s.Find("a[href^='/platforms']").Text()
		riskTmp := s.Find("[class*=tips_risk_color_]").Text()
		risk := strings.Split(riskTmp, " ")[len(strings.Split(riskTmp, " "))-1]
		cve := s.Find("a[href^='/cve']").Next().Text()
		if cve != "" {
			re := regexp.MustCompile(`CVE-\d+-\d+`)
			matches := re.FindAllString(strings.ToUpper(cve), -1)
			cve = strings.Join(matches, ",")
		}
		poc := s.Find(".tips_price_0").Text()
		if poc != "" {
			poc = "https://0day.today/exploit/" + id
		}
		if !newVul[name] {
			newVul[name] = true
			err := database.InsertZerodayDB(id, name, string2date(date), category, cve, risk, poc)
			if err != nil {
				log.Fatal(err)
			}
			result.WriteString(fmt.Sprintf("%s  %s\n", name, cve))
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   name,
				CVE:    cve,
				URL:    poc,
				Source: "0day.today",
			})
		}

	})
	if len(newVul) == 0 {
		result.WriteString("Already up to date.")
		utils.PrintLog("info", "0day.today is up to date")
	} else {
		utils.PrintLog("success", "0day.today update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	}
	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), ZerodayUpdateInfoPath)
	return result.String()
}
