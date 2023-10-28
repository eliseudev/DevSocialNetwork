package controllers

import (
	"api/src/autenticacao"
	"api/src/db"
	"api/src/models"
	"api/src/repository"
	"api/src/response"
	"api/src/security"
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

func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorId, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	parametros := mux.Vars(r)
	usuarioId, err := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	if seguidorId == usuarioId {
		response.Err(w, http.StatusForbidden, errors.New("não é possivel seguir você mesmo"))
		return
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repository.RepositoryUsuarios(db)
	if err = repository.Seguir(usuarioId, seguidorId); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func ParaSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorId, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	usuarioId, err := strconv.ParseUint(params["usuarioId"], 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	if seguidorId == usuarioId {
		response.Err(w, http.StatusForbidden, errors.New("não é possivel deixar de seguir você mesmo"))
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()
	repo := repository.RepositoryUsuarios(db)
	if err = repo.PararSeguir(usuarioId, seguidorId); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parms := mux.Vars(r)
	usuarioId, err := strconv.ParseUint(parms["usuarioId"], 10, 64)
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
	seguidores, err := repo.BuscarSeguidores(usuarioId)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, seguidores)
}

func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parms := mux.Vars(r)
	usuarioId, err := strconv.ParseUint(parms["usuarioId"], 10, 64)
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
	seguindo, err := repo.BuscarSeguindo(usuarioId)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, seguindo)
}

func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIdToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	parms := mux.Vars(r)
	usuarioIdRequest, err := strconv.ParseUint(parms["usuarioId"], 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	if usuarioIdToken != usuarioIdRequest {
		response.Err(w, http.StatusForbidden, errors.New("você não está autorizado a atualizar esse usuario"))
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)

	var senha models.Senha
	if err = json.Unmarshal(bodyRequest, &senha); err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repos := repository.RepositoryUsuarios(db)
	senhaBanco, err := repos.BuscarSenha(usuarioIdToken)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = security.VerificarSenha(senhaBanco, senha.Atual); err != nil {
		response.Err(w, http.StatusUnauthorized, errors.New("a senha atual esta incorreta"))
		return
	}

	senhaHash, err := security.Hash(senha.Nova)
	if err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = repos.AtualizarSenha(usuarioIdToken, string(senhaHash)); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
