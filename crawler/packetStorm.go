package crawler

import (
	"crush/database"
	"crush/utils"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	PacketstormlogDir         = filepath.Join(utils.GetCurrentPath(), "data", "packetstorm")
	PacketstormupdateInfoPath = filepath.Join(PacketstormlogDir, "packetstorm_update_info.log")
	PacketstormPocDir         = cfg.PacketStorm.PacketstormPocDir
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
	vulnerabilities := []Vulnerability{}
	doc.Find("dl[id]").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		name := s.Find("dt a").Text()
		dateSelector := s.Find("dd.datetime a[href^='/files/date']")
		dateTmp, _ := dateSelector.Attr("href")
		dateParts := strings.Split(dateTmp, "/")
		date := dateParts[len(dateParts)-2]

		cveLinks := s.Find("dd.cve a")
		var cves []string
		cveLinks.Each(func(j int, cveLink *goquery.Selection) {
			cve := cveLink.Text()
			cves = append(cves, cve)
		})
		cveStr := strings.Join(cves, ",")

		description := s.Find("dd.detail p").Text()
		pocSelector := s.Find("dd.act-links a[href^='/files/download']")
		pocTmp, _ := pocSelector.Attr("href")
		poc := "https://packetstormsecurity.com" + pocTmp
		pocParts := strings.Split(poc, "/")
		pocName := pocParts[len(pocParts)-1]
		if date == Yesterday.Format("2006-01-02") {
			updated = true

			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:   name,
				CVE:    cveStr,
				URL:    poc,
				Source: "packetstorm",
			})

			err := database.InsertPacketstormDB(id, string2date(date), name, cveStr, poc, description)
			if err != nil {
				log.Fatal(err)
			}

			if DownloadPOC == true {
				err = DownloadPacketstorm(poc, PacketstormPocDir+id+"/"+pocName)
				if err != nil {
					log.Fatal(err)
				}
			}
			result.WriteString(fmt.Sprintf("%s %s\n", name, cveStr))
		}
	})

	if !updated {
		utils.PrintLog("info", "packetstorm is up to date")
		result.WriteString("Already up to date.")
	} else {
		utils.PrintLog("success", "packetstorm update")
		jsonData, _ := json.MarshalIndent(vulnerabilities, "", "    ")
		utils.WriteToLog(string(jsonData), JsonlogPath)
	}

	utils.WriteToLog(Yesterday.Format("2006-01-02")+"\n"+result.String(), PacketstormupdateInfoPath)
	return result.String()
}
