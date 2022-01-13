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
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/catfishlabs/goOpenCNPJ/consts"
	"github.com/catfishlabs/goOpenCNPJ/model"
)

const DBTESTURI = "mongodb://opencnpjtest:testSecret@localhost:27017/opendatatest"

func cities(t *testing.T) {
	fmt.Println("Running Cities import...")
	inputCSV, err := filepath.Abs("../test-data/F.K03200$Z.D10710.MUNIC.csv")
	if err != nil {
		t.Error(err)
	}
	md := model.NewMongoDatabase(DBTESTURI)
	err = md.Connect()
	if err != nil {
		t.Error(err)
	}
	defer md.Close()

	err = CitiesFromCSV(inputCSV, md)
	if err != nil {
		t.Error(err)
	}

	cityID := int64(8327)
	cityName := "SAO JOSE"
	result, err := md.FindOneCityById(cityID)
	if err != nil {
		t.Error(err)
	}
	if result.NomeMunicipio != cityName {
		t.Errorf("Expected %s, Got %s", cityName, result.NomeMunicipio)
	}
}

func statusDescription(t *testing.T) {
	fmt.Println("Running Status Description import...")
	inputCSV, err := filepath.Abs("../test-data/motivosituacao.csv")
	if err != nil {
		t.Error(err)
	}
	md := model.NewMongoDatabase(DBTESTURI)
	err = md.Connect()
	if err != nil {
		t.Error(err)
	}
	defer md.Close()

	err = StatusFromCSV(inputCSV, md)
	if err != nil {
		t.Error(err)
	}
	// Testing a query
	descriptionExpected := "INAPTIDAO (LEI 11.941/2009 ART.54)"
	result, err := md.FindOneStatusDescriptionById(71)
	if err != nil {
		t.Error(err)
	}
	if result.Motivo != descriptionExpected {
		t.Errorf("Expected: %s, receive: %s\n", descriptionExpected, result.Motivo)
	}
}

func nr04(t *testing.T) {
	fmt.Println("Running NR04 (Risk Level) import...")
	inputPDF, err := filepath.Abs("../test-data/NR-04.pdf")
	if err != nil {
		t.Error(err)
	}
	md := model.NewMongoDatabase(DBTESTURI)
	err = md.Connect()
	if err != nil {
		t.Error(err)
	}
	defer md.Close()

	err = NR04FromPDF(inputPDF, md)
	if err != nil {
		t.Error(err)
	}
	result, err := md.FindOneRiskLevelById("07103")
	if err != nil {
		t.Error(err)
	}
	if result.GrauRisco != "4" {
		t.Errorf("Expected: 4, Got: %s", result.GrauRisco)
	}
}

func companies(t *testing.T) {
	fmt.Println("Running Companies import...")
	inputFile01, err := filepath.Abs("../test-data/K03200Y0.ESTABELE.csv")
	if err != nil {
		t.Error(err)
	}
	inputFile02, err := filepath.Abs("../test-data/K03200Y0.EMPRECSV.csv")
	if err != nil {
		t.Error(err)
	}
	companyLayout, err := filepath.Abs("../config/cnpj-schema.json")
	if err != nil {
		t.Error(err)
	}
	md := model.NewMongoDatabase(DBTESTURI)
	err = md.Connect()
	if err != nil {
		t.Error(err)
	}
	defer md.Close()

	ci := NewCompanyImporter(companyLayout, md)
	gError := make(chan error)
	fInputs := []string{inputFile01, inputFile02}
	for _, f := range fInputs {
		go func(e chan<- error, fileTest string) {
			fmt.Println("-- Importing file:", fileTest)
			err := ci.CompaniesFromCSV(fileTest)
			e <- err
		}(gError, f)
	}
	// Wait
	for range fInputs {
		err := <-gError
		if err != nil {
			t.Error(err)
		}
	}
	baseID := "65747887"
	razaoSocial := "FULANO DA SILVA"
	fullID := "65747887000121"
	cnaeFiscal := "6120501"
	cnaeSecundario := "8599604"
	lenCnaeSecundario := 7
	dataInicioAtividade, _ := time.Parse(consts.DateLayoutSchema, "20171009")
	baseResult, err := md.FindOneBaseCompanyById(baseID)
	if err != nil {
		t.Error(err)
	}
	if baseResult.RazaoSocial != razaoSocial {
		t.Errorf("Expected: [%s], Got: [%s]", razaoSocial, baseResult.RazaoSocial)
	}
	result, err := md.FindOneCompanyById(fullID)
	if err != nil {
		t.Error(err)
	}
	dataInicioGot := time.Time(result.DataInicioAtividade)
	if dataInicioAtividade != dataInicioGot {
		t.Errorf("Expected: %v, Got: %v", dataInicioAtividade, dataInicioGot)
	}
	dataSituacaoGot := time.Time(result.DataSituacaoEspecial)
	if !dataSituacaoGot.IsZero() {
		t.Error("DataSituacaoEspecial is not Zero!!")
	}
	if result.CNAEFiscal != cnaeFiscal {
		t.Errorf("Expected: [%s], Got: [%s]", cnaeFiscal, result.CNAEFiscal)
	}
	if len(result.CNAEsSecundarios) == lenCnaeSecundario {
		if result.CNAEsSecundarios[4] != cnaeSecundario {
			t.Errorf("Expected: [%s], Got: [%s]", cnaeSecundario, result.CNAEsSecundarios[4])
		}
	} else {
		t.Errorf("Len CNAEsSecundarios is %d, Expected %d", len(result.CNAEsSecundarios), lenCnaeSecundario)
	}
}

func TestImporters(t *testing.T) {
	t.Run("Status", statusDescription)
	t.Run("NR04", nr04)
	t.Run("Cities", cities)
	t.Run("Companies", companies)
}
