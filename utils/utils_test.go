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

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type FileDataTest struct {
	zipName   string
	fileName  string
	sha256Sum string
}

var unzipFilesTest = []FileDataTest{
	{
		zipName:   "../test-data/gopher.zip",
		fileName:  "gopher.jpeg",
		sha256Sum: "f46d2b67ba712ac545712539fca4a769615c2df8b77087f439eadbc5c8ee27fd",
	},
	{
		zipName:   "../test-data/page.zip",
		fileName:  "page.html",
		sha256Sum: "c568cf185ff26419b21b3644b30d5dab5fd927fea0ca65be8e2cdc9db0c3cd5a",
	},
}

func checkSHA256Sum(src string, hash string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return err
	}
	strHash := hex.EncodeToString(h.Sum(nil))
	// fmt.Printf("Read %s against %s\n", strHash, hash)
	if strHash != hash {
		return errors.New("Checksum fail")
	}
	return nil
}

func unzipTest(fileDataTest FileDataTest, unzipDir string) error {
	fmt.Printf("-> Unzip [%s]\n", fileDataTest.zipName)
	zipName, err := filepath.Abs(fileDataTest.zipName)
	if err != nil {
		return err
	}
	err = Unzip(zipName, unzipDir)
	if err != nil {
		return err
	}
	fmt.Printf("--> Checking SHA256 sum [%s]\n", fileDataTest.fileName)
	fileName, err := filepath.Abs(filepath.Join(unzipDir, fileDataTest.fileName))
	if err != nil {
		return err
	}
	return checkSHA256Sum(fileName, fileDataTest.sha256Sum)
}

func TestLoadCompanyDownloadConfig(t *testing.T) {
	fmt.Println("LoadCompanyDownloadConfig tests...")
	want := CompanyDownloadConfig{
		CompaniesMainUrl: "https://receita.economia.gov.br/orientacao/tributaria/cadastros/cadastro-nacional-de-pessoas-juridicas-cnpj/dados-publicos-cnpj",
		CompaniesMirrorUrls: []string{
			"https://data.brasil.io/mirror/socios-brasil/",
		},
		NR04Url: "https://www.gov.br/trabalho/pt-br/inspecao/seguranca-e-saude-no-trabalho/normas-regulamentadoras/nr-04.pdf/@@download/file/NR-04.pdf",
	}
	configFileTest := "../test-data/companies-download.json"
	got, err := LoadCompanyDownloadConfig(configFileTest)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, want: %v", got, want)
	}
}

func TestUnzip(t *testing.T) {
	fmt.Println("Utils Unzip tests...")
	unzipDir := "../test-data/unzip_tests"
	for _, v := range unzipFilesTest {
		err := unzipTest(v, unzipDir)
		if err != nil {
			t.Error(err)
		}
	}
	os.RemoveAll(unzipDir)
}

func TestFileDownload(t *testing.T) {
	fmt.Println("Utils FileDownload test...")
	toDir := "../test-data/tmp"
	urlTest := "https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/320px-Go_Logo_Blue.svg.png"
	toFile := filepath.Join(toDir, "320px-Go_Logo_Blue.svg.png")
	fmt.Printf("Attempt to download [%s] to [%s]\n", urlTest, toFile)
	err := FileDownload(urlTest, toDir)
	if err != nil {
		t.Error(err)
	}
	if _, err = os.Stat(toFile); !os.IsNotExist(err) {
		fmt.Println("File downloaded. Now removing...")
		os.Remove(toFile)
	} else {
		t.Error(toFile, "not downloaded!!")
	}
}
