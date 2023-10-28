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

func (repo UsuarioDB) BuscarSeguidores(usuarioId uint64) ([]models.Usuario, error) {
	sqlQuery := `select u.id, u.nome, u.nick, u.email, u.created_at from usuarios
					u inner join seguidores s on u.id = s.seguidor_id  where s.usuario_id = ?`
	lines, err := repo.db.Query(sqlQuery, usuarioId)
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

func (repo UsuarioDB) BuscarSeguindo(usuarioId uint64) ([]models.Usuario, error) {
	sqlQuery := `select u.id, u.nome, u.nick, u.email, u.created_at from usuarios u
					inner join seguidores s on u.id = s.usuario_id  where s.seguidor_id = ?`
	lines, err := repo.db.Query(sqlQuery, usuarioId)
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

func (repo UsuarioDB) Seguir(usuarioId, seguidorId uint64) error {
	sqlQuery := "insert ignore into seguidores (usuario_id, seguidor_id) values(?, ?)"
	statement, err := repo.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuarioId, seguidorId); err != nil {
		return err
	}

	return nil
}

func (repo UsuarioDB) PararSeguir(usuarioId, seguidorId uint64) error {
	sqlQuery := "delete from seguidores where usuario_id = ? and seguidor_id = ?"
	statement, err := repo.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(usuarioId, seguidorId); err != nil {
		return err
	}

	return nil
}

func (repo UsuarioDB) BuscarSenha(usuarioId uint64) (string, error) {
	sqlQuery := "select senha from usuarios where id = ?"
	senha, err := repo.db.Query(sqlQuery, usuarioId)
	if err != nil {
		return "", err
	}
	defer senha.Close()

	var usuario models.Usuario

	if senha.Next() {
		if err = senha.Scan(&usuario.Senha); err != nil {
			return "", err
		}
	}

	return usuario.Senha, nil
}

func (repo UsuarioDB) AtualizarSenha(usuarioId uint64, senha string) error {
	sqlQuery := "update usuarios set senha = ? where id = ?"
	statement, err := repo.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(senha, usuarioId); err != nil {
		return err
	}

	return nil
}
