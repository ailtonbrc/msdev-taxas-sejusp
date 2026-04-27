package servicos

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"taxas-sejusp/backend/modelos"

	_ "modernc.org/sqlite"
)

type UsuarioService struct {
	db *sql.DB
}

var GlobalUsuarios *UsuarioService

func IniciarServicoUsuario() error {
	dbPath := "dados_cache/usuarios.db"
	
	// Garantir que a pasta existe
	if _, err := os.Stat("dados_cache"); os.IsNotExist(err) {
		os.Mkdir("dados_cache", 0777)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir banco de usuarios: %w", err)
	}

	// Criar tabela se não existir
	query := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		cpf TEXT UNIQUE NOT NULL,
		role TEXT NOT NULL,
		instituicoes TEXT -- Armazenado como string separada por vírgula
	);`
	
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("erro ao criar tabela de usuarios: %w", err)
	}

	GlobalUsuarios = &UsuarioService{db: db}

	// Garante usuários administradores iniciais
	adminsIniciais := []modelos.Usuario{
		{Nome: "Administrador Sistema", CPF: "admin", Role: "admin"},
		{Nome: "Ailton Branco", CPF: "62789678987", Role: "admin"},
	}

	for _, adm := range adminsIniciais {
		existente, _ := GlobalUsuarios.BuscarPorCPF(adm.CPF)
		if existente == nil {
			log.Printf("[USUARIO] Criando usuário administrador inicial: %s (%s)", adm.Nome, adm.CPF)
			GlobalUsuarios.Salvar(adm)
		}
	}

	return nil
}

func (s *UsuarioService) BuscarPorCPF(cpf string) (*modelos.Usuario, error) {
	var u modelos.Usuario
	var insts string
	err := s.db.QueryRow("SELECT id, nome, cpf, role, instituicoes FROM usuarios WHERE cpf = ?", cpf).
		Scan(&u.ID, &u.Nome, &u.CPF, &u.Role, &insts)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if insts != "" {
		u.Instituicoes = strings.Split(insts, ",")
	} else {
		u.Instituicoes = []string{}
	}

	return &u, nil
}

func (s *UsuarioService) ListarTodos() ([]modelos.Usuario, error) {
	rows, err := s.db.Query("SELECT id, nome, cpf, role, instituicoes FROM usuarios ORDER BY nome")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []modelos.Usuario
	for rows.Next() {
		var u modelos.Usuario
		var insts string
		if err := rows.Scan(&u.ID, &u.Nome, &u.CPF, &u.Role, &insts); err != nil {
			continue
		}
		if insts != "" {
			u.Instituicoes = strings.Split(insts, ",")
		} else {
			u.Instituicoes = []string{}
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *UsuarioService) Salvar(u modelos.Usuario) error {
	insts := strings.Join(u.Instituicoes, ",")
	
	if u.ID > 0 {
		_, err := s.db.Exec("UPDATE usuarios SET nome = ?, cpf = ?, role = ?, instituicoes = ? WHERE id = ?", 
			u.Nome, u.CPF, u.Role, insts, u.ID)
		return err
	}
	
	_, err := s.db.Exec("INSERT INTO usuarios (nome, cpf, role, instituicoes) VALUES (?, ?, ?, ?)", 
		u.Nome, u.CPF, u.Role, insts)
	return err
}

func (s *UsuarioService) Excluir(id int) error {
	_, err := s.db.Exec("DELETE FROM usuarios WHERE id = ?", id)
	return err
}
