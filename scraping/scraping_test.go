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
	"fmt"
	"path/filepath"
	"testing"
)

func TestGetCNPJDataFilesLocal(t *testing.T) {
	localHTMLPage, error := filepath.Abs("../test-data/page-to-scrap.html")
	if error != nil {
		t.Error("Local HTML page file does not exist!!!")
	}
	fmt.Println("Testing SCRAPER against local file:")
	fmt.Println(localHTMLPage)
	shouldBe := struct {
		DataFiles  []string
		LastUpdate string
		StatusFile string
		CitiesFile string
	}{
		DataFiles: []string{
			"http://200.152.38.155/CNPJ/K3241.K03200Y0.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y1.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y2.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y3.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y4.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y5.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y6.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y7.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y8.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y9.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y0.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y1.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y2.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y3.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y4.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y5.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y6.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y7.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y8.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y9.D10710.ESTABELE.zip",
		},
		LastUpdate: "16/07/2021",
		StatusFile: "http://200.152.38.155/CNPJ/F.K03200$Z.D10710.MOTICSV.zip",
		CitiesFile: "http://200.152.38.155/CNPJ/F.K03200$Z.D10710.MUNICCSV.zip",
	}
	got := NewLocalFile(localHTMLPage)
	got.GetCNPJData()
	// fmt.Println(got)
	fmt.Println(got.DataFiles)
	fmt.Println("----------------------------------------")
	fmt.Println(got.LastUpdate)
	for _, f := range got.DataFiles {
		hasHref := false
		for _, sb := range shouldBe.DataFiles {
			if f == sb {
				hasHref = true
			}
		}
		if !hasHref {
			t.Error("Data files missing!!!")
		}
	}
	if got.LastUpdate != shouldBe.LastUpdate {
		t.Errorf("Last update date: should be %s, got %s", shouldBe.LastUpdate, got.LastUpdate)
	}
	if got.StatusFile != shouldBe.StatusFile {
		t.Errorf("Status file: should be %s, got %s", shouldBe.StatusFile, got.StatusFile)
	}
	if got.CitiesFile != shouldBe.CitiesFile {
		t.Errorf("Cities file: should %s, got %s", shouldBe.CitiesFile, got.CitiesFile)
	}
	fmt.Println("----------------------------------------")
	fmt.Println("")
}

func TestGetGetCNPJDataFilesURL(t *testing.T) {
	url := "https://www.gov.br/receitafederal/pt-br/assuntos/orientacao-tributaria/cadastros/consultas/dados-publicos-cnpj"
	fmt.Println("Testing SCRAPER against URL:")
	fmt.Println(url)
	got := New(url)
	got.GetCNPJData()
	// fmt.Println(got)
	fmt.Println(got.DataFiles)
	fmt.Println("----------------------------------------")
	fmt.Println(got.LastUpdate)
	fmt.Println("----------------------------------------")
	fmt.Println(got.StatusFile)
	fmt.Println("----------------------------------------")
	fmt.Println(got.CitiesFile)
	fmt.Println("----------------------------------------")
	fmt.Println("")
}
