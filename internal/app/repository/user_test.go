package repository

import (
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("s3cret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hash == "s3cret" {
		t.Fatal("password is stored in plain text")
	}
	if !CheckPasswordHash("s3cret", hash) {
		t.Fatal("correct password was rejected")
	}
	if CheckPasswordHash("wrong", hash) {
		t.Fatal("wrong password was accepted")
	}
}

func TestGenerateToken(t *testing.T) {
	t.Setenv("JWT_KEY", "test-key")
	id := uuid.New()

	tokenString, err := GenerateToken(id, true)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return []byte("test-key"), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("token does not verify with the configured key: %v", err)
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims["user_id"] != id.String() {
		t.Fatalf("user_id = %v, want %v", claims["user_id"], id)
	}
	if claims["is_moderator"] != true {
		t.Fatalf("is_moderator = %v, want true", claims["is_moderator"])
	}
}

func TestBlacklistKeyForToken(t *testing.T) {
	key := blacklistKeyForToken("token-a")

	if key != blacklistKeyForToken("token-a") {
		t.Fatal("key is not deterministic")
	}
	if key == blacklistKeyForToken("token-b") {
		t.Fatal("different tokens share a key")
	}
	if !strings.HasPrefix(key, "blacklist:") {
		t.Fatalf("unexpected key format: %s", key)
	}
	if strings.Contains(key, "token-a") {
		t.Fatal("raw token leaks into the Redis key")
	}
}
