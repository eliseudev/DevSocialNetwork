package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//Carregando configurações para conexão banco de dados
	config.Carregar()
	fmt.Println(config.StringConnection)
	//Gerando Rotas API
	r := router.Gerar()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
