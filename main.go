package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	instructions "github.com/calvinalex/smartlet/Instructions"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega o .env se existir
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, prosseguindo...")
	}

	// Pega as variáveis de ambiente
	url := os.Getenv("API_URL")
	method := os.Getenv("API_METHOD")
	auth := os.Getenv("API_AUTH") // Tipo "Bearer xyz" ou vazio

	// Valida se as variáveis de ambiente necessárias estão presentes
	if url == "" || method == "" {
		log.Fatal("API_URL e API_METHOD são obrigatórios")
	}

	// Configura os headers para a request
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Monta a estrutura da request
	request := instructions.RequestInput{
		URL:     url,
		Method:  method,
		Auth:    auth,
		Headers: headers,
		Body:    nil, // poderia montar um JSON []byte se fosse POST
	}

	// Faz a requisição
	response, err := instructions.Fetcher(request)
	if err != nil {
		log.Fatalf("Erro ao fazer a requisição: %v", err)
	}

	// Define os campos que queremos extrair
	fields := []string{"orderId", "price", "discount", "sku"}

	// Extrai e agrupa os campos
	grouped, err := instructions.ExtractFields(response, fields)
	if err != nil {
		log.Fatalf("Erro ao extrair campos: %v", err)
	}

	// Formata e imprime o resultado
	output, _ := json.MarshalIndent(grouped, "", "  ")
	fmt.Println(string(output))
}
