CREATE DATABASE IF NOT EXISTS db_go;
USE db_go;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(50) not null,
    email varchar(50) not null unique,
    password varchar(100) not null
) ENGINE=INNODB;
