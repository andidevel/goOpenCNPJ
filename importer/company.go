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
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/catfishlabs/goOpenCNPJ/consts"
	"github.com/catfishlabs/goOpenCNPJ/model"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type fieldConvFunc func(string) interface{}

var converter = map[string]fieldConvFunc{
	"float": func(v string) interface{} {
		// v is a string like "067000000000,00", replace "," by ".", so s = "067000000000.00"
		// l := len(v)
		// s := v[:l-2] + "." + v[l-2:]
		s := strings.ReplaceAll(v, ",", ".")
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
		return 0.0
	},
	"int": func(v string) interface{} {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
		return nil
	},
	"timestamp": func(v string) interface{} {
		if v != "" {
			if t, err := time.Parse(consts.DateLayoutSchema, v); err == nil {
				return model.DateTime(t)
			}
		}
		return nil
	},
}

// CNPJFieldMap maps a field into JSON layout object
type CNPJFieldMap struct {
	FieldType string `json:"field_type"`
	Position  int    `json:"position"`
}

// CNPJLayoutJSONMap maps a JSON object representing a layout configuration
type CNPJLayoutJSONMap struct {
	Type     string                  `json:"type"`
	Document map[string]CNPJFieldMap `json:"document"`
}

type CompanyImporter struct {
	layoutJSONFile string
	md             model.IDataStorage
}

func GetSchemaTypeByName(csvFileName string) string {
	if consts.SchemaType0.MatchString(csvFileName) {
		return "0"
	}
	if consts.SchemaType1.MatchString(csvFileName) {
		return "1"
	}
	return ""
}

func NewCompanyImporter(layoutJSONFile string, md model.IDataStorage) *CompanyImporter {
	ci := CompanyImporter{
		layoutJSONFile: layoutJSONFile,
		md:             md,
	}
	return &ci
}

// loadLayoutSchema load a json file with layout map schema
func (ci *CompanyImporter) loadLayoutSchema() ([]CNPJLayoutJSONMap, error) {
	jsonFile, err := os.Open(ci.layoutJSONFile)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	var layout []CNPJLayoutJSONMap
	err = json.NewDecoder(jsonFile).Decode(&layout)
	return layout, err
}

func (ci *CompanyImporter) findMapType(mapType string, layout []CNPJLayoutJSONMap) map[string]CNPJFieldMap {
	for _, v := range layout {
		if v.Type == mapType {
			return v.Document
		}
	}
	return nil
}

func (ci *CompanyImporter) saveCompany(co model.Company) {
	_, err := ci.md.FindOneUpsertCompany(co)
	if err != nil && err != model.ErrNoRows {
		log.Println("Error inserting/updating Company table:", err)
	}
}

func (ci *CompanyImporter) saveBaseCompany(bc model.BaseCompany) {
	_, err := ci.md.FindOneUpsertBaseCompany(bc)
	if err != nil && err != model.ErrNoRows {
		log.Println("Error inserting/updating BaseCompany table:", err)
	}
}

func mapFromSchema(row []string, schema map[string]CNPJFieldMap) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range schema {
		value := row[v.Position]
		fn, keyExists := converter[v.FieldType]
		if keyExists {
			result[k] = fn(value)
		} else {
			result[k] = value
		}
	}
	return result
}
func (ci *CompanyImporter) CompaniesFromCSV(csvFileName string) error {
	var err error
	fCSV, err := os.Open(csvFileName)
	if err != nil {
		return err
	}
	defer fCSV.Close()

	companySchema, err := ci.loadLayoutSchema()
	if err != nil {
		return err
	}
	// It's known CSV file is ISO-8859-15 encoded
	csvFileCharEncoded := transform.NewReader(fCSV, charmap.ISO8859_15.NewDecoder())
	// csvReader := csv.NewReader(fCSV)
	csvReader := csv.NewReader(csvFileCharEncoded)
	csvReader.Comma = ';'
	mapType := GetSchemaTypeByName(csvFileName)
	schemaType := ci.findMapType(mapType, companySchema)
	for {
		row, read_err := csvReader.Read()
		if read_err == io.EOF {
			break
		}
		doc := mapFromSchema(row, schemaType)
		switch mapType {
		case "0":
			// Remove CPF from "razao_social" if personal company
			if doc["codigo_natureza_juridica"].(int64) == 2135 {
				cpfMatch := consts.CPFMEIER.FindAllStringSubmatch(doc["razao_social"].(string), -1)
				for i := 0; i < len(cpfMatch); i++ {
					doc["razao_social"] = strings.Trim(strings.ReplaceAll(doc["razao_social"].(string), cpfMatch[i][1], ""), " ")
				}
			}
			var bc model.BaseCompany
			model.DecodeFromMap(doc, &bc)
			ci.saveBaseCompany(bc)

		case "1":
			doc["_id"] = strings.Join(
				[]string{
					doc["empresa_base_id"].(string),
					doc["id_ordem"].(string),
					doc["id_dv"].(string),
				},
				"",
			)
			doc["cnaes_secundarios"] = strings.Split(doc["cnaes_secundarios"].(string), ",")
			doc["motivo_situacao_cadastral"] = ""
			doc["grau_risco"] = ""
			status, err := ci.md.FindOneStatusDescriptionById(doc["codigo_situacao_cadastral"].(int64))
			if err == nil {
				doc["motivo_situacao_cadastral"] = status.Motivo
			}
			if rl, ok := doc["cnae_fiscal"]; ok && rl != "" {
				riskLevel, err := ci.md.FindOneRiskLevelById(rl.(string)[:5])
				if err == nil {
					doc["grau_risco"] = riskLevel.GrauRisco
				}
			}
			if ct, ok := doc["codigo_municipio"]; ok {
				city, err := ci.md.FindOneCityById(ct.(int64))
				if err == nil {
					doc["nome_municipio"] = city.NomeMunicipio
				}
			}
			var co model.Company
			model.DecodeFromMap(doc, &co)
			ci.saveCompany(co)
		}
	}
	return nil
}

// DEPRECATED - This format isn't used by Federal Revenue anymore
// func (ci *CompanyImporter) getCNAEs(s string) []string {
// 	var cnaes []string
// 	first := 0
// 	last := 7
// 	length := 7
// 	if (len(s) % length) == 0 {
// 		for i := 0; i < (len(s) / length); i++ {
// 			cnae := s[first:last]
// 			if cnae != "0000000" {
// 				cnaes = append(cnaes, cnae)
// 			}
// 			first = first + length
// 			last = last + length
// 		}
// 	}
// 	return cnaes
// }

// DEPRECATED - This format isn't used by Federal Revenue anymore
// func (ci *CompanyImporter) processLine(lineP string, documentMap map[string]CNPJFieldMap) map[string]interface{} {
// 	result := map[string]interface{}{}
// 	for k, v := range documentMap {
// 		// Convert into proper type !?!?!
// 		value := strings.Trim(lineP[v.Position[0]-1:v.Position[1]], " ")
// 		fn, keyExists := converter[v.FieldType]
// 		if keyExists {
// 			result[k] = fn(value)
// 		} else {
// 			result[k] = value
// 		}
// 	}
// 	return result
// }

// CompaniesFromFixedWidthFile imports company data from a standardized fixed width file, following the standard in docs/LAYOUT_DADOS_ABERTOS_CNPJ.pdf
// TODO: [Refactor] Try to use DRY (Don't Repeat Yourself), there are some code blocks repeating
// DEPRECATED - This format isn't used by Federal Revenue anymore
// func (ci *CompanyImporter) CompaniesFromFixedWidthFile(fixedWidthFile string, bst *BinarySearchTree) error {
// 	var err error
// 	layoutFile, err := os.Open(fixedWidthFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer layoutFile.Close()

// 	companySchema, err := ci.loadLayoutSchema()
// 	if err == nil {
// 		var companyDoc map[string]interface{}
// 		var partnersDoc []model.Partner
// 		lineReader := bufio.NewScanner(layoutFile)
// 		for lineReader.Scan() {
// 			lineP := lineReader.Text()
// 			// Process line
// 			mapType := ci.findMapType(string(lineP[0]), companySchema)
// 			if mapType != nil {
// 				doc := ci.processLine(lineP, mapType)
// 				if companyDoc != nil && doc["_id"] != companyDoc["_id"] {
// 					// Save document
// 					companyDoc["socios"] = partnersDoc
// 					var co model.Company
// 					model.DecodeFromMap(companyDoc, &co)
// 					ci.saveCompany(co)
// 					companyDoc = nil
// 					partnersDoc = nil
// 				}
// 				switch lineP[0] {
// 				case '1':
// 					if companyDoc == nil {
// 						companyDoc = doc
// 						companyDoc["motivo_situacao_cadastral"] = ""
// 						companyDoc["grau_risco"] = ""
// 						// TODO: Using assertion to cast interface{} to int/string. There is a better way?
// 						status, err := ci.md.FindOneStatusDescriptionById(companyDoc["codigo_situacao_cadastral"].(int64))
// 						if err == nil {
// 							companyDoc["motivo_situacao_cadastral"] = status.Motivo
// 						}
// 						if rl, ok := companyDoc["cnae_fiscal"]; ok && rl != "" {
// 							riskLevel, err := ci.md.FindOneRiskLevelById(rl.(string)[:5])
// 							if err == nil {
// 								companyDoc["grau_risco"] = riskLevel.GrauRisco
// 							}
// 						}
// 					}

// 				case '2':
// 					partner := model.Partner{}
// 					model.DecodeFromMap(doc, &partner)
// 					if companyDoc != nil {
// 						partnersDoc = append(partnersDoc, partner)
// 					} else {
// 						// Probably we have an orphan record, save for later
// 						node := bst.Find(doc["_id"].(string))
// 						if node == nil {
// 							data := map[string]interface{}{}
// 							data["partner"] = []model.Partner{partner}
// 							bst.Add(doc["_id"].(string), data)
// 						} else {
// 							if partnerData, ok := node.Data["partner"]; ok {
// 								node.Data["partner"] = append(partnerData.([]model.Partner), partner)
// 							} else {
// 								node.Data["partner"] = []model.Partner{partner}
// 							}
// 						}
// 					}

// 				case '6':
// 					// Last line. Every 7 digits, is a CNAE code
// 					cnaes_sec := ci.getCNAEs(doc["cnaes_secundarios"].(string))
// 					if len(cnaes_sec) > 0 {
// 						if companyDoc != nil {
// 							companyDoc["cnaes_secundarios"] = cnaes_sec
// 						} else {
// 							// Probably we have an orphan record, save for later
// 							node := bst.Find(doc["_id"].(string))
// 							if node == nil {
// 								data := map[string]interface{}{}
// 								data["cnaes"] = cnaes_sec
// 								bst.Add(doc["_id"].(string), data)
// 							} else {
// 								if cnaesSlice, ok := node.Data["cnaes"]; ok {
// 									node.Data["cnaes"] = append(cnaesSlice.([]string), cnaes_sec...)
// 								} else {
// 									node.Data["cnaes"] = cnaes_sec
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 		if err = lineReader.Err(); err != nil {
// 			log.Println(err)
// 		} else {
// 			// Last record
// 			if companyDoc != nil {
// 				companyDoc["socios"] = partnersDoc
// 				var co model.Company
// 				model.DecodeFromMap(companyDoc, &co)
// 				ci.saveCompany(co)
// 			}
// 		}
// 	}
// 	return err
// }

// func (ci *CompanyImporter) CheckOrphans(bst *BinarySearchTree) {
// 	if bst != nil {
// 		bst.Walk(func(node *BTreeNode) {
// 			co, err := ci.md.FindOneCompanyById(node.ID)
// 			if err == nil {
// 				if cnaes, ok := node.Data["cnaes"]; ok {
// 					co.CNAEsSecundarios = cnaes.([]string)
// 				}
// 				if partners, ok := node.Data["partner"]; ok {
// 					co.Socios = append(co.Socios, partners.([]model.Partner)...)
// 				}
// 				ci.saveCompany(co)
// 			}
// 		})
// 	}
// }
