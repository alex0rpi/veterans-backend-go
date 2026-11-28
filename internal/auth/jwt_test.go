package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateTokenExpiration(t *testing.T) {
	manager := NewJWTManager("test-secret")
	tests := []struct {
		name    string
		expires *jwt.NumericDate
		wantErr error
	}{
		{name: "valid", expires: jwt.NewNumericDate(time.Now().Add(time.Hour))},
		{name: "missing", wantErr: jwt.ErrTokenRequiredClaimMissing},
		{name: "expired", expires: jwt.NewNumericDate(time.Now().Add(-time.Hour)), wantErr: jwt.ErrTokenExpired},
	}

	for _, testToken := range tests {
		t.Run(testToken.name, func(t *testing.T) {
			claims := Claims{
				RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: testToken.expires},
				TokenType:        TokenTypeAccess,
			}
			token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(manager.secretKey))
			if err != nil {
				t.Fatal(err)
			}
			_, err = manager.ValidateToken(token, TokenTypeAccess)
			if !errors.Is(err, testToken.wantErr) {
				t.Fatalf("ValidateToken error = %v; want %v", err, testToken.wantErr)
			}
		})
	}
}
