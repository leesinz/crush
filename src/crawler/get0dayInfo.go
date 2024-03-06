package crawler

import (
	"context"
	"crush/database"
	"crush/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

var (
	zerodaylogDir         = filepath.Join(utils.GetParentPath(), "data", "zeroday", "log")
	zerodayupdateInfoPath = filepath.Join(zerodaylogDir, "zeroday_update_info.log")
)

func Check0dayUpdate() string {
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
	var content string
	var result, rstHTML strings.Builder
	yesterday := time.Now().AddDate(0, 0, -1)
	date_format := yesterday.Format("02-01-2006")
	url := fmt.Sprintf("https://0day.today/date/%s", date_format)
	agree_xpath := "/html/body/div/div[1]/div[14]/div[3]/form/input"
	close_xpath := "/html/body/div[5]/div/div/a"
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(agree_xpath, chromedp.BySearch),
		chromedp.Click(agree_xpath),
		chromedp.WaitVisible(close_xpath, chromedp.BySearch),
		chromedp.Click(close_xpath),
		chromedp.InnerHTML("html", &content),
	)
	if err != nil {
		errMsg := fmt.Sprintf("crawling %v err:%v", url, err)
		log.Fatal(errMsg)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	vulnerabilities := make(map[string]bool)
	doc.Find(".ExploitTableContent").Each(func(i int, s *goquery.Selection) {
		date_str := s.Find(".td a").First().Text()
		date_tmp, _ := time.Parse("02-01-2006", date_str)
		date := date_tmp.Format("2006-01-02")
		idSelection := s.Find("h3 a[href^='/exploit/description']")
		id_tmp, _ := idSelection.Attr("href")
		idParts := strings.Split(id_tmp, "/")
		id := idParts[len(idParts)-1]
		name := s.Find("h3 a").Text()
		category := s.Find("a[href^='/platforms']").Text()
		risk_tmp := s.Find("[class*=tips_risk_color_]").Text()
		risk := strings.Split(risk_tmp, " ")[len(strings.Split(risk_tmp, " "))-1]
		cve := s.Find("a[href^='/cve']").Next().Text()
		if cve != "" {
			cve = strings.ReplaceAll(strings.TrimSpace(cve), "\n", "")
		} else {
			cve = ""
		}
		poc := s.Find(".tips_price_0").Text()
		if poc != "" {
			poc = "https://0day.today/exploit/" + id
		}
		if !vulnerabilities[name] {
			vulnerabilities[name] = true
			err := database.InsertZerodayDB(id, name, string2date(date), category, cve, risk, poc)
			if err != nil {
				log.Fatal(err)
			}
			result.WriteString(id + "  " + name + "  " + cve + "\n")
			rstHTML.WriteString(fmt.Sprintf("%s   <a href=\"%s\">%s</a>\n", id, poc, name))

		}

	})

	if len(vulnerabilities) == 0 {
		result.WriteString("Already up to date.")
		rstHTML = result
	} else {
		utils.PrintColor("success", "0day.today Updated")
	}
	utils.WriteToLog(yesterday.Format("2006-01-02")+"\n"+result.String(), zerodayupdateInfoPath)
	return rstHTML.String()
}
