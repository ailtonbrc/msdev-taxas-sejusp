package servicos

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// ObterConexaoNoahContribuinte abre/cria o arquivo [CPF]_contrib.noah
func ObterConexaoNoahContribuinte(cpfCnpj string) (*sql.DB, error) {
	caminhoPasta := "dados_cache"
	if _, err := os.Stat(caminhoPasta); os.IsNotExist(err) {
		os.Mkdir(caminhoPasta, 0755)
	}

	caminhoArquivo := filepath.Join(caminhoPasta, fmt.Sprintf("%s_contrib.noah", cpfCnpj))
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL", caminhoArquivo)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro abrir Noah Contribuinte: %w", err)
	}

	// Tabela de Dados Contribuinte
	sqlData := `
	CREATE TABLE IF NOT EXISTS view_contribuinte (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cpf_login TEXT,
		ano_mes_int INTEGER,
		cod_mun INTEGER,
		municipio TEXT,
		cnpj_cpf_contrib TEXT,
		ano TEXT,
		mes TEXT,
		cnpj_administradora TEXT,
		razao_social_admin TEXT,
		vlr_mensal_bruto REAL,
		adicao REAL,
		cancel REAL,
		valor_liq REAL,
		credito REAL,
		debito REAL,
		boleto REAL,
		transf_rec REAL,
		dinheiro_outra_estrut REAL,
		pix REAL,
		voucher REAL,
		saque REAL,
		outros REAL,
		deposito REAL,
		bol_guias_rec REAL,
		data_realizacao_pesquisa TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_contrib_periodo ON view_contribuinte (ano_mes_int);
	CREATE INDEX IF NOT EXISTS idx_contrib_admin ON view_contribuinte (razao_social_admin);
	`
	if _, err := db.Exec(sqlData); err != nil {
		return nil, err
	}

	// Tabela de Controle
	sqlControl := `
	CREATE TABLE IF NOT EXISTS controle_carga (
		periodo INTEGER PRIMARY KEY,
		total_registros INTEGER,
		data_inicio TEXT,
		data_fim TEXT,
		duracao_segundos INTEGER,
		concluido INTEGER
	);
	`
	if _, err := db.Exec(sqlControl); err != nil {
		return nil, err
	}

	return db, nil
}

// ObterConexaoNoahAdministradora abre/cria o arquivo [CPF]_admin.noah
func ObterConexaoNoahAdministradora(cpfCnpj string) (*sql.DB, error) {
	caminhoPasta := "dados_cache"
	if _, err := os.Stat(caminhoPasta); os.IsNotExist(err) {
		os.Mkdir(caminhoPasta, 0755)
	}

	caminhoArquivo := filepath.Join(caminhoPasta, fmt.Sprintf("%s_admin.noah", cpfCnpj))
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL", caminhoArquivo)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro abrir Noah Administradora: %w", err)
	}

	// Tabela de Dados Administradora
	sqlData := `
	CREATE TABLE IF NOT EXISTS view_administradora (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cpf_login TEXT,
		ano_mes_int INTEGER,
		cod_mun INTEGER,
		municipio TEXT,
		cnpj_administradora TEXT,
		razao_social_admin TEXT,
		ano INTEGER,
		mes INTEGER,
		valor_liq REAL
	);
	CREATE INDEX IF NOT EXISTS idx_admin_periodo ON view_administradora (ano_mes_int);
	CREATE INDEX IF NOT EXISTS idx_admin_nome ON view_administradora (razao_social_admin);
	`
	if _, err := db.Exec(sqlData); err != nil {
		return nil, err
	}

	// Tabela de Controle (mesma estrutura)
	sqlControl := `
	CREATE TABLE IF NOT EXISTS controle_carga (
		periodo INTEGER PRIMARY KEY,
		total_registros INTEGER,
		data_inicio TEXT,
		data_fim TEXT,
		duracao_segundos INTEGER,
		concluido INTEGER
	);
	`
	if _, err := db.Exec(sqlControl); err != nil {
		return nil, err
	}

	return db, nil
}

