package repository

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type UsuarioDB struct {
	db *sql.DB
}

func RepositoryUsuarios(db *sql.DB) *UsuarioDB {
	return &UsuarioDB{db: db}
}

func (repo UsuarioDB) Criar(usuario models.Usuario) (uint64, error) {
	sqlQuery := "insert into usuarios(nome, nick, email, senha) values(? ,? ,? ,?)"
	statement, err := repo.db.Prepare(sqlQuery)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	resultado, err := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if err != nil {
		return 0, err
	}

	ultimoId, err := resultado.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ultimoId), nil
}

func (repo UsuarioDB) Buscar(nomeNick string) ([]models.Usuario, error) {
	nomeNick = fmt.Sprintf("%%%s%%", nomeNick)
	sqlQuery := "select id, nome, nick, email, created_at from usuarios where nome LIKE ? or nick LIKE ?"
	lines, err := repo.db.Query(sqlQuery, nomeNick, nomeNick)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var usuarios []models.Usuario
	for lines.Next() {
		var usuario models.Usuario
		if err = lines.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CreatedAt,
		); err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repo UsuarioDB) BuscarUsuarioID(usuarioId uint64) (models.Usuario, error) {
	sqlQuery := "select id, nome, nick, email, created_at from usuarios where id = ?"
	lines, err := repo.db.Query(sqlQuery, usuarioId)
	if err != nil {
		return models.Usuario{}, err
	}
	defer lines.Close()

	var usuario models.Usuario

	if lines.Next() {
		if err = lines.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CreatedAt,
		); err != nil {
			return models.Usuario{}, err
		}
	}

	return usuario, nil
}

func (repo UsuarioDB) Atualizar(usuarioId uint64, usuario models.Usuario) error {
	sqlQuery := "update usuarios set nome = ?, nick = ?, email = ? where id = ?"
	statement, err := repo.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuarioId); err != nil {
		return err
	}

	return nil
}

func (repo UsuarioDB) Deletar(usuarioId uint64) error {
	sqlQuery := "delete from usuarios where id = ?"
	statement, err := repo.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuarioId); err != nil {
		return err
	}

	return nil
}

func (repo UsuarioDB) BuscarUsuarioEmail(email string) (models.Usuario, error) {
	sqlQuery := "select id, senha from usuarios where email = ?"
	line, err := repo.db.Query(sqlQuery, email)
	if err != nil {
		return models.Usuario{}, err
	}
	defer line.Close()

	var usuario models.Usuario

	if line.Next() {
		if err = line.Scan(&usuario.ID, &usuario.Senha); err != nil {
			return models.Usuario{}, err
		}
	}
	return usuario, nil
}

func (respo UsuarioDB) Seguir(usuarioId, seguidorId uint64) error {
	sqlQuery := "insert ignore into seguidores (usuario_id, seguidor_id) values(?, ?)"
	statement, err := respo.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuarioId, seguidorId); err != nil {
		return err
	}

	return nil
}
