package function

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"
)

func SnakeToCamel(snakeCase string) string {
	words := strings.Split(snakeCase, "_")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func CamelToSnake(input string) string {
	var result strings.Builder

	for i, char := range input {
		if unicode.IsUpper(char) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
