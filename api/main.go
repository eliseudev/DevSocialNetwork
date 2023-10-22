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
	//Gerando Rotas API
	r := router.Gerar()
	fmt.Printf("Run API: localhost:%d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
