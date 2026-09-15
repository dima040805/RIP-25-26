package apitypes

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/dima040805/RIP-25-26/internal/app/ds"
)

func TestUserToJSONHidesPassword(t *testing.T) {
	user := ds.User{ID: uuid.New(), Login: "alice", Password: "$2a$10$hash", IsModerator: true}

	body, err := json.Marshal(UserToJSON(user))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "password") || strings.Contains(string(body), "hash") {
		t.Fatalf("password hash leaked into the response: %s", body)
	}
	if !strings.Contains(string(body), `"login":"alice"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}
