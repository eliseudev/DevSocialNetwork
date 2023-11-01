package controllers

import (
	"api/src/autenticacao"
	"api/src/db"
	"api/src/models"
	"api/src/repository"
	"api/src/response"
	"encoding/json"
	"io"
	"net/http"
)

func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioId, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		response.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var publicacao models.Publicacao
	if err = json.Unmarshal(bodyRequest, &publicacao); err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	publicacao.AutorID = usuarioId

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repository.RepositoryPublicacoes(db)
	publicacao.ID, err = repository.Criar(publicacao)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, publicacao)

}

func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {

}

func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {

}

func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {

}

func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {

}
