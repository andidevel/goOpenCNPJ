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
	"bytes"
	"os"
	"path/filepath"

	"github.com/ledongthuc/pdf"
)

func openPDF(pdfFile string) (*os.File, *pdf.Reader, error) {
	pdfFileAbs, err := filepath.Abs(pdfFile)
	if err != nil {
		return nil, nil, err
	}
	f, reader, err := pdf.Open(pdfFileAbs)
	if err != nil {
		return nil, nil, err
	}
	return f, reader, nil
}

func PDFPlainText(pdfFile string) (string, error) {
	f, r, err := openPDF(pdfFile)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()
	var buffer bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buffer.ReadFrom(b)
	return buffer.String(), nil
}

func PDFText(pdfFile string) ([]pdf.Rows, error) {
	f, r, err := openPDF(pdfFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()
	pages := r.NumPage()
	var result []pdf.Rows
	for i := 1; i <= pages; i++ {
		p := r.Page(i)
		if !p.V.IsNull() {
			row, _ := p.GetTextByRow()
			result = append(result, row)
		}
	}
	return result, nil
}
