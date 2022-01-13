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
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type CompanyDownloadConfig struct {
	CompaniesMainUrl    string   `json:"companies.main.url"`
	CompaniesMirrorUrls []string `json:"companies.mirrors.url"`
	NR04Url             string   `json:"nr04.url"`
}

// LoadCompanyDownloadConfig loads the companies configuration JSON file
func LoadCompanyDownloadConfig(configFile string) (CompanyDownloadConfig, error) {
	result := CompanyDownloadConfig{}
	configFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		return result, err
	}
	fpConf, err := os.Open(configFileAbs)
	if err != nil {
		return result, err
	}
	defer fpConf.Close()

	err = json.NewDecoder(fpConf).Decode(&result)
	return result, err
}

// RemoveChars returns a copy of s removing each character in chars
func RemoveChars(s, chars string) string {
	result := s
	for _, i := range chars {
		result = strings.ReplaceAll(result, string(i), "")
	}
	return result
}

// Unzip uncompress a file.zip
// From: https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang#24792688
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// FileDownload get a file from url
func FileDownload(url string, toPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	downTo := filepath.Base(response.Request.URL.String())
	downTo = filepath.Join(toPath, downTo)
	fout, err := os.Create(downTo)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, response.Body)

	return err
}

// HeadersFrom get http headers from url
func HeadersFrom(url string) (http.Header, error) {
	h, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	return h.Header, err
}
