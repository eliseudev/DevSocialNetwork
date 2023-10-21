CREATE DATABASE devsocial if not exists;
USE devsocial;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios(
    id int auto_increment primary key,
    nome varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    senha varchar(20) not null unique,
    created_at timestamp default current_timestamp()
)ENGINE=INNODB;