package controllers

import (
	"api/src/autenticacao"
	"api/src/db"
	"api/src/models"
	"api/src/repository"
	"api/src/response"
	"api/src/security"
	"encoding/json"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	db, err := db.Conectar()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.RepositoryUsuarios(db)
	usuarioBanco, err := repo.BuscarUsuarioEmail(usuario.Email)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerificarSenha(usuarioBanco.Senha, usuario.Senha); err != nil {
		response.Err(w, http.StatusUnauthorized, err)
		return
	}

	token, _ := autenticacao.CriarToken(usuarioBanco.ID)
	w.Write([]byte(token))
}
