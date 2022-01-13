package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/catfishlabs/goOpenCNPJ/model"
	"github.com/catfishlabs/goOpenCNPJ/utils"
	"github.com/gorilla/mux"
)

type CompanyResponse struct {
	RazaoSocial string `json:"razao_social"`
	model.Company
}

func GetCompany(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"data":  nil,
		"error": "",
	}
	vars := mux.Vars(r)
	cnpj, keyExists := vars["cnpj"]
	if !keyExists {
		response["error"] = "invalid parameter"
		json.NewEncoder(w).Encode(response)
		return
	}
	cnpj = utils.RemoveChars(cnpj, ".-/")
	err := model.DB.Connect()
	if err != nil {
		response["error"] = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	defer model.DB.Close()

	company, err := model.DB.FindOneCompanyById(cnpj)
	if err != nil {
		response["error"] = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	companyResponse := CompanyResponse{"", company}
	baseCompany, err := model.DB.FindOneBaseCompanyById(cnpj[:8])
	if err == nil {
		// add base company data to response
		companyResponse.RazaoSocial = baseCompany.RazaoSocial
	}
	response["data"] = companyResponse
	json.NewEncoder(w).Encode(response)
}
