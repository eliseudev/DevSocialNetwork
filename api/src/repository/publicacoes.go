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

func (repo Publicacoes) BuscarPublicacaoId(publicacaoID uint64) (models.Publicacao, error) {
	sql := "select p.*, u.nick from publicacoes p inner join usuarios u on u.id = p.autor_id where p.id = ?"
	line, err := repo.db.Query(sql, publicacaoID)
	if err != nil {
		return models.Publicacao{}, err
	}
	defer line.Close()

	var publicacao models.Publicacao
	if line.Next() {
		if err = line.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CreatedAt,
			&publicacao.AutorNick,
		); err != nil {
			return models.Publicacao{}, err
		}
	}

	return publicacao, nil
}

func (repo Publicacoes) Buscar(usuarioId uint64) ([]models.Publicacao, error) {
	sql := "select distinct p.*, u.nick from publicacoes p inner join usuarios u on u.id = p.autor_id inner join seguidores s on p.autor_id = s.usuario_id where u.id = ? or s.seguidor_id = ? order by 1 desc;"
	lines, err := repo.db.Query(sql, usuarioId, usuarioId)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var publicacoes []models.Publicacao
	if lines.Next() {
		var publicacao models.Publicacao
		if err = lines.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CreatedAt,
			&publicacao.AutorNick,
		); err != nil {
			return nil, err
		}
		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (repo Publicacoes) Atualizar(publicacaoId uint64, publicacao models.Publicacao) error {
	sql := "update publicacoes se titulo = ?, conteudo = ? where id = ?"
	statement, err := repo.db.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacaoId); err != nil {
		return err
	}

	return nil
}

func (repo Publicacoes) Deletar(publicacaoId uint64) error {
	sql := "delete from publicacoes where id = ?"
	statement, err := repo.db.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacaoId); err != nil {
		return err
	}

	return nil
}

func (repo Publicacoes) BuscarPublicacaoUsuario(usuarioId uint64) ([]models.Publicacao, error) {
	sql := "select p.*, u.nick from publicacoes p join usuarios u on u.id = p.autor_id where p.autor_id = ?"
	lines, err := repo.db.Query(sql, usuarioId)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var publicacoes []models.Publicacao
	if lines.Next() {
		var publicacao models.Publicacao
		if err = lines.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CreatedAt,
			&publicacao.AutorNick,
		); err != nil {
			return nil, err
		}
		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (repo Publicacoes) Curtir(publicacaoId uint64) error {
	sql := "update publicacoes set curtidas = curtidas + 1 where id = ?"
	statement, err := repo.db.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacaoId); err != nil {
		return err
	}

	return nil

}

func (repo Publicacoes) Descurtir(publicacaoId uint64) error {
	sql := "update publicacoes set curtidas = CASE WHEN curtidas > 0 THEN curtidas -1 ELSE 0 END where id = ?"
	statement, err := repo.db.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicacaoId); err != nil {
		return err
	}

	return nil

}
