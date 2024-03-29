package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const nvdURL = "https://nvd.nist.gov/vuln/detail/"

func GetCurrentPath() string {
	dir, _ := filepath.Abs(".")
	return dir
}

func ConvertToString(slice []interface{}) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprint(v)
	}
	return strings.Join(strSlice, ",")
}

func SetHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func ExtractCVE(name string) string {
	pattern := `(CVE-\d+-\d+)`
	re := regexp.MustCompile(pattern)
	return re.FindString(strings.ToUpper(name))
}

func extractScore(doc *goquery.Document, selector string) string {
	scoreText := doc.Find(selector).Text()
	if scoreText == "" {
		return "0.0"
	}
	return strings.Fields(scoreText)[0]
}

func GetCVEInfo(cve string) (exists bool, desc, cvss2, cvss3, cna string, err error) {
	url := nvdURL + cve
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error establishing connect:", err)
		return false, "", "", "", "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Error parsing html:", err)
		return false, "", "", "", "", err
	}
	exists = doc.Find("#vulnDetailPanel").Length() > 0
	cvss3 = extractScore(doc, "[data-testid=vuln-cvss3-panel-score]")
	cna = extractScore(doc, "[data-testid=vuln-cvss3-cna-panel-score]")
	cvss2 = extractScore(doc, "span.severityDetail a#Cvss2CalculatorAnchor")
	desc = doc.Find("[data-testid=vuln-description]").Text()

	return exists, desc, cvss2, cvss3, cna, nil
}
func GitClone(url, path string) error {
	cmd := exec.Command("git", "clone", url, path)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
