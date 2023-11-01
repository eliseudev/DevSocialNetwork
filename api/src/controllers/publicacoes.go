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

	"github.com/gorilla/mux"
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

	if err = publicacao.Preparar(); err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

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
	usuarioId, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.RepositoryPublicacoes(db)
	publicacoes, err := repo.Buscar(usuarioId)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, publicacoes)
}

func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoId, err := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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

	repo := repository.RepositoryPublicacoes(db)
	publicacao, err := repo.BuscarPublicacaoId(publicacaoId)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, publicacao)
}

func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioId, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	publicacaoId, err := strconv.ParseUint(params["publicacaoId"], 10, 64)
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

	repo := repository.RepositoryPublicacoes(db)
	publicacaoDb, err := repo.BuscarPublicacaoId(publicacaoId)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	if publicacaoDb.AutorID != usuarioId {
		response.Err(w, http.StatusForbidden, errors.New("não é possivel atualizar essa publicação"))
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

	if err = publicacao.Preparar(); err != nil {
		response.Err(w, http.StatusBadRequest, err)
		return
	}

	repo = repository.RepositoryPublicacoes(db)

	if err = repo.Atualizar(publicacaoId, publicacao); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioId, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	publicacaoId, err := strconv.ParseUint(params["publicacaoId"], 10, 64)
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

	repo := repository.RepositoryPublicacoes(db)
	publicacaoDb, err := repo.BuscarPublicacaoId(publicacaoId)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	if publicacaoDb.AutorID != usuarioId {
		response.Err(w, http.StatusForbidden, errors.New("não é possivel atualizar essa publicação"))
		return
	}

	if err = repo.Deletar(publicacaoId); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func BuscarPublicacoesUsuario(w http.ResponseWriter, r *http.Request) {
	parametos := mux.Vars(r)
	usuarioId, err := strconv.ParseUint(parametos["usuarioId"], 10, 64)
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

	repo := repository.RepositoryPublicacoes(db)
	publicacoes, err := repo.BuscarPublicacaoUsuario(usuarioId)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, publicacoes)
}

func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametos := mux.Vars(r)
	publicacaoId, err := strconv.ParseUint(parametos["publicacaoId"], 10, 64)
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

	repo := repository.RepositoryPublicacoes(db)
	if err = repo.Curtir(publicacaoId); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametos := mux.Vars(r)
	publicacaoId, err := strconv.ParseUint(parametos["publicacaoId"], 10, 64)
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

	repo := repository.RepositoryPublicacoes(db)
	if err = repo.Descurtir(publicacaoId); err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
