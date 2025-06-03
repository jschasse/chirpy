package auth

import (
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
    password := "testpassword123"
    
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword failed: %v", err)
    }
    
    if hash == "" {
        t.Error("Expected non-empty hash")
    }
    
    if hash == password {
        t.Error("Hash should not equal original password")
    }
}

func TestCheckPasswordHash(t *testing.T) {
    password := "testpassword123"
    wrongPassword := "wrongpassword"
    
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword failed: %v", err)
    }
    
    // Test correct password
    err = CheckPasswordHash(hash, password)
    if err != nil {
        t.Errorf("CheckPasswordHash should succeed with correct password: %v", err)
    }
    
    // Test wrong password
    err = CheckPasswordHash(hash, wrongPassword)
    if err == nil {
        t.Error("CheckPasswordHash should fail with wrong password")
    }
}

func TestMakeJWT(t *testing.T) {
    userID := uuid.New()
    secret := "test-secret-key"
    expiresIn := time.Hour
    
    tokenString, err := MakeJWT(userID, secret, expiresIn)
    if err != nil {
        t.Fatalf("MakeJWT failed: %v", err)
    }
    
    if tokenString == "" {
        t.Error("Expected non-empty token string")
    }
}

func TestValidateJWT_ValidToken(t *testing.T) {
    userID := uuid.New()
    secret := "test-secret-key"
    expiresIn := time.Hour
    
    // Create a token
    tokenString, err := MakeJWT(userID, secret, expiresIn)
    if err != nil {
        t.Fatalf("MakeJWT failed: %v", err)
    }
    
    // Validate the token
    extractedUserID, err := ValidateJWT(tokenString, secret)
    if err != nil {
        t.Fatalf("ValidateJWT failed: %v", err)
    }
    
    if extractedUserID != userID {
        t.Errorf("Expected userID %v, got %v", userID, extractedUserID)
    }
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
    userID := uuid.New()
    secret := "test-secret-key"
    
    // Create an expired token (negative duration means it's already expired)
    expiredDuration := -time.Hour
    tokenString, err := MakeJWT(userID, secret, expiredDuration)
    if err != nil {
        t.Fatalf("MakeJWT failed: %v", err)
    }
    
    // Try to validate the expired token
    _, err = ValidateJWT(tokenString, secret)
    if err == nil {
        t.Error("ValidateJWT should fail with expired token")
    }
}

func TestValidateJWT_WrongSecret(t *testing.T) {
    userID := uuid.New()
    correctSecret := "correct-secret-key"
    wrongSecret := "wrong-secret-key"
    expiresIn := time.Hour
    
    // Create a token with the correct secret
    tokenString, err := MakeJWT(userID, correctSecret, expiresIn)
    if err != nil {
        t.Fatalf("MakeJWT failed: %v", err)
    }
    
    // Try to validate with the wrong secret
    _, err = ValidateJWT(tokenString, wrongSecret)
    if err == nil {
        t.Error("ValidateJWT should fail with wrong secret")
    }
}

func TestValidateJWT_InvalidToken(t *testing.T) {
    secret := "test-secret-key"
    invalidToken := "invalid.token.string"
    
    _, err := ValidateJWT(invalidToken, secret)
    if err == nil {
        t.Error("ValidateJWT should fail with invalid token string")
    }
}

func TestValidateJWT_MalformedToken(t *testing.T) {
    secret := "test-secret-key"
    malformedToken := "not-a-jwt-token"
    
    _, err := ValidateJWT(malformedToken, secret)
    if err == nil {
        t.Error("ValidateJWT should fail with malformed token")
    }
}

func TestValidateJWT_EmptyToken(t *testing.T) {
    secret := "test-secret-key"
    emptyToken := ""
    
    _, err := ValidateJWT(emptyToken, secret)
    if err == nil {
        t.Error("ValidateJWT should fail with empty token")
    }
}

func TestValidateJWT_TokenWithInvalidUserID(t *testing.T) {
    secret := "test-secret-key"
    
    // Create a token with invalid subject (not a valid UUID)
    claims := jwt.RegisteredClaims{
        Issuer:    "chirpy",
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
        Subject:   "not-a-valid-uuid",
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        t.Fatalf("Failed to create test token: %v", err)
    }
    
    _, err = ValidateJWT(tokenString, secret)
    if err == nil {
        t.Error("ValidateJWT should fail with invalid user ID in token")
    }
}

func TestMakeJWT_Integration(t *testing.T) {
    testCases := []struct {
        name       string
        userID     uuid.UUID
        secret     string
        expiresIn  time.Duration
        shouldFail bool
    }{
        {
            name:       "valid token",
            userID:     uuid.New(),
            secret:     "valid-secret",
            expiresIn:  time.Hour,
            shouldFail: false,
        },
        {
            name:       "short expiration",
            userID:     uuid.New(),
            secret:     "valid-secret",
            expiresIn:  time.Second * 30,
            shouldFail: false,
        },
        {
            name:       "long expiration",
            userID:     uuid.New(),
            secret:     "valid-secret",
            expiresIn:  time.Hour * 24 * 365, // 1 year
            shouldFail: false,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            tokenString, err := MakeJWT(tc.userID, tc.secret, tc.expiresIn)
            
            if tc.shouldFail {
                if err == nil {
                    t.Error("Expected MakeJWT to fail, but it succeeded")
                }
                return
            }
            
            if err != nil {
                t.Fatalf("MakeJWT failed: %v", err)
            }
            
            // Validate the created token
            extractedUserID, err := ValidateJWT(tokenString, tc.secret)
            if err != nil {
                t.Fatalf("ValidateJWT failed: %v", err)
            }
            
            if extractedUserID != tc.userID {
                t.Errorf("Expected userID %v, got %v", tc.userID, extractedUserID)
            }
        })
    }
}