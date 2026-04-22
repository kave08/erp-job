package config

import "testing"

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
