package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type Usuario struct {
	ID        uint64    `json:"id,omitempty"`
	Nome      string    `json:"nome,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Senha     string    `json:"senha,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (usuario *Usuario) ValidFormat(etapa string) error {
	if err := usuario.isValid(etapa); err != nil {
		return err
	}

	if err := usuario.isFormat(etapa); err != nil {
		return err
	}

	return nil
}

func (usuario *Usuario) isValid(etapa string) error {
	if usuario.Nome == "" {
		return errors.New("o nome e obrigatório")
	}

	if usuario.Nick == "" {
		return errors.New("o nick é obrigatório")
	}

	if usuario.Email == "" {

		return errors.New("o e-mail é obrigatório")
	}

	if err := checkmail.ValidateFormat(usuario.Email); err != nil {
		return errors.New("por favor digite um e-mail válido")
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("a senha é obrigatório")
	}

	return nil
}

func (usuario *Usuario) isFormat(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if etapa == "cadastro" {
		senhaHash, err := security.Hash(usuario.Senha)
		if err != nil {
			return err
		}
		usuario.Senha = string(senhaHash)
	}

	return nil
}
