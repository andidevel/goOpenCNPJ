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

package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/catfishlabs/goOpenCNPJ/consts"
	"github.com/catfishlabs/goOpenCNPJ/importer"
	"github.com/catfishlabs/goOpenCNPJ/model"
	"github.com/catfishlabs/goOpenCNPJ/scraping"
	"github.com/catfishlabs/goOpenCNPJ/utils"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

const Version = "0.0.1"

type threadStatus struct {
	err        error
	threadInfo string
}

type DownloadAction struct {
	md                model.IDataStorage
	ws                *scraping.CNPJDataScrape
	downloadTo        string
	companySchemaFile string
	companyConf       utils.CompanyDownloadConfig
	ci                *importer.CompanyImporter
}

func NewDownloadAction(dataEnv map[string]string, companiesConfFile, companiesSchemaFile string, md model.IDataStorage) (*DownloadAction, error) {
	var err error
	result := DownloadAction{
		md:                md,
		downloadTo:        dataEnv["DATA_DOWNLOAD_PATH"],
		companySchemaFile: companiesSchemaFile,
	}
	result.ci = importer.NewCompanyImporter(companiesSchemaFile, md)
	// Load companies config file
	result.companyConf, err = utils.LoadCompanyDownloadConfig(companiesConfFile)
	if err != nil {
		return &result, err
	}
	result.ws = scraping.New(result.companyConf.CompaniesMainUrl)

	return &result, err
}

func (da *DownloadAction) importCitiesFromCSV(csvFileName string) error {
	csvFileDownloaded := filepath.Join(da.downloadTo, csvFileName)
	return importer.CitiesFromCSV(csvFileDownloaded, da.md)
}

func (da *DownloadAction) downloadAndImportCitiesFile() error {
	err := utils.FileDownload(da.ws.CitiesFile, da.downloadTo)
	if err == nil {
		// Cities file is a zip file
		csvZipFileName := filepath.Join(da.downloadTo, filepath.Base(da.ws.CitiesFile))
		csvFileName, err := firstZipFile(csvZipFileName)
		if err == nil {
			utils.Unzip(csvZipFileName, da.downloadTo)
			return da.importCitiesFromCSV(csvFileName)
		}

	}
	return err
}

func (da *DownloadAction) importStatusFromCSV(csvFileName string) error {
	// Import File
	csvFileDownloaded := filepath.Join(da.downloadTo, csvFileName)
	return importer.StatusFromCSV(csvFileDownloaded, da.md)
}

func (da *DownloadAction) downloadAndImportStatusFile() error {
	// Download status file
	err := utils.FileDownload(da.ws.StatusFile, da.downloadTo)
	if err == nil {
		// Status file is a zip file
		csvZipFileName := filepath.Join(da.downloadTo, filepath.Base(da.ws.StatusFile))
		csvFileName, err := firstZipFile(csvZipFileName)
		if err == nil {
			utils.Unzip(csvZipFileName, da.downloadTo)
			return da.importStatusFromCSV(csvFileName)
		}

	}
	return err
}

func (da *DownloadAction) importNR04FromPDF(pdfFileName string) error {
	nr04FileDownloaded := filepath.Join(da.downloadTo, pdfFileName)
	return importer.NR04FromPDF(nr04FileDownloaded, da.md)
}

func (da *DownloadAction) downloadAnImportNR04File() error {
	var err error
	nr04FileName := filepath.Base(da.companyConf.NR04Url)
	err = utils.FileDownload(da.companyConf.NR04Url, da.downloadTo)
	if err == nil {
		return da.importNR04FromPDF(nr04FileName)
	}
	return err
}

// Download and import auxiliary tables
func (da *DownloadAction) auxiliaryTables() {
	auxTablesFunc := []func(chan<- threadStatus){
		func(c chan<- threadStatus) {
			log.Printf("Importing [%s]...\n", da.ws.StatusFile)
			ts := threadStatus{}
			err := da.downloadAndImportStatusFile()
			ts.err = err
			ts.threadInfo = "Importing Status from CSV file"
			c <- ts
		},
		func(c chan<- threadStatus) {
			ts := threadStatus{}
			// Supposed to be a PDF. Extract text and import
			err := da.downloadAnImportNR04File()
			ts.err = err
			ts.threadInfo = "Importing NR04 from PDF file"
			c <- ts
		},
		func(c chan<- threadStatus) {
			log.Printf("Importing [%s]...\n", da.ws.CitiesFile)
			ts := threadStatus{}
			err := da.downloadAndImportCitiesFile()
			ts.err = err
			ts.threadInfo = "Importing Cities from CSV file"
			c <- ts
		},
	}
	fq := make(chan threadStatus)
	for _, f := range auxTablesFunc {
		go f(fq)
	}

	// Block process, we need these auxiliary tables before go any further
	for range auxTablesFunc {
		ts := <-fq
		log.Print(ts.threadInfo)
		if ts.err != nil {
			log.Println("...Error:", ts.err)
		} else {
			log.Println("...Done!")
		}
	}
}

func (da *DownloadAction) downloadAndUnzipOneFile(c chan<- threadStatus, fileURL string) {
	ts := threadStatus{}
	var err error
	baseURL, dataFile := path.Split(fileURL)
	urlsToTry := []string{baseURL}
	urlsToTry = append(urlsToTry, da.companyConf.CompaniesMirrorUrls...)
	canDownload := false
	// Check if file URL is a downloadable zip file. If not, try fallback mirror
	// Typical value for "Content-Type" key = "application/zip", "application/zip; field=value"
	for _, u := range urlsToTry {
		urlObj, err := url.Parse(u)
		if err == nil {
			urlObj.Path = path.Join(urlObj.Path, dataFile)
			fileURL = urlObj.String()
			fileHeader, err := utils.HeadersFrom(fileURL)
			if err == nil {
				if checkHeader(fileHeader, "Content-Type", "application/zip") {
					canDownload = true
					break
				}
			}
		}
	}
	if canDownload {
		log.Printf(" |-> Downloading %s\n", fileURL)
		err = utils.FileDownload(fileURL, da.downloadTo)
		if err == nil {
			// Unzip downloaded file
			zipFile := filepath.Join(da.downloadTo, dataFile)
			// We are expecting a one CSV ziped file
			var csvFileName string
			csvFileName, err = firstZipFile(zipFile)
			if err == nil {
				// Unzip
				log.Println(" |-> Unziping:", csvFileName)
				utils.Unzip(zipFile, da.downloadTo)
				err = da.ci.CompaniesFromCSV(filepath.Join(da.downloadTo, csvFileName))
			}
		}
	} else {
		err = fmt.Errorf("not an application/zip file [%s]", fileURL)
	}
	ts.err = err
	ts.threadInfo = fmt.Sprintf("Data File: %s", dataFile)
	c <- ts
}

func (da *DownloadAction) downloadAll(forceDownload bool, n int) {
	canProcess := forceDownload
	da.ws.GetCNPJData()
	dtUpdated, err := time.Parse(consts.DateLayoutBR, da.ws.LastUpdate)
	if err != nil {
		log.Fatal("Error parsing updated date:", err)
	}
	if !forceDownload {
		canProcess = isTimeToUpdate(da.md, dtUpdated)
	}
	log.Println("Updated (from Federal Revenue site):", da.ws.LastUpdate)
	if canProcess {
		log.Println("Status descriptions file:", da.ws.StatusFile)
		log.Printf("Time to update (last update: %s). This can take a while!!!\n", da.ws.LastUpdate)
		log.Println("First, update auxiliary tables:")
		da.auxiliaryTables()
		log.Println("Now, the main files:")
		fq := make(chan threadStatus)
		dataFiles := da.ws.DataFiles
		if n > 0 {
			// Now there are two data files with data about the same company: *EMPRECSV* and *ESTABELE*
			dataFiles = []string{}
			n0 := 0
			n1 := 0
			for _, f := range da.ws.DataFiles {
				schemaType := importer.GetSchemaTypeByName(f)
				if schemaType == "0" && n0 < n {
					dataFiles = append(dataFiles, f)
					n0++
				} else {
					if schemaType == "1" && n1 < n {
						dataFiles = append(dataFiles, f)
						n1++
					}
				}
			}
		}
		for _, f := range dataFiles {
			go da.downloadAndUnzipOneFile(fq, f)
		}
		// Wait file processing or errors
		for range dataFiles {
			ts := <-fq
			log.Print(ts.threadInfo)
			if ts.err != nil {
				log.Println("...Error:", err)
			} else {
				log.Println("...Downloaded and parsed!")
			}
		}
		log.Println("Done!?!")
	} else {
		log.Println("Not yet! Last time was", da.ws.LastUpdate)
	}
}

func dateTimeToTime(dt primitive.DateTime) time.Time {
	return dt.Time().UTC()
}

// firstZipFile tries to get the name of first file found in a file ziped
// TODO: Is there a better way to know first file name in a Zipfile?
func firstZipFile(zipFile string) (string, error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			// Skip until find first file
			continue
		}
		return f.Name, nil
	}
	return "", errors.New("no files in first level")
}

func checkHeader(headers http.Header, headerKey, headerValue string) bool {
	for _, v := range headers.Values(headerKey) {
		if v != "" && strings.Contains(v, headerValue) {
			return true
		}
	}
	return false
}

func isTimeToUpdate(md model.IDataStorage, dt time.Time) bool {
	result := false
	param := model.Parameter{
		ID:    "cnpj.update.date",
		Value: dt,
	}
	updateParam, err := md.FindOneUpsertParameter(param)
	if err != nil && err != model.ErrNoRows {
		log.Println("Error finding parameter:", err)
	}
	if err == model.ErrNoRows {
		// Fist time, download anyway
		result = true
	} else {
		// Using assert to convert value. Is that a best pratice??
		decodedDate := dateTimeToTime(updateParam.Value.(primitive.DateTime))
		result = dt.After(decodedDate)
	}
	return result
}

func main() {
	app := &cli.App{
		Name:  fmt.Sprintf("get-companies v%s", Version),
		Usage: "download and import Brazil's companies database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "config/companies-download.json",
				Aliases: []string{"c"},
				Usage:   "path to companies JSON config file",
			},
			&cli.StringFlag{
				Name:    "schema",
				Value:   "config/cnpj-schema.json",
				Aliases: []string{"s"},
				Usage:   "path to companies JSON schema file",
			},
			&cli.BoolFlag{
				Name:    "aux",
				Value:   false,
				Aliases: []string{"a"},
				Usage:   "download and parse auxiliary tables only",
			},
			&cli.BoolFlag{
				Name:    "force",
				Value:   false,
				Aliases: []string{"f"},
				Usage:   "force download and parse",
			},
			&cli.Int64Flag{
				Name:    "nfiles",
				Value:   0,
				Aliases: []string{"n"},
				Usage:   "number of data files to download",
			},
		},
		Action: func(c *cli.Context) error {
			// Load env config
			envConfig, err := godotenv.Read()
			if err != nil {
				log.Fatal("Error reading configuration:", err)
			}
			// Database connect
			md := model.NewMongoDatabase(envConfig["DBURI"])
			err = md.Connect()
			if err != nil {
				return err
			}
			defer md.Close()

			da, err := NewDownloadAction(envConfig, c.String("config"), c.String("schema"), md)
			if c.Bool("aux") {
				da.ws.GetCNPJData()
				da.auxiliaryTables()
				return nil
			}
			da.downloadAll(c.Bool("force"), int(c.Int64("nfiles")))
			return err
		},
	}

	// Parse command line arguments
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
