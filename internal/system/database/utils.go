package database

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// genMigrationVersionThingName generates a name for the migration version tracker
// based on the collection name.
func genMigrationVersionThingName(collectionName string) string {
	return fmt.Sprintf("%s_migration_version", collectionName)
}

// createMigrationMap takes in a path to the folder with ti files and returns a map
// of it's version mapped to it's content
func createMigrationMap(path string) (map[int]string, error) {
	migrationMap := map[int]string{}

	fileEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range fileEntries {
		nameSlices := strings.Split(entry.Name(), "_")

		fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", path, entry.Name()))
		if err != nil {
			return nil, err
		}

		migrationVersion, err := strconv.Atoi(nameSlices[0])
		if err != nil {
			return nil, err
		}

		migrationMap[migrationVersion] = string(fileContent)
	}

	return migrationMap, nil
}

// filterMigrationMap takes in the migration map[int]string and a version. It then returns a list of scripts that need to be run.
func filterMigrationMap(migrationMap map[int]string, currentVersion int) map[int]string {
	scripts := map[int]string{}

	for k, v := range migrationMap {
		if k <= currentVersion {
			continue
		}

		scripts[k] = v
	}

	return scripts
}
