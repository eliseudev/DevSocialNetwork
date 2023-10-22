package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON retorna uma resposta em JSON para requisição
func JSON(w http.ResponseWriter, statuscode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)

	if dados != nil {
		if err := json.NewEncoder(w).Encode(dados); err != nil {
			log.Fatal(err)
		}
	}
}

// Err retorna um erro em formato JSON
func Err(w http.ResponseWriter, statuscode int, err error) {
	JSON(w, statuscode, struct {
		Err string `json:"err"`
	}{
		Err: err.Error(),
	})
}
