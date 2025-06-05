package auth

import(
	"net/http"
	"fmt"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("Authorization header does not exist")
	}
	
	return strings.TrimSpace(strings.TrimPrefix(apiKey, "ApiKey")), nil
}