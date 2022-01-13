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

package importer

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/catfishlabs/goOpenCNPJ/model"
	"github.com/catfishlabs/goOpenCNPJ/utils"
)

func CitiesFromCSV(csvFileName string, md model.IDataStorage) error {
	f, err := os.Open(csvFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ID, err := strconv.ParseInt(row[0], 10, 64)
		if err == nil {
			ct := model.City{
				ID:            ID,
				NomeMunicipio: strings.TrimSpace(row[1]),
			}
			_, err := md.FindOneUpsertCity(ct)
			if err != nil && err != model.ErrNoRows {
				log.Println("Error inserting/update Cities table:", err)
			}
		}
	}
	return err
}

func StatusFromCSV(csvFileName string, md model.IDataStorage) error {
	f, err := os.Open(csvFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// It's known CSV file is ISO-8859-15 encoded
	// 2021-07-29: Not anymore. The format changed and now is UTF-8
	// csvFile := transform.NewReader(f, charmap.ISO8859_15.NewDecoder())
	// csvReader := csv.NewReader(csvFile)
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ID, err := strconv.ParseInt(row[0], 10, 64)
		if err == nil {
			sd := model.StatusDescription{
				ID:     ID,
				Motivo: strings.TrimSpace(row[1]),
			}
			_, err := md.FindOneUpsertStatusDescription(sd)
			if err != nil && err != model.ErrNoRows {
				log.Println("Error inserting/updating Status table:", err)
			}
		}
	}
	return err
}

func NR04FromPDF(pdfFile string, md model.IDataStorage) error {
	pdfTexts, err := utils.PDFText(pdfFile)
	if err != nil {
		return err
	}
	cnaePattern, _ := regexp.Compile(`([0-9]{2}\.[0-9]{2}\-[0-9]{1}).*([0-9]){1}`)
	for _, texts := range pdfTexts {
		for _, text := range texts {
			var s strings.Builder
			for _, word := range text.Content {
				s.WriteString(word.S)
			}
			match := cnaePattern.FindStringSubmatch(s.String())
			if len(match) > 0 {
				cnae := utils.RemoveChars(strings.TrimSpace(match[1]), ".-")
				risk_level := strings.TrimSpace(match[2])
				// fmt.Printf("CNAE: %s [%s], Risk Level: %s\n", match[1], cnae, risk_level)
				rl := model.RiskLevel{
					ID:        cnae,
					GrauRisco: risk_level,
				}
				_, err := md.FindOneUpsertRiskLevel(rl)
				if err != nil && err != model.ErrNoRows {
					log.Println("Error inseting/updating Risk Level table:", err)
				}
			}
		}
	}
	return err
}
