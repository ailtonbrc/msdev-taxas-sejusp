package handlers

// ChaveUsuarioCpfCnpj é a chave usada no contexto do Gin para armazenar o usuário logado.
// Centralizada aqui para ser acessível por middleware.go, handlerTaxas.go, etc.
const ChaveUsuarioCpfCnpj = "usuarioCpfCnpj"
const ChaveUsuarioInfo = "usuarioInfo"

