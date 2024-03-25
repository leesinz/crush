package crawler

import (
	"crush/database"
	"crush/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var (
	packetstormlogDir         = filepath.Join(utils.GetParentPath(), "data", "packetstorm", "log")
	packetstormupdateInfoPath = filepath.Join(packetstormlogDir, "packetstorm_update_info.log")
)

func CheckPacketstormUpdate() string {
	var result, rstHTML strings.Builder
	updated := false
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	url := "https://packetstormsecurity.com/files/tags/exploit/"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("request err:", err)
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	doc.Find("dl[id]").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		name := s.Find("dt a").Text()
		date_selector := s.Find("dd.datetime a[href^='/files/date']")
		date_tmp, _ := date_selector.Attr("href")
		dateParts := strings.Split(date_tmp, "/")
		date := dateParts[len(dateParts)-2]
		cveLinks := s.Find("dd.cve a")

		var cves []string
		cveLinks.Each(func(j int, cveLink *goquery.Selection) {
			cve := cveLink.Text()
			cves = append(cves, cve)
		})

		cveStr := strings.Join(cves, ",")
		description := s.Find("dd.detail p").Text()
		poc_selector := s.Find("dd.act-links a[href^='/files/download']")
		poc_tmp, _ := poc_selector.Attr("href")
		poc := "https://packetstormsecurity.com" + poc_tmp
		if date == yesterday {
			updated = true
			err := database.InsertPacketstormDB(id, string2date(date), name, cveStr, poc, description)
			if err != nil {
				log.Fatal(err)
			}
			result.WriteString(id + " " + name + " " + cveStr + "\n")
			rstHTML.WriteString(fmt.Sprintf("%s   <a href=\"%s\">%s</a>\n", id, poc, name))
		}
	})

	if !updated {
		result.WriteString("Already up to date.")
		rstHTML = result
	} else {
		utils.PrintColor("success", "PacketStorm Updated")
	}
	utils.WriteToLog(yesterday+"\n"+result.String(), packetstormupdateInfoPath)
	return rstHTML.String()
}
