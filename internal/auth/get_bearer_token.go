package auth

import(
	"net/http"
	"fmt"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("Authorization header does not exist")
	}
	
	return strings.TrimSpace(strings.TrimPrefix(bearer, "Bearer")), nil
}