package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestExtractTokenFromHeader(t *testing.T) {
	cases := map[string]string{
		"":               "",
		"Bearer abc.def": "abc.def",
		"Basic abc":      "",
		"Bearer":         "",
		"bearer abc":     "",
	}
	for header, want := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if header != "" {
			req.Header.Set("Authorization", header)
		}
		if got := extractTokenFromHeader(req); got != want {
			t.Errorf("header %q: got %q, want %q", header, got, want)
		}
	}
}

func TestCORSMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	preflight := httptest.NewRecorder()
	router.ServeHTTP(preflight, httptest.NewRequest(http.MethodOptions, "/ping", nil))
	if preflight.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d, want %d", preflight.Code, http.StatusNoContent)
	}
	if got := preflight.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}

	get := httptest.NewRecorder()
	router.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/ping", nil))
	if get.Code != http.StatusOK || get.Body.String() != "pong" {
		t.Fatalf("GET passed through as %d %q", get.Code, get.Body.String())
	}
}

func TestModeratorMiddlewareRejectsRequests(t *testing.T) {
	t.Setenv("JWT_KEY", "server-key")

	foreign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      "00000000-0000-0000-0000-000000000001",
		"is_moderator": true,
		"exp":          time.Now().Add(time.Hour).Unix(),
	})
	foreignToken, err := foreign.SignedString([]byte("attacker-key"))
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"no token":          "",
		"foreign signature": "Bearer " + foreignToken,
		"malformed token":   "Bearer not-a-jwt",
	}
	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			h := &Handler{}
			router := gin.New()
			router.GET("/private", h.ModeratorMiddleware(false), func(c *gin.Context) { c.Status(http.StatusOK) })

			req := httptest.NewRequest(http.MethodGet, "/private", nil)
			if header != "" {
				req.Header.Set("Authorization", header)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestTokenTTLFromClaims(t *testing.T) {
	future := time.Now().Add(time.Hour).Unix()

	ttl, err := tokenTTLFromClaims(jwt.MapClaims{"exp": float64(future)})
	if err != nil || ttl <= 0 || ttl > time.Hour {
		t.Fatalf("float exp: ttl=%v err=%v", ttl, err)
	}
	if _, err := tokenTTLFromClaims(jwt.MapClaims{"exp": json.Number(strconv.FormatInt(future, 10))}); err != nil {
		t.Fatalf("json.Number exp: %v", err)
	}
	if _, err := tokenTTLFromClaims(jwt.MapClaims{}); err == nil {
		t.Fatal("missing exp must be an error")
	}
	if _, err := tokenTTLFromClaims(jwt.MapClaims{"exp": float64(time.Now().Add(-time.Minute).Unix())}); err == nil {
		t.Fatal("expired token must be an error")
	}
	if _, err := tokenTTLFromClaims(jwt.MapClaims{"exp": "tomorrow"}); err == nil {
		t.Fatal("string exp must be an error")
	}
}
