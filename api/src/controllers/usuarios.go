package controllers

import (
	"api/src/autenticacao"
	"api/src/db"
	"api/src/models"
	"api/src/repository"
	"api/src/response"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		response.Err(w, http.StatusUnprocessableEntity, err)
		return
	}
	var usuario models.Usuario
	if err = json.Unmarshal(bodyRequest, &usuario); err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = usuario.ValidFormat("cadastro"); err != nil {
		response.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repository.RepositoryUsuarios(db)
	usuario.ID, err = repository.Criar(usuario)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, usuario)
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.RepositoryUsuarios(db)
	usuarios, err := repo.Buscar(nomeNick)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, usuarios)
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	usuarioID, err := strconv.ParseUint(params["usuarioId"], 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.RepositoryUsuarios(db)
	usuario, err := repo.BuscarUsuarioID(usuarioID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, usuario)
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(params["usuarioId"], 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	usuarioIdToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	if usuarioID != usuarioIdToken {
		response.Err(w, http.StatusForbidden, errors.New("não é possivel atualizar esse usuario"))
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		response.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var usuario models.Usuario
	if err = json.Unmarshal(bodyRequest, &usuario); err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = usuario.ValidFormat("edicao"); err != nil {
		response.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repository.RepositoryUsuarios(db)
	if err = repository.Atualizar(usuarioID, usuario); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, err := strconv.ParseUint(params["usuarioId"], 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	usuarioIdToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	if usuarioID != usuarioIdToken {
		response.Err(w, http.StatusForbidden, errors.New("não é possivel deletar esse usuario"))
		return
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repository.RepositoryUsuarios(db)
	if err = repository.Deletar(usuarioID); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
