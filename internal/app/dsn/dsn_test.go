package dsn

import "testing"

func TestFromEnv(t *testing.T) {
	t.Setenv("DB_HOST", "")
	if got := FromEnv(); got != "" {
		t.Fatalf("without DB_HOST got %q, want empty string", got)
	}

	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "app")
	t.Setenv("DB_PASS", "pass")
	t.Setenv("DB_NAME", "exoplanets")

	want := "host=db port=5432 user=app password=pass dbname=exoplanets sslmode=disable"
	if got := FromEnv(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
