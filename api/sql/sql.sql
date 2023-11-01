CREATE DATABASE devsocial if not exists;
USE devsocial;

DROP TABLE IF EXISTS publicacoes;
DROP TABLE IF EXISTS seguidores;
DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios(
    id int auto_increment primary key,
    nome varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    senha varchar(100) not null,
    created_at timestamp default current_timestamp()
)ENGINE=INNODB;

CREATE TABLE seguidores(
    usuario_id int not null,
    FOREIGN KEY (usuario_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,
    seguidor_id int not null,
    FOREIGN KEY(seguidor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,
    PRIMARY KEY(usuario_id, seguidor_id)
)ENGINE=INNODB;


CREATE TABLE publicacoes(
    id int auto_increment primary key,
    titulo varchar(50) not null,
    conteudo varchar(300),
    autor_id int not null,
    FOREIGN key(autor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,
    curtidas int default 0,
    created_at timestamp default current_timestamp()
)ENGINE=INNODB;