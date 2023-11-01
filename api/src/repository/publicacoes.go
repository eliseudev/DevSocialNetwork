package repository

import (
	"api/src/models"
	"database/sql"
)

type Publicacoes struct {
	db *sql.DB
}

func RepositoryPublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

func (repo Publicacoes) Criar(publicacao models.Publicacao) (uint64, error) {
	sql := "insert into publicacoes (titulo, conteudo, autor_id) values(?, ?, ?)"
	statement, err := repo.db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	resultado, err := statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacao.AutorID)
	if err != nil {
		return 0, err
	}

	ultimoId, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ultimoId), nil
}
