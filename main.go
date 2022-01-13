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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/catfishlabs/goOpenCNPJ/controllers"
	"github.com/catfishlabs/goOpenCNPJ/model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const Version = "0.0.1"

func initDatabaseInterface(dbURI string) {
	model.DB = model.NewMongoDatabase(dbURI)
}

func main() {
	var addr string
	// Load env config
	envConfig, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error reading configuration:", err)
	}
	initDatabaseInterface(envConfig["DBURI"])
	addr = fmt.Sprintf("%s:%s", envConfig["HOST"], envConfig["PORT"])
	router := mux.NewRouter()
	router.HandleFunc(
		"/about",
		func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{"version": Version})
		},
	).
		Methods("GET")

	router.HandleFunc(
		"/cnpj/{cnpj}",
		controllers.GetCompany,
	).
		Methods("GET")

	router.HandleFunc(
		"/nr04/{cnae}",
		controllers.GetNR04,
	).
		Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("goOpenCNPJ server v%s listening on %s\n", Version, addr)
	log.Fatal(srv.ListenAndServe())
}
