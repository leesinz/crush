package crawler

import (
	"crush/database"
	"crush/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	packetstormlogDir         = filepath.Join(utils.GetParentPath(), "data", "packetstorm", "log")
	packetstormupdateInfoPath = filepath.Join(packetstormlogDir, "packetstorm_update_info.log")
	packetstormDir            = cfg.PacketStorm.PacketstormDir
)

func DownloadPacketstorm(url string, filePath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func CheckPacketstormUpdate() string {
	var result strings.Builder
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
		cve := s.Find("dd.cve a").Text()
		description := s.Find("dd.detail p").Text()
		poc_selector := s.Find("dd.act-links a[href^='/files/download']")
		poc_tmp, _ := poc_selector.Attr("href")
		poc := "https://packetstormsecurity.com" + poc_tmp
		pocParts := strings.Split(poc, "/")
		poc_name := pocParts[len(pocParts)-1]
		if date == yesterday {
			updated = true
			err := database.InsertPacketstormDB(id, string2date(date), name, cve, poc, description)
			if err != nil {
				log.Fatal(err)
			}
			err = DownloadPacketstorm(poc, packetstormDir+id+"/"+poc_name)
			if err != nil {
				log.Fatal(err)
			}
			result.WriteString(id + " " + name + " " + cve + "\n")
		}
	})

	if !updated {
		result.WriteString("Already up to date.")
	}

	utils.WriteToLog(yesterday+"\n"+result.String(), packetstormupdateInfoPath)
	return result.String()
}
