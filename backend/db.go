package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Caminho do arquivo do banco. Fica dentro da pasta backend/.
const dbPath = "crm.db"

// abrirBanco conecta ao SQLite, valida a conexão e cria as tabelas.
func abrirBanco() (*sql.DB, error) {
	// _pragma=foreign_keys(1) liga a checagem de chaves estrangeiras.
	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}

	// SQLite grava em um único arquivo; uma conexão evita erro de "database is locked".
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("conectar ao banco: %w", err)
	}

	if err := criarTabelas(db); err != nil {
		return nil, err
	}
	return db, nil
}

// criarTabelas cria clientes, oportunidades e vendas (se ainda não existirem).
func criarTabelas(db *sql.DB) error {
	const esquema = `
	CREATE TABLE IF NOT EXISTS clientes (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		nome           TEXT NOT NULL,
		empresa        TEXT NOT NULL,
		email          TEXT NOT NULL UNIQUE,
		telefone       TEXT,
		status         TEXT NOT NULL,
		responsavel    TEXT NOT NULL,
		ultimo_contato TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS oportunidades (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		cliente_id    INTEGER NOT NULL REFERENCES clientes(id),
		titulo        TEXT NOT NULL,
		etapa_funil   TEXT NOT NULL,
		valor         REAL NOT NULL,
		data_criacao  TEXT NOT NULL,
		UNIQUE (cliente_id, titulo)
	);

	CREATE TABLE IF NOT EXISTS vendas (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		cliente_id  INTEGER NOT NULL REFERENCES clientes(id),
		valor       REAL NOT NULL,
		data        TEXT NOT NULL,
		produto     TEXT NOT NULL,
		UNIQUE (cliente_id, data, produto)
	);`

	if _, err := db.Exec(esquema); err != nil {
		return fmt.Errorf("criar tabelas: %w", err)
	}
	return nil
}
