package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/catfishlabs/goOpenCNPJ/model"
	"github.com/gorilla/mux"
)

func GetNR04(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"data":  nil,
		"error": "",
	}
	vars := mux.Vars(r)
	cnae, keyExists := vars["cnae"]
	if !keyExists {
		response["error"] = "invalid parameter"
		json.NewEncoder(w).Encode(response)
		return
	}

	err := model.DB.Connect()
	if err != nil {
		response["error"] = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	defer model.DB.Close()

	riskLevel, err := model.DB.FindOneRiskLevelById(cnae[:5])
	if err != nil {
		response["error"] = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response["data"] = riskLevel
	json.NewEncoder(w).Encode(response)
}
