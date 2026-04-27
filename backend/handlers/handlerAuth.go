package handlers

import (
	"taxas-sejusp/backend/config"
	"taxas-sejusp/backend/modelos"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// IniciarLogin redireciona o usuário para a tela de login do provedor OIDC via Authorization Code Flow.
func IniciarLogin(c *gin.Context) {
	cfg := config.AppConfig

	// --- SUPORTE AO MODO MOCK ---
	if cfg.ModoAutenticacao == "mock" {
		log.Println("[AUTH MOCK] Simulando sucesso de login.")
		c.SetCookie("auth_token", "mock-token-secret-123", 3600, "/", "", false, true)
		c.SetCookie("user_cpf", cfg.MockCpfCnpj, 3600, "/", "", false, false)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Constrói a URL de autorização baseada no provedor configurado
	// Ex: https://des.id.ms.gov.br/auth/realms/ms/protocol/openid-connect/auth
	authEndpoint := fmt.Sprintf("%s/protocol/openid-connect/auth", cfg.OidcProviderURL)

	authURL, err := url.Parse(authEndpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao analisar URL de auth: " + err.Error()})
		return
	}

	query := authURL.Query()
	query.Set("client_id", cfg.EfazendaClientID)
	query.Set("redirect_uri", cfg.EfazendaRedirectURI)
	query.Set("response_type", "code")
	query.Set("scope", "openid profile email") // Solicitamos escopos padrão

	// TODO: Implementar state e nonce para segurança (CSRF)
	query.Set("state", "random_state_string")

	authURL.RawQuery = query.Encode()

	c.Redirect(http.StatusFound, authURL.String())
}

// CallbackLogin processa o retorno do provedor OIDC, troca o code por token e cria a sessão.
func CallbackLogin(c *gin.Context) {
	cfg := config.AppConfig
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Código de autorização não fornecido"})
		return
	}

	// Validação básica de state (deveria comparar com o armazenado na sessão antes do redirect)
	if state != "random_state_string" {
		log.Println("Aviso: State mismatch (ignorado por enquanto)")
	}

	// 1. Trocar Code por Access Token
	tokenEndpoint := fmt.Sprintf("%s/protocol/openid-connect/token", cfg.OidcProviderURL)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", cfg.EfazendaClientID)
	data.Set("client_secret", cfg.EfazendaClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.EfazendaRedirectURI)

	req, err := http.NewRequestWithContext(context.Background(), "POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao criar requisição de token"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "Falha ao comunicar com provedor de token"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Erro na troca de token. Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "Falha na troca de token", "detalhes": string(bodyBytes)})
		return
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao decodificar resposta do token"})
		return
	}

	// 2. Obter UserInfo para pegar o CPF/CNPJ e Nome
	userInfoEndpoint := fmt.Sprintf("%s/protocol/openid-connect/userinfo", cfg.OidcProviderURL)
	reqUserInfo, _ := http.NewRequestWithContext(context.Background(), "GET", userInfoEndpoint, nil)
	reqUserInfo.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	respUserInfo, err := client.Do(reqUserInfo)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "Falha ao obter dados do usuário"})
		return
	}
	defer respUserInfo.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(respUserInfo.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao decodificar dados do usuário"})
		return
	}

	// 3. Extrair dados
	// Tenta encontrar o CPF/CNPJ em campos comuns ou claims customizadas
	nome := "Usuário"
	if n, ok := userInfo["name"].(string); ok {
		nome = n
	} else if n, ok := userInfo["preferred_username"].(string); ok {
		nome = n
	}

	// Estratégia de busca do CPF: preferred_username, cpf, sub, etc.
	// Assumindo que o preferred_username pode ser o CPF ou existe um claim específico.
	// Vamos logar o userInfo para debug se necessário, mas por ora vamos tentar 'preferred_username' e 'cpf'
	cpfCnpj := ""
	if v, ok := userInfo["cpf"].(string); ok {
		cpfCnpj = v
	} else if v, ok := userInfo["preferred_username"].(string); ok {
		cpfCnpj = v // Keycloak muitas vezes usa o username como identificador principal
	} else if v, ok := userInfo["sub"].(string); ok {
		cpfCnpj = v // Fallback para o ID interno
	}

	// 4. Cria Sessão (Cookies)
	c.SetCookie("auth_token", tokenResp.AccessToken, 3600*8, "/", "", false, true)
	c.SetCookie("user_cpf", cpfCnpj, 3600*8, "/", "", false, false)
	c.SetCookie("user_name", nome, 3600*8, "/", "", false, false) // ✅ Novo cookie com o nome real do e-Fazenda

	// Log de sucesso
	log.Printf("Login e-Fazenda realizado: %s (%s)", nome, cpfCnpj)
	c.Redirect(http.StatusFound, "/")
}

// SairLogin limpa os cookies de sessão.
func SairLogin(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.SetCookie("user_cpf", "", -1, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"mensagem": "Logoff realizado com sucesso"})
}

// VerificarUsuario retorna os dados do usuário logado.
func VerificarUsuario(c *gin.Context) {
	// Obtém as informações do usuário injetadas pelo middleware
	userObj, exists := c.Get(ChaveUsuarioInfo)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"autenticado": false})
		return
	}

	user := userObj.(modelos.Usuario)

	c.JSON(http.StatusOK, modelos.RespostaUsuario{
		ID:           user.ID,
		Nome:         user.Nome,
		CPF:          user.CPF,
		Role:         user.Role,
		Instituicoes: user.Instituicoes,
		Autenticado:  true,
	})
}

