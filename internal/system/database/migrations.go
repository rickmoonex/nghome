package database

import (
	"errors"
	"fmt"
)

// TODO: Logging in this here is quite shit, needs a bunch of improving...

// requiredCollections represents the collections that need to be presents on the DB
// These names will also be used to identify files in the migrations folder.
// The `db_man` collection will also hold migration version information, it's created seperatly.
// These version 'Things' will be named using the collection name.
var requiredCollections = []string{
	"system_space",
	"user_space",
}

// AutoMigrate will automatically handle the migration of the database.
// It takes in a path to the migrations folder holding the .ti scripts.
func (c *Client) AutoMigrate(path string) error {
	// Ensure the `db_man` collection exists
	err := c.ensureCollection("db_man")
	if err != nil {
		return err
	}

	// We make sure that the requiredCollections exists
	for _, v := range requiredCollections {
		err := c.ensureCollection(v)
		if err != nil {
			return err
		}
	}

	// Make sure the migration version counters exist
	err = c.createMigrationVersionThings(requiredCollections)
	if err != nil {
		return err
	}

	// Migrate each collection seperatly
	for _, collection := range requiredCollections {
		collectionPath := fmt.Sprintf("%s/%s", path, collection)

		migMap, err := createMigrationMap(collectionPath)
		if err != nil {
			return err
		}

		version, err := c.getMigrationVersion(collection)
		if err != nil {
			return err
		}

		scripts, index := filterMigrationMap(migMap, version)

		for i := range index {
			err := c.runMigrationScript(collection, scripts[i])
			if err != nil {
				return err
			}
			err = c.updateMigrationVersion(collection, i)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// runMigrationScript takes a migration script content and a collection name,
// it then runs this script on the ThingsDB server.
func (c *Client) runMigrationScript(collectionName, scriptContent string) error {
	scope := fmt.Sprintf("//%s", collectionName)

	_, err := c.Query(scope, scriptContent, nil)
	return err
}

// updateMigrationVersion takes a collection name and a new version, it then updates the version tracker
func (c *Client) updateMigrationVersion(collectionName string, version int) error {
	thingsName := genMigrationVersionThingName(collectionName)

	vars := map[string]interface{}{
		"version": version,
	}

	_, err := c.Query("//db_man", fmt.Sprintf(".%s = version", thingsName), vars)
	return err
}

// getMigrationVersion takes in a collection name and returns the current migration version for that collection
func (c *Client) getMigrationVersion(collectionName string) (int, error) {
	thingName := genMigrationVersionThingName(collectionName)

	res, err := c.Query("//db_man", fmt.Sprintf(".%s", thingName), nil)
	if err != nil {
		return 0, err
	}

	version, ok := res.(int8)
	if !ok {
		return 0, errors.New("error casting migration version response to int8")
	}

	return int(version), nil
}

// createMigrationVersionThings will create a new thing in the `db_man` collection to keep track of migrations
func (c *Client) createMigrationVersionThings(collectionNames []string) error {
	for _, v := range collectionNames {
		thingName := genMigrationVersionThingName(v)

		res, err := c.Query("//db_man", fmt.Sprintf("!is_err(try(.%s))", thingName), nil)
		if err != nil {
			return err
		}

		exists, ok := res.(bool)
		if !ok {
			return errors.New("error casting thing exists response into bool")
		}

		if exists {
			continue
		}

		_, err = c.Query("//db_man", fmt.Sprintf(".%s = 0", thingName), nil)
		if err != nil {
			return err
		}
	}

	return nil
}
