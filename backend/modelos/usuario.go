package modelos

type Usuario struct {
	ID           int      `json:"id"`
	Nome         string   `json:"nome"`
	CPF          string   `json:"cpf"`
	Role         string   `json:"role"`         // 'admin' ou 'user'
	Instituicoes []string `json:"instituicoes"` // Vertentes permitidas
}

// RespostaUsuario para o frontend
type RespostaUsuario struct {
	ID           int      `json:"id"`
	Nome         string   `json:"nome"`
	CPF          string   `json:"cpf"`
	Role         string   `json:"role"`
	Instituicoes []string `json:"instituicoes"`
	Autenticado  bool     `json:"autenticado"`
}
