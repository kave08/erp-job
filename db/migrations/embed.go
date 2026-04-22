package migrations

import (
	"embed"
	"io/fs"
	"sort"
)

//go:embed *.sql
var files embed.FS

func FS() fs.FS {
	return files
}

func UpFiles() ([]string, error) {
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) >= len(".up.sql") && name[len(name)-len(".up.sql"):] == ".up.sql" {
			names = append(names, name)
		}
	}

	sort.Strings(names)
	return names, nil
}
