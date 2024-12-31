package database

import (
	"fmt"
	"os"
	"slices"
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
// It will also return a sorted list of the map's keys, this can be used to run the migrations in the correct order.
func filterMigrationMap(migrationMap map[int]string, currentVersion int) (map[int]string, []int) {
	scripts := map[int]string{}

	for k, v := range migrationMap {
		if k <= currentVersion {
			continue
		}

		scripts[k] = v
	}

	index := make([]int, 0, len(scripts))

	for k := range scripts {
		index = append(index, k)
	}

	slices.Sort(index)

	return scripts, index
}
