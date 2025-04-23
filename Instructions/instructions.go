package instructions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RequestInput representa uma requisição
type RequestInput struct {
	URL     string
	Method  string
	Auth    string
	Body    []byte
	Headers map[string]string
}

// Fetcher realiza uma requisição HTTP baseada no RequestInput
func Fetcher(input RequestInput) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(input.Method, input.URL, bytes.NewBuffer(input.Body))
	if err != nil {
		return nil, err
	}

	if input.Auth != "" {
		req.Header.Add("Authorization", input.Auth)
	}
	for key, value := range input.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// SKUData é um tipo para agrupar dados dinâmicos
type SKUData map[string]interface{}

// ExtractFields percorre e extrai os campos desejados da resposta
func ExtractFields(responseBody []byte, fields []string) (map[string]SKUData, error) {
	var data interface{}
	err := json.Unmarshal(responseBody, &data)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string]SKUData)
	CollectFieldsV2(data, "", fields, grouped)
	return grouped, nil
}

// CollectFieldsV2 percorre os dados e agrupa apenas no nível correto
func CollectFieldsV2(data interface{}, currentPath string, fields []string, grouped map[string]SKUData) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newPath := key
			if currentPath != "" {
				newPath = currentPath + "." + key
			}

			// Se o campo é um dos que queremos (price, discount, etc.)
			if contains(fields, key) {
				groupPath := cleanPath(currentPath)
				if groupPath == "" {
					groupPath = "root"
				}
				if _, exists := grouped[groupPath]; !exists {
					grouped[groupPath] = SKUData{}
				}
				grouped[groupPath][key] = value
			}

			// Continua descendo (sem criar agrupamento vazio)
			CollectFieldsV2(value, newPath, fields, grouped)
		}
	case []interface{}:
		for i, item := range v {
			indexedPath := fmt.Sprintf("%s[%d]", currentPath, i)
			CollectFieldsV2(item, indexedPath, fields, grouped)
		}
	}
}

// cleanPath limpa o caminho removendo índices de arrays
func cleanPath(path string) string {
	cleaned := strings.ReplaceAll(path, "[0]", "")
	if cleaned == "" {
		return "root"
	}
	return cleaned
}

// contains verifica se um array contém uma string
func contains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}
