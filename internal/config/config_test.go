package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRejectsEnabledOTelWithoutEndpoint(t *testing.T) {
	t.Parallel()

	cfg := validConfig()
	cfg.OTel.Enabled = true
	cfg.OTel.Endpoint = ""

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected OTEL validation error")
	}
}

func TestValidateAcceptsValidConfig(t *testing.T) {
	t.Parallel()

	cfg := validConfig()
	cfg.OTel.Enabled = true

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}

func TestLoadDatabaseAcceptsDatabaseOnlyConfig(t *testing.T) {
	t.Parallel()

	configPath := writeConfigFile(t, `
DATABASE:
  USERNAME: user
  PASSWORD: pass
  HOST: localhost
  PORT: "3306"
  DB_NAME: erp_job
  MAX_OPEN_CONNECTIONS: 5
  MAX_IDLE_CONNECTIONS: 2
`)

	databaseConfig, err := LoadDatabase(configPath)
	if err != nil {
		t.Fatalf("LoadDatabase returned error: %v", err)
	}
	if databaseConfig.DBName != "erp_job" {
		t.Fatalf("unexpected database config: %#v", databaseConfig)
	}
}

func TestLoadRejectsDatabaseOnlyConfig(t *testing.T) {
	t.Parallel()

	configPath := writeConfigFile(t, `
DATABASE:
  USERNAME: user
  PASSWORD: pass
  HOST: localhost
  PORT: "3306"
  DB_NAME: erp_job
  MAX_OPEN_CONNECTIONS: 5
  MAX_IDLE_CONNECTIONS: 2
`)

	if _, err := Load(configPath); err == nil {
		t.Fatal("expected transfer config validation error")
	}
}

func writeConfigFile(t *testing.T, body string) string {
	t.Helper()

	configPath := filepath.Join(t.TempDir(), "env.yml")
	if err := os.WriteFile(configPath, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	return configPath
}

func validConfig() Config {
	return Config{
		AryanApp: AryanApp{
			BaseURL: "https://aryan.example.com",
			APIKey:  "aryan-key",
			Timeout: 5,
		},
		FararavandApp: FararavandApp{
			BaseURL: "https://fararavand.example.com",
			APIKey:  "fararavand-key",
			Timeout: 5,
		},
		Database: Database{
			Username:           "user",
			Password:           "pass",
			Host:               "localhost",
			Port:               "3306",
			DBName:             "erp_job",
			MaxOpenConnections: 5,
			MaxIdleConnections: 2,
		},
		OTel: OTel{
			Enabled:     false,
			Endpoint:    "http://otel-collector:4318",
			ServiceName: "erp-job",
			Environment: "test",
		},
	}
}
