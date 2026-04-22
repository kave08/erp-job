package migrations

import (
	"io/fs"
	"strings"
	"testing"
)

func TestDeliveryStateRepairMigrationUsesCopySwapStrategy(t *testing.T) {
	t.Parallel()

	script, err := fs.ReadFile(FS(), "000003_delivery_state_entity_keys.up.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	body := string(script)
	requiredFragments := []string{
		"delivery_state_rebuild_000003",
		"delivery_state_legacy_000003",
		"INSERT INTO delivery_state_rebuild_000003",
		"RENAME TABLE",
		"GROUP BY operation_name",
	}
	for _, fragment := range requiredFragments {
		if !strings.Contains(body, fragment) {
			t.Fatalf("expected migration to contain %q", fragment)
		}
	}

	disallowedFragments := []string{
		"ALTER TABLE delivery_state\n\tDROP PRIMARY KEY",
		"ADD PRIMARY KEY (operation_name, entity_key);",
	}
	for _, fragment := range disallowedFragments {
		if strings.Contains(body, fragment) {
			t.Fatalf("expected migration to avoid in-place key mutation fragment %q", fragment)
		}
	}
}
