package models

import (
	"errors"
	"strings"
	"time"
)

type Usuario struct {
	ID        uint64    `json:"id,omitempty"`
	Nome      string    `json:"nome,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Senha     string    `json:"senha,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (usuario *Usuario) ValidFormat(stage string) error {
	if err := usuario.isValid(stage); err != nil {
		return err
	}

	usuario.isFormat()
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
		return errors.New("o email é obrigatório")
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("a senha é obrigatório")
	}

	return nil
}

func (usuario *Usuario) isFormat() {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)
}
