package docker_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readDockerFile(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(".", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func TestComposeDefinesPostgresNetworkAndVolume(t *testing.T) {
	t.Parallel()
	content := readDockerFile(t, "docker-compose.yml")

	required := []string{
		"name: client-manager",
		"postgres:",
		"image: postgres:16-alpine",
		"POSTGRES_PASSWORD:",
		"healthcheck:",
		"pg_isready",
		"networks:",
		"clientmanager:",
		"postgres_data:",
		"no-new-privileges:true",
	}
	for _, want := range required {
		if !strings.Contains(content, want) {
			t.Errorf("docker-compose.yml missing %q", want)
		}
	}

	forbidden := []string{
		"init-schemas.sql",
		"docker-entrypoint-initdb.d",
		"CREATE SCHEMA",
		"POSTGRES_PASSWORD: postgres",
		"POSTGRES_PASSWORD: password",
	}
	for _, bad := range forbidden {
		if strings.Contains(content, bad) {
			t.Errorf("docker-compose.yml must not contain %q", bad)
		}
	}
}

func TestComposeDoesNotShipBootstrapSchemaSQL(t *testing.T) {
	t.Parallel()
	path := filepath.Join("postgres", "init-schemas.sql")
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("%s must not exist: schemas come from per-service migrations", path)
	}
}

func TestServiceEntrypointAutomatesMigrations(t *testing.T) {
	t.Parallel()
	content := readDockerFile(t, filepath.Join("bin", "service-entrypoint.sh"))

	required := []string{
		"MIGRATE_CMD",
		"DATABASE_URL",
		"SKIP_MIGRATE",
		"wait_for_tcp",
		`exec "$@"`,
	}
	for _, want := range required {
		if !strings.Contains(content, want) {
			t.Errorf("service-entrypoint.sh missing %q", want)
		}
	}

	info, err := os.Stat(filepath.Join("bin", "service-entrypoint.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0o111 == 0 {
		t.Error("service-entrypoint.sh must be executable")
	}
}

func TestComposeDefinesConfigurationService(t *testing.T) {
	t.Parallel()
	content := readDockerFile(t, "docker-compose.yml")

	required := []string{
		"configuration:",
		"services/configuration/Dockerfile",
		"DATABASE_URL:",
		"HTTP_ADDR:",
		"condition: service_healthy",
		"search_path=configuration",
	}
	for _, want := range required {
		if !strings.Contains(content, want) {
			t.Errorf("docker-compose.yml missing %q", want)
		}
	}
}

func TestConfigurationDockerfileExists(t *testing.T) {
	t.Parallel()
	path := filepath.Join("..", "services", "configuration", "Dockerfile")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("configuration Dockerfile: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, want := range []string{
		"FROM golang:1.25-alpine AS build",
		"service-entrypoint.sh",
		"USER appuser",
		"/app/configuration",
		"/app/migrate",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("Dockerfile missing %q", want)
		}
	}
}

func TestEnvExampleHasNoRealSecretAndDocumentsCopy(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", ".env.example")
	data, err := os.ReadFile(root)
	if err != nil {
		t.Fatalf("read .env.example: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "POSTGRES_PASSWORD=") {
		t.Fatal(".env.example must define POSTGRES_PASSWORD")
	}
	if !strings.Contains(content, "change-me") {
		t.Error(".env.example should use an obvious placeholder password")
	}
	if !strings.Contains(content, "cp .env.example .env") {
		t.Error(".env.example should document copying to .env")
	}
}
