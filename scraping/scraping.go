//    Copyright 2021 Anderson Rodrigues do Livramento

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package scraping

import (
	"net/http"
	"regexp"

	"github.com/gocolly/colly/v2"
)

// CNPJDataScrape is a struct to return a list of URL of the files to download and updated date
type CNPJDataScrape struct {
	DataFiles    []string
	StatusFile   string
	CitiesFile   string
	LastUpdate   string
	url          string
	scrapeParser *colly.Collector
}

var (
	cnpjFileNameD1ER, _ = regexp.Compile(`.*EMPRECSV\.zip`)
	cnpjFileNameD2ER, _ = regexp.Compile(`.*ESTABELE\.zip`)
	cnpjLastDateER, _   = regexp.Compile(`Data.*:.*([0-9]{2}/[0-9]{2}/[0-9]{4})`)
	cnpjStatusFile, _   = regexp.Compile(`.*MOTICSV\.zip`)
	citiesFile, _       = regexp.Compile(`.*MUNICCSV\.zip`)
)

// New instantiate a default CNPJDataScrape object
func New(url string) *CNPJDataScrape {
	cds := &CNPJDataScrape{}
	cds.scrapeParser = colly.NewCollector()
	cds.url = url
	return cds
}

// NewLocalFile instantiate a CNPJDataScrape object to scrape in local html files
func NewLocalFile(filePath string) *CNPJDataScrape {
	cds := &CNPJDataScrape{}
	tp := &http.Transport{}
	tp.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	cds.scrapeParser = colly.NewCollector()
	cds.scrapeParser.WithTransport(tp)
	cds.url = "file://" + filePath
	return cds
}

// GetCNPJDataFiles scrape an HTML document for CNPJ files
func (cds *CNPJDataScrape) getCNPJDataFiles() {
	cds.scrapeParser.OnHTML("div[id=\"parent-fieldname-text\"]", func(e *colly.HTMLElement) {
		for _, href := range e.ChildAttrs("a.external-link[href]", "href") {
			// CNPJ data file
			// fmt.Println("Parsed href:", href)
			if cnpjFileNameD1ER.MatchString(href) || cnpjFileNameD2ER.MatchString(href) {
				cds.DataFiles = append(cds.DataFiles, href)
			} else {
				if cds.StatusFile == "" && cnpjStatusFile.MatchString(href) {
					cds.StatusFile = href
				}
				if cds.CitiesFile == "" && citiesFile.MatchString(href) {
					cds.CitiesFile = href
				}
			}
		}
		// for _, href := range e.ChildAttrs("a.internal-link[href]", "href") {
		// 	// Status descriptions file
		// 	// fmt.Println("Parsed href:", href)
		// 	if cds.StatusFile == "" && cnpjStatusFile.MatchString(href) {
		// 		cds.StatusFile = href
		// 	}
		// }
		for _, pText := range e.ChildTexts("p") {
			// fmt.Println("P:", pText)
			match := cnpjLastDateER.FindStringSubmatch(pText)
			if len(match) > 0 {
				cds.LastUpdate = match[1]
			}
		}
	})
	cds.scrapeParser.Visit(cds.url)
}

// GetCNPJData gets CNPJ data from an URL
func (cds *CNPJDataScrape) GetCNPJData() {
	cds.getCNPJDataFiles()
}
